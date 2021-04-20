package gii

import (
	"fmt"
	"strings"
)

type HandlerFunc func(*Context)

// roots key eg, roots['GET'] roots['POST']
// handlers key eg, handlers['GET-/p/:lang/doc'], handlers['POST-/p/book']
type router struct {
	roots    map[string]*node
	handlers map[string]HandlerFunc
}

func newRouter() *router {
	return &router{
		roots:    make(map[string]*node),
		handlers: make(map[string]HandlerFunc),
	}
}

func parsePattern(pattern string) []string {
	vs := strings.Split(pattern, "/")

	parts := make([]string, 0)
	for _, item := range vs {
		if item != "" {
			parts = append(parts, item)
			if item[0] == '*' {
				break
			}
		}
	}
	return parts
}

func (r *router) addRoute(method, pattern string, handler HandlerFunc) {
	parts := parsePattern(pattern)
	_, ok := r.roots[method]
	if !ok {
		r.roots[method] = &node{}
	}

	r.roots[method].insert(pattern, parts, 0)
	key := method + "-" + pattern
	r.handlers[key] = handler
}

func (r *router) getRoute(method, pattern string) (*node, map[string]string) {
	sparts := parsePattern(pattern)
	params := make(map[string]string)
	root, ok := r.roots[method]
	if !ok {
		return nil, nil
	}

	n := root.search(sparts, 0)
	if n != nil {
		parts := parsePattern(n.pattern)
		for i, v := range parts {
			if v[0] == ':' {
				params[v[1:]] = sparts[i]
			}
			if v[0] == '*' && len(v) > 1 {
				params[v[1:]] = strings.Join(sparts[i:], "/")
				break
			}
		}
		return n, params
	}

	return nil, nil
}

func (r *router) handle(c *Context) {
	n, params := r.getRoute(c.Method, c.Path)
	if n != nil {
		c.Params = params
		key := c.Method + "-" + n.pattern
		c.handlers = append(c.handlers, r.handlers[key])
	} else {
		c.handlers = append(c.handlers, func(c *Context) {
			fmt.Fprintf(c.Writer, "404 NOT FOUND: %s\n", c.Path)
		})
	}

	c.Next()
}
