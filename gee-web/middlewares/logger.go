package middlewares

import (
	"gee-web/gee"
	"log"
	"time"
)

func Logger(c *gee.Context) {
	t := time.Now()
	c.Next()
	log.Printf("[%d] %s in %v", c.StatusCode, c.Request.RequestURI, time.Since(t))
}