package gee

import (
	"log"
	"net/http"
	"strings"
)

type HandleFunc func(c *Context)

type RouterGroup struct {
	prefix      string
	middlewares []HandleFunc
	parent      *RouterGroup
	engine      *Engine
}

type Engine struct {
	*RouterGroup
	router *Router
	groups []*RouterGroup
}

func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var middlewares []HandleFunc
	for _, group := range engine.groups {
		if strings.HasPrefix(req.URL.Path, group.prefix) {
			middlewares = append(middlewares, group.middlewares...)
		}
	}
	c := newContext(w, req)
	c.handlers = middlewares
	engine.router.handle(c)
}

func New() *Engine {
	engine := &Engine{
		router: newRouter(),
	}
	engine.RouterGroup = &RouterGroup{engine: engine}
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
}

func (group *RouterGroup) Group(prefix string) *RouterGroup {
	engine := group.engine
	newGroup := &RouterGroup{
		engine: engine,
		prefix: group.prefix + prefix,
		parent: group,
	}
	engine.groups = append(engine.groups, newGroup)
	return newGroup
}

func (group *RouterGroup) addRoute(method string, path string, handle HandleFunc) {
	pattern := group.prefix + path
	// fmt.Println(path, "     ", group.prefix)
	log.Printf("Route %4s - %s", method, pattern)
	group.engine.router.addRoute(method, pattern, handle)
}

func (group *RouterGroup) GET(pattern string, handle HandleFunc) {
	group.addRoute("GET", pattern, handle)
}

func (group *RouterGroup) POST(pattern string, handle HandleFunc) {
	group.addRoute("POST", pattern, handle)
}

// middlewares
func (group *RouterGroup) Use(middlewares ...HandleFunc) {
	group.middlewares = append(group.middlewares, middlewares...)
}

func (engine *Engine) Run(add string) (err error) {
	return http.ListenAndServe(add, engine)
}
