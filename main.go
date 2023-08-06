package main

import (
	"flag"
	"tmp/controllers"
	"tmp/framework"
)

func main() {
	host := flag.String("host", "localhost", "Host address for web server")
	port := flag.Uint("port", 8080, "Port number for web server")
	flag.Parse()

	engine := framework.NewEngine(*host, uint16(*port))
	router := engine.Router()
	router.Get("/users", controllers.UsersController)
	router.Get("/list", controllers.ListController)
	router.Get("/list/:item_id/:item_name", controllers.ListItemController)
	router.Get("/students", controllers.StudentsController)
	router.Get("/", controllers.TopPageController)
	router.Post("/post", controllers.PostController)
	router.Post("/users/:user_id/posts", controllers.UserPostController)
	router.Get("/json_p", controllers.JsonpController)

	router.Use(framework.TimeoutMiddleware)
	router.Use(framework.AuthUserMiddleware)
	router.Use(framework.TimeCostMiddleware)
	router.Use(framework.StaticFileMiddleware)
	router.UseNotFound(func(ctx *framework.MyContext) {
		ctx.WriteString("page is not found...")
	})

	engine.Run()
}
