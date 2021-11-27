package test

import "gee"

func newTestRouter() *gee.Engine {
	g := gee.New()
	g.GET("/", nil)
	return g
}
