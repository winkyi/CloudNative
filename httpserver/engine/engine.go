package engine

import (
	"net/http"
)

// HandlerFunc 定义路由隐射方法
type HandlerFunc func(*Context)

type engine struct {
	// 路由
	router *router
}

func New() *engine {
	return &engine{router: newRouter()}
}

// addRoute 添加路由
func (e *engine) addRoute(method string, pattern string, handler HandlerFunc) {
	e.router.addRouter(method, pattern, handler)
}

// POST post请求
func (e *engine) POST(pattern string, handler HandlerFunc) {
	e.addRoute("POST", pattern, handler)
}

// GET get请求
func (e *engine) GET(pattern string, handler HandlerFunc) {
	e.addRoute("GET", pattern, handler)
}

func (e *engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := newContext(w, r)
	e.router.handle(c)
}

func (e *engine) Run(addr string) error {
	return http.ListenAndServe(addr, e)
}
