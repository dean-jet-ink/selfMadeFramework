package framework

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
)

type Engine struct {
	host   string
	port   uint16
	router *Router
}

func NewEngine(host string, port uint16) *Engine {
	return &Engine{
		host:   host,
		port:   port,
		router: NewRouter(),
	}
}

func (e *Engine) Router() *Router {
	return e.router
}

func (e *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	pathName := strings.TrimSuffix(r.URL.Path, "/")

	routingTable := e.router.routingTables[strings.ToLower(r.Method)]
	node := routingTable.Search(pathName)

	ctx := NewMyContext(w, r)

	var targetHandler func(*MyContext)
	if node == nil || node.handler == nil {
		targetHandler = e.router.NotFoundHandler()
	} else {
		targetHandler = node.handler
		pathParams := routingTable.ParsePath(pathName, node)
		ctx.SetPathParams(pathParams)
	}

	handlers := append(e.router.Middlewares(), targetHandler)
	ctx.SetHandlers(handlers)
	ctx.Next()
}

func (e *Engine) Run() {
	ch := make(chan os.Signal)
	signal.Notify(ch)

	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", e.host, e.port),
		Handler: e,
	}
	go func() {
		server.ListenAndServe()
	}()

	<-ch
	log.Println("shutdown")

	if err := server.Shutdown(context.Background()); err != nil {
		log.Println(err)
	}

	log.Println("shutdown successfully")
}

type Router struct {
	routingTables   map[string]*TrieTree
	middlewares     []func(*MyContext)
	notFoundHandler func(*MyContext)
}

func NewRouter() *Router {
	return &Router{
		routingTables: map[string]*TrieTree{
			"get":    NewTrieTree(),
			"post":   NewTrieTree(),
			"put":    NewTrieTree(),
			"delete": NewTrieTree(),
			"patch":  NewTrieTree(),
		},
		notFoundHandler: func(ctx *MyContext) {
			ctx.WriteHeader(http.StatusNotFound)
			log.Println("page not found")
		},
	}
}

func (r *Router) Middlewares() []func(*MyContext) {
	return r.middlewares
}

func (r *Router) register(method string, path string, handler func(*MyContext)) {
	routingTable := r.routingTables[method]
	path = strings.TrimSuffix(path, "/")

	node := routingTable.Search(path)
	if node != nil && node.handler != nil {
		panic(fmt.Sprintf("%s: handler already exist", path))
	}

	routingTable.Insert(path, handler)
}

func (r *Router) Get(path string, handler func(*MyContext)) {
	r.register("get", path, handler)
}

func (r *Router) Post(path string, handler func(*MyContext)) {
	r.register("post", path, handler)
}

func (r *Router) Put(path string, handler func(*MyContext)) {
	r.register("put", path, handler)
}

func (r *Router) Delete(path string, handler func(*MyContext)) {
	r.register("delete", path, handler)
}

func (r *Router) Patch(path string, handler func(*MyContext)) {
	r.register("patch", path, handler)
}

func (r *Router) Use(fn func(*MyContext)) {
	r.middlewares = append(r.middlewares, fn)
}

func (r *Router) UseNotFound(fn func(*MyContext)) {
	r.notFoundHandler = fn
}

func (r *Router) NotFoundHandler() func(*MyContext) {
	return r.notFoundHandler
}
