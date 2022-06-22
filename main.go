package main

import (
	"fmt"
	"gee"
	"html/template"
	"log"
	"net/http"
	"time"
)

type student struct {
	Name string
	Age  int8
}

func FormatAsDate(t time.Time) string {
	year, month, day := t.Date()
	return fmt.Sprintf("%d-%02d-%02d", year, month, day)
}

func firstForV2() gee.HandlerFunc {
	return func(c *gee.Context) {
		log.Printf("[%d] %s first for group v2", c.StatusCode, c.Req.RequestURI)
		c.Next()
		log.Printf("[%d] %s fourth for group v2", c.StatusCode, c.Req.RequestURI)
	}
}

func secondForV2() gee.HandlerFunc {
	return func(c *gee.Context) {
		log.Printf("[%d] %s second for group v2", c.StatusCode, c.Req.RequestURI)
		c.Next()
		log.Printf("[%d] %s third for group v2", c.StatusCode, c.Req.RequestURI)
	}
}

func main() {
	r := gee.New()
	r.Use(gee.Logger())
	r.Static("/assets", "./static")
	// 自定义映射函数
	r.SetFuncMap(template.FuncMap{
		"FormatAsDate": FormatAsDate,
	})
	// 引入所有的模板
	r.LoadHTMLGlob("templates/*")
	{
		r.GET("/", func(c *gee.Context) {
			c.HTML(http.StatusOK, "css.tmpl", nil)
		})

		r.GET("/date", func(c *gee.Context) {
			c.HTML(http.StatusOK, "custom_func.tmpl", gee.Map{
				"title": "gee",
				"now":   time.Now(),
			})
		})
	}

	v1 := r.Group("/v1")
	{
		stu1 := &student{Name: "CC", Age: 20}
		stu2 := &student{Name: "Jack", Age: 22}
		v1.GET("/students", func(c *gee.Context) {
			c.HTML(http.StatusOK, "arr.tmpl", gee.Map{
				"title":  "gee",
				"stuArr": [2]*student{stu1, stu2},
			})
		})
		v1.GET("/hello", func(c *gee.Context) {
			// expect /hello?name=gcc
			c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
		})
	}
	v2 := r.Group("/v2")
	v2.Use(firstForV2(), secondForV2())
	{
		v2.GET("/hello/:name", func(c *gee.Context) {
			// expect /hello/cc
			c.String(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
		})
		v2.GET("/file/*filepath", func(c *gee.Context) {
			// expect /file/asserts/cc/jquery.js
			c.String(http.StatusOK, "visit %s, you're at %s\n", c.Param("filepath"), c.Path)
		})
		v2.POST("/login", func(c *gee.Context) {
			c.JSON(http.StatusOK, gee.Map{
				"username": c.PostForm("username"),
				"password": c.PostForm("password"),
			})
		})

	}

	err := r.Run(":9999")
	if err != nil {
		return
	}
}
