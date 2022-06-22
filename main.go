package main

import (
	"gee"
	"log"
	"net/http"
)

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
	{
		r.GET("/", func(c *gee.Context) {
			c.HTML(http.StatusOK, "<h1>Hello Gee</h1>")
		})

		r.GET("/index", func(c *gee.Context) {
			c.HTML(http.StatusOK, "<h1>Index Page</h1>")
		})
	}

	v1 := r.Group("/v1")
	{
		v1.GET("/", func(c *gee.Context) {
			c.HTML(http.StatusOK, "<h1>Hello Gee</h1>")
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
