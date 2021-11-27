package gee

import (
	"net/http"
	"strings"
)

type Router struct {
	handlers map[string]HandleFunc
	roots    map[string]*node
}

// parsePattern 根据路径分隔符 / 解析路由路径
// eg /p/:lang/*filepath/ return []string{"p",":lang", "*filepath"}
func parsePattern(pattern string) []string {
	vs := strings.Split(pattern, "/")
	parts := make([]string, 0)
	for _, part := range vs {
		if part != "" {
			parts = append(parts, part)
			if part[0] == '*' {
				break
			}
		}
	}
	return parts
}

// roots key eg, roots['GET'] roots['POST']
// handlers key eg, handlers['GET-/p/:lang/doc'], handlers['POST-/p/book']
func newRouter() *Router {
	return &Router{
		handlers: make(map[string]HandleFunc),
		roots:    make(map[string]*node),
	}
}
func (r *Router) addRoute(method string, pattern string, handle HandleFunc) {
	parts := parsePattern(pattern)
	key := strings.Join([]string{method, "-", pattern}, "")
	_, ok := r.roots[method]
	if !ok {
		r.roots[method] = &node{}
	}
	r.roots[method].insert(pattern, parts, 0)
	r.handlers[key] = handle
}

func (r *Router) getRoute(method string, path string) (*node, map[string]string) {
	searchParts := parsePattern(path)
	params := make(map[string]string)
	root, ok := r.roots[method]
	if !ok {
		return nil, nil
	}
	n := root.search(searchParts, 0)
	if n != nil {
		parts := parsePattern(n.pattern)
		for index, part := range parts {
			if part[0] == ':' {
				params[part[1:]] = searchParts[index]
			}
			if part[0] == '*' {
				params[part[1:]] = strings.Join(searchParts[index:], "/")
			}
		}
		return n, params
	}
	return nil, nil
}

func (r *Router) handle(c *Context) {
	n, params := r.getRoute(c.Method, c.Path)
	if n == nil {
		c.handlers = append(c.handlers, func(c *Context) {
			c.String(http.StatusNotFound, "404 Not Found: %s\n", c.Path)
		})
	} else {
		c.Params = params
		key := strings.Join([]string{c.Method, "-", n.pattern}, "")
		c.handlers = append(c.handlers, func(c *Context) {
			r.handlers[key](c)
		})
	}
	c.Next()
}
