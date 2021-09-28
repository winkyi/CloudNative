package engine

type router struct {
	handlers map[string]HandlerFunc
}

func newRouter() *router {
	return &router{handlers: make(map[string]HandlerFunc)}
}

func (r *router) addRouter(method string, pattern string, handler HandlerFunc) {
	key := method + "-" + pattern
	r.handlers[key] = handler
}

func (r *router) handle(c *Context) {
	key := c.Method + "-" + c.Pattern
	handler, ok := r.handlers[key]
	if !ok {
		c.SetHeaders(c.R.Header)
		c.SetEnvToResponseHeader("VERSION")
		c.HTML(404, "<h1>404 页面不存在</h1>")
		c.Log("404 页面不存在")
	} else {
		handler(c)
	}
}
