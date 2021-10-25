package main

import (
	"gee-web/gee"
	"net/http"
)

var r *gee.Engine



func main() {
	r := gee.New()
	r.GET("/", func(c *gee.Context) {
		c.Html(http.StatusOK, "<h2>Hello Gee</h2>\n")
	})

	r.GET("/about", func(c *gee.Context) {
		c.String(http.StatusOK, "about me ?\n")
	})

	r.POST("/user", func(c *gee.Context) {
		c.Json(http.StatusOK, map[string]string{
			"a" : "b",
			"c" : "d",
		})
	})
	r.Run(":9999")
}
