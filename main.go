package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/inciner8r/todo_backend_go/routes"
)

func main() {
	r := gin.Default()
	routes.SetRoutes(r)
	r.Use(cors.Default())
	r.Run(":8080")
}
