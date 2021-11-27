package main

import (
	"gee"
	"net/http"
)

func main() {
	r := gee.New()
	r.Use(gee.Logger())
	r.GET("/", func(c *gee.Context) {
		c.HTML(http.StatusOK, "<h1>Hello H1</h1>")
	})

	r.GET("/hello", func(c *gee.Context) {
		c.String(http.StatusOK, "hello %s, you are at %s \n", c.Query("name"), c.Path)
	})
	r.GET("/hello/:name", func(c *gee.Context) {
		c.String(http.StatusOK, "hello %s, you are at %s \n", c.Query("name"), c.Path)
	})

	r.GET("/assets/*filapath", func(c *gee.Context) {
		c.JSON(http.StatusOK, gee.H{
			"filapath": c.Param("filapath"),
		})
	})

	r.POST("/login", func(c *gee.Context) {
		c.JSON(http.StatusOK, gee.H{
			"username": c.PostForm("usernam"),
			"password": c.PostForm("password"),
		})
	})

	v1 := r.Group("/v1")
	v1.GET("/", func(c *gee.Context) {
		c.HTML(http.StatusOK, "<h1>v1</h1>")
	})

	r.Run(":9999")
}
