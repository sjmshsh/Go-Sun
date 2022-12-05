package sun

import (
	"errors"
	"github.com/Go-Sun/sun/sun/render"
	"log"
	"net/http"
	"net/url"
	"strings"
)
// 32M
const defaultMaxMemory = 32 << 20

type Context struct {
	W          http.ResponseWriter
	R          *http.Request
	engine     *Engine
	queryCache url.Values
	formCache  url.Values
}

func (c *Context) GetQuery(key string) string {
	c.initQueryCache()
	return c.queryCache.Get(key)
}

func (c *Context) GetQueryArray(key string) ([]string, bool) {
	c.initQueryCache()
	values, ok := c.queryCache[key]
	return values, ok
}

func (c *Context) QueryArray(key string) (values []string) {
	c.initQueryCache()
	values, _ = c.queryCache[key]
	return
}

func (c *Context) GetDefaultQuery(key, defaultValue string) string {
	values, ok := c.GetQueryArray(key)
	if !ok {
		return defaultValue
	}
	return values[0]
}

func (c *Context) QueryMap(key string) (dicts map[string]string) {
	dicts, _ = c.GetQueryMap(key)
	return
}

func (c *Context) GetQueryMap(key string) (map[string]string, interface{}) {
	c.initQueryCache()
	return c.get(c.queryCache, key)
}

func (c *Context) initPostFormCache() {
	if c.R != nil {
		if err := c.R.ParseMultipartForm(defaultMaxMemory); err != nil {
			if !errors.Is(err, http.ErrNotMultipart) {
				log.Panicln(err)
			}
		}
		c.formCache = c.R.URL.Query()
	} else {
		c.formCache = url.Values{}
	}
}

// http://localhost:8080/queryMap?user[id]=1&user[name]=张三
func (c *Context) get(cache map[string][]string, key string) (map[string]string, bool) {
	// user[id]=1&user[name]=张三
	dicts := make(map[string]string)
	exist := false
	for k, value := range cache {
		// 判断有没有中括号，是不是map类型的
		if i := strings.IndexByte(k, '['); i >= 1 && k[0:i] == key {
			if j := strings.IndexByte(k[i+1:], ']'); j >= 0 {
				exist = true
				dicts[k[i+1:][:j]] = value[0]
			}
		}
	}
	return dicts, exist
}

func (c *Context) initQueryCache() {
	if c.R != nil {
		c.queryCache = c.R.URL.Query()
	} else {
		c.queryCache = url.Values{}
	}
}

func (c *Context) File(fileName string) {
	http.ServeFile(c.W, c.R, fileName)
}

func (c *Context) FileAttachment(filepath, filename string) {
	if isASCII(filename) {
		c.W.Header().Set("Content-Disposition", `attachment; filename="`+filename+`"`)
	} else {
		c.W.Header().Set("Content-Disposition", `attachment; filename*=UTF-8''`+url.QueryEscape(filename))
	}
	http.ServeFile(c.W, c.R, filepath)
}

// FileFromFS filepath是相对文件系统的路径
func (c *Context) FileFromFS(filepath string, fs http.FileSystem) {
	defer func(old string) {
		c.R.URL.Path = old
	}(c.R.URL.Path)

	c.R.URL.Path = filepath
	http.FileServer(fs).ServeHTTP(c.W, c.R)
}

func (c *Context) String(status int, format string, values ...any) error {
	err := c.Render(status, &render.String{
		Format: format,
		Data:   values,
	})
	return err
}

func (c *Context) XML(status int, data any) error {
	return c.Render(status, &render.XML{
		Data: data,
	})
}

func (c *Context) JSON(status int, data any) error {
	return c.Render(status, &render.JSON{Data: data})
}

func (c *Context) Render(statusCode int, s render.Render) error {
	err := s.Render(c.W)
	// 因为如果是200的话，那么系统会自动的给你去写，你这里再次写的话就会报错：superfluous。虽然不影响但是我们的框架中是不可以出现这种现象的
	if statusCode != http.StatusOK {
		c.W.WriteHeader(statusCode)
	}
	return err
}

func (c *Context) HTML(status int, html string) {
	c.Render(status, &render.HTML{
		IsTemplate: false,
		Data:       html,
	})
}

func (c *Context) HTMLTemplate(name string, data any) {
	c.Render(http.StatusOK, &render.HTML{
		IsTemplate: true,
		Name:       name,
		Data:       data,
		Template:   c.engine.HTMLRender.Template,
	})
}

func (c *Context) Redirect(status int, location string) {
	c.Render(status, &render.Redirect{
		Code:     status,
		Request:  c.R,
		Location: location,
	})
}
