package main

import (
	"fmt"
	"github.com/Go-Sun/sun/sun"
	"net/http"
)

func Log(next sun.HandlerFunc) sun.HandlerFunc {
	return func(ctx *sun.Context) {
		fmt.Println("pre handler")
		next(ctx)
		fmt.Println("post handler")
	}
}

func main() {
	r := sun.New()
	g := r.Group("user")
	g.Get("/json", func(ctx *sun.Context) {
		user := &struct {
			Name     string
			Password string
		}{
			Name:     "lxy",
			Password: "520",
		}
		ctx.JSON(http.StatusOK, user)
	})
	g.Get("/xml", func(ctx *sun.Context) {
		user := &struct {
			Name     string `xml:"name"`
			Password string `xml:"password"`
		}{
			Name:     "lxy",
			Password: "520",
		}
		ctx.XML(http.StatusOK, user)
	})
	g.Get("/excel", func(ctx *sun.Context) {
		ctx.File("test.xlsx")
	})
	g.Get("/excelName", func(ctx *sun.Context) {
		ctx.FileAttachment("test.xlsx", "aaa")
	})
	g.Get("/fs", func(ctx *sun.Context) {
		ctx.FileFromFS("test/xlsx", http.Dir("tpl"))
	})
	g.Get("/redirect", func(ctx *sun.Context) {
		ctx.Redirect(http.StatusFound, "/user/json")
	})
	g.Get("/string", func(ctx *sun.Context) {
		ctx.JSON(http.StatusOK, "heheheheheh")
	})
	g.Get("/add", func(ctx *sun.Context) {
		id, _ := ctx.GetQueryArray("id")
		fmt.Println(id)
	})
	g.Get("/queryMap", func(ctx *sun.Context) {
		m, _ := ctx.GetQueryMap("user")
		ctx.JSON(http.StatusOK, m)
	})
	r.Run()
}
