package main

import (
	"fmt"
	"log"
	"my-web/gee"
	"net/http"
	"time"
)

func logger(c *gee.Context) {
	t := time.Now()
	c.Next()

	log.Println(time.Since(t))
}

func capture(c *gee.Context) {
	defer func() {
		if err := recover(); err != nil {
			message := fmt.Sprintf("%s", err)
			log.Println(message)
			c.Fail("Server error.")
		}
	}()

	c.Next()
}

func main() {
	r := gee.New()

	r.GET("/index", func(c *gee.Context) {
		c.String(http.StatusOK, "index")
	})
	hello := r.Group("/hello")
	hello.Use(logger, capture)
	hello.GET("/:name", func(c *gee.Context) {
		c.String(http.StatusOK, fmt.Sprintf("hello %v\n", c.Param("name")))
	})
	hello.GET("/:name/111", func(c *gee.Context) {
		data := []string{"a"}
		c.String(http.StatusOK, data[1])
	})
	//hello.GET("/world", func(c *gee.Context) {
	//	c.JSON(http.StatusOK, gee.H{
	//		"hello": "world",
	//	})
	//})

	r.Run(":9999")
}
