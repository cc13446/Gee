package gee

import (
	"log"
	"time"
)

// Logger 日志中间件
func Logger() HandlerFunc {
	return func(c *Context) {
		// Start timer
		t := time.Now()
		// 下一个中间件或者URL处理函数
		c.Next()
		// Calculate resolution time
		log.Printf("[%d] %s in %v", c.StatusCode, c.Req.RequestURI, time.Since(t))
	}
}
