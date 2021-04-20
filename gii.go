package gii

import (
	"net/http"
	"strings"
)

type Engine struct {
	*RouterGroup

	router *router
	groups []*RouterGroup
}

type RouterGroup struct {
	prefix      string
	parent      *RouterGroup
	middlewares []HandlerFunc
	engine      *Engine
}

func New() *Engine {
	e := &Engine{router: newRouter()}
	e.RouterGroup = &RouterGroup{engine: e, prefix: "/"}
	e.groups = []*RouterGroup{e.RouterGroup}

	return e
}

func Default() *Engine {
	e := New()
	e.Use(Logger(), Recovery())
	return e
}

func (g *RouterGroup) Group(prefix string) *RouterGroup {
	engine := g.engine
	new := &RouterGroup{
		prefix: g.prefix + prefix,
		parent: g,
		engine: engine,
	}
	engine.groups = append(engine.groups, new)
	return new
}

func (g *RouterGroup) addRoute(method, pattern string, handler HandlerFunc) {
	pattern = g.prefix + pattern
	g.engine.router.addRoute(method, pattern, handler)
}

func (g *RouterGroup) GET(pattern string, handler HandlerFunc) {
	g.addRoute("GET", pattern, handler)
}

func (g *RouterGroup) POST(pattern string, handler HandlerFunc) {
	g.addRoute("POST", pattern, handler)
}

func (g *RouterGroup) Use(middlewares ...HandlerFunc) {
	g.middlewares = append(g.middlewares, middlewares...)
}

func (e *Engine) Run(addr ...string) error {
	adds := ":8001"
	if len(addr) > 0 {
		adds = addr[0]
	}
	return http.ListenAndServe(adds, e)
}

func (e *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var handlers []HandlerFunc
	for _, g := range e.groups {
		if strings.HasPrefix(req.URL.Path, g.prefix) {
			handlers = append(handlers, g.middlewares...)
		}
	}

	c := newContext(w, req)
	c.handlers = handlers
	e.router.handle(c)
}
