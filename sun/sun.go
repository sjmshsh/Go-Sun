package sun

import (
	"fmt"
	"github.com/Go-Sun/sun/sun/render"
	"html/template"
	"log"
	"net/http"
	"sync"
)

const ANY = "ANY"

// HandlerFunc 处理函数的定义
type HandlerFunc func(ctx *Context)

type MiddlewareFunc func(handlerFunc HandlerFunc) HandlerFunc

// 分组
type routerGroup struct {
	name               string
	handleFuncMap      map[string]map[string]HandlerFunc
	middlewaresFuncMap map[string]map[string][]MiddlewareFunc
	handlerMethodMap   map[string][]string
	treeNode           *treeNode
	middlewares        []MiddlewareFunc
}

func (r *routerGroup) Use(middlewareFunc ...MiddlewareFunc) {
	r.middlewares = append(r.middlewares, middlewareFunc...)
}

func (r *routerGroup) methodHandle(name string, method string, h HandlerFunc, ctx *Context) {
	//组通用中间件
	if r.middlewares != nil {
		for _, middlewareFunc := range r.middlewares {
			h = middlewareFunc(h)
		}
	}
	//组路由级别
	middlewareFuncs := r.middlewaresFuncMap[name][method]
	if middlewareFuncs != nil {
		for _, middlewareFunc := range middlewareFuncs {
			h = middlewareFunc(h)
		}
	}
	h(ctx)
}

// user   get->handle
// goods
// order
type router struct {
	routerGroups []*routerGroup
}

func (r *router) Group(name string) *routerGroup {
	routerGroup := &routerGroup{
		name:               name,
		handleFuncMap:      make(map[string]map[string]HandlerFunc),
		middlewaresFuncMap: make(map[string]map[string][]MiddlewareFunc),
		handlerMethodMap:   make(map[string][]string),
		treeNode: &treeNode{
			name:     "/",
			children: make([]*treeNode, 0),
		},
	}
	r.routerGroups = append(r.routerGroups, routerGroup)
	return routerGroup
}

func (r *routerGroup) handle(name string, method string, handleFunc HandlerFunc, middlewareFunc ...MiddlewareFunc) {
	_, ok := r.handleFuncMap[name]
	if !ok {
		r.handleFuncMap[name] = make(map[string]HandlerFunc)
		r.middlewaresFuncMap[name] = make(map[string][]MiddlewareFunc)
	}
	_, ok = r.handleFuncMap[name][method]
	if ok {
		log.Fatal("有重复的路由")
	}
	r.handleFuncMap[name][method] = handleFunc
	r.middlewaresFuncMap[name][method] = append(r.middlewaresFuncMap[name][method], middlewareFunc...)
	r.treeNode.Put(name)
}

func (r *routerGroup) Any(name string, handlerFunc HandlerFunc, middlewareFunc ...MiddlewareFunc) {
	r.handle(name, ANY, handlerFunc, middlewareFunc...)
}

func (r *routerGroup) Get(name string, handlerFunc HandlerFunc, middlewareFunc ...MiddlewareFunc) {
	r.handle(name, http.MethodGet, handlerFunc, middlewareFunc...)
}

func (r *routerGroup) Post(name string, handlerFunc HandlerFunc, middlewareFunc ...MiddlewareFunc) {
	r.handle(name, http.MethodPost, handlerFunc, middlewareFunc...)
}

func (r *routerGroup) Delete(name string, handlerFunc HandlerFunc, middlewareFunc ...MiddlewareFunc) {
	r.handle(name, http.MethodPost, handlerFunc, middlewareFunc...)
}

func (r *routerGroup) Patch(name string, handlerFunc HandlerFunc, middlewareFunc ...MiddlewareFunc) {
	r.handle(name, http.MethodPost, handlerFunc, middlewareFunc...)
}

func (r *routerGroup) Options(name string, handlerFunc HandlerFunc, middlewareFunc ...MiddlewareFunc) {
	r.handle(name, http.MethodPost, handlerFunc, middlewareFunc...)
}

func (r *routerGroup) Head(name string, handlerFunc HandlerFunc, middlewareFunc ...MiddlewareFunc) {
	r.handle(name, http.MethodPost, handlerFunc, middlewareFunc...)
}

type Engine struct {
	router
	funcMap    template.FuncMap
	HTMLRender render.HTMLRender
	pool       sync.Pool
}

func New() *Engine {
	engine := &Engine{
		router: router{},
	}
	engine.pool.New = func() any {
		return engine.allocateContext()
	}
	return engine
}

func (e *Engine) allocateContext() any {
	return &Context{engine: e}
}

func (e *Engine) SetFuncMap(funcMap template.FuncMap) {
	e.funcMap = funcMap
}

func (e *Engine) LoadTemplate(pattern string) {
	t := template.Must(template.New("").Funcs(e.funcMap).ParseGlob(pattern))
	e.SetHtmlTemplate(t)
}

func (e *Engine) SetHtmlTemplate(t *template.Template) {
	e.HTMLRender = render.HTMLRender{Template: t}
}

func (c *Context) Template(name string, data any) error {
	c.W.Header().Set("Content-Type", "text/html; charset=utf-8")
	err := c.engine.HTMLRender.Template.ExecuteTemplate(c.W, name, data)
	return err
}

func (e *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := e.pool.Get().(*Context)
	ctx.W = w
	ctx.R = r
	e.httpRequestHandle(ctx, w, r)
	e.pool.Put(ctx)
}

func (e *Engine) httpRequestHandle(ctx *Context, w http.ResponseWriter, r *http.Request) {
	method := r.Method
	for _, group := range e.routerGroups {
		routerName := SubStringLast(r.URL.Path, "/"+group.name)
		// get/1
		node := group.treeNode.Get(routerName)

		if node != nil && node.isEnd == true {
			// 先判单ANY有没有
			handle, ok := group.handleFuncMap[node.routerName][ANY]
			if ok {
				// 有，就直接处理
				group.methodHandle(node.routerName, ANY, handle, ctx)
				return
			}
			// 然后判断method有没有
			handle, ok = group.handleFuncMap[node.routerName][method]
			if ok {
				group.methodHandle(node.routerName, method, handle, ctx)
				return
			}
			w.WriteHeader(http.StatusMethodNotAllowed)
			fmt.Fprintf(w, "%s %s not allowed \n", r.RequestURI, method)
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprintf(w, "%s not Found \n", r.RequestURI)
}

func (e *Engine) Run() {
	http.Handle("/", e)
	err := http.ListenAndServe(":8011", nil)
	if err != nil {
		log.Fatal(err)
	}
}
