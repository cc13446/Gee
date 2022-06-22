package gee

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Map map[string]interface{}

type Context struct {

	// origin objects
	Writer http.ResponseWriter
	Req    *http.Request

	// request info
	Path   string
	Method string
	Params map[string]string

	// response info
	StatusCode int

	// middleware
	handlers []HandlerFunc
	index    int
}

func newContext(w http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		Writer: w,
		Req:    req,
		Path:   req.URL.Path,
		Method: req.Method,
		index:  -1,
	}
}

func (c *Context) status(code int) {
	c.StatusCode = code
	c.Writer.WriteHeader(code)
}

// Param 获取URL参数
func (c *Context) Param(key string) string {
	if value, ok := c.Params[key]; ok {
		return value
	}
	return ""
}

// PostForm 从Post和Put主体中获取value
func (c *Context) PostForm(key string) string {
	return c.Req.FormValue(key)
}

// Query 从URL获取value
func (c *Context) Query(key string) string {
	return c.Req.URL.Query().Get(key)
}

func (c *Context) SetHeader(key string, value string) {
	c.Writer.Header().Set(key, value)
}

func (c *Context) String(code int, format string, values ...interface{}) {
	c.SetHeader("Content-Type", "text/plain")
	c.status(code)
	_, err := c.Writer.Write([]byte(fmt.Sprintf(format, values...)))
	if err != nil {
		panic(err)
	}
}

func (c *Context) JSON(code int, obj interface{}) {
	c.SetHeader("Content-Type", "application/json")
	c.status(code)
	encoder := json.NewEncoder(c.Writer)
	if err := encoder.Encode(obj); err != nil {
		panic(err)
	}
}

func (c *Context) Data(code int, data []byte) {
	c.status(code)
	_, err := c.Writer.Write(data)
	if err != nil {
		panic(err)
	}
}

func (c *Context) HTML(code int, html string) {
	c.SetHeader("Content-Type", "text/html")
	c.status(code)
	_, err := c.Writer.Write([]byte(html))
	if err != nil {
		panic(err)
	}
}

func (c *Context) Fail(code int, err string) {
	c.index = len(c.handlers)
	c.JSON(code, Map{"message": err})
}

// Next 依次调用中间件处理函数
func (c *Context) Next() {
	c.index++
	s := len(c.handlers)
	// 这里循环调用而不是简单调用下一个的原因是：
	// 不是所有的 handler都会调用 Next()
	// 手工调用 Next()，一般用于在请求前后各实现一些行为。
	// 如果中间件只作用于请求前，中间件可以省略调用Next()
	for ; c.index < s; c.index++ {
		c.handlers[c.index](c)
	}
}
