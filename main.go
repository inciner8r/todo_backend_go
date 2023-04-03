package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/inciner8r/todo_backend_go/routes"
)

func main() {
	r := gin.Default()
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowHeaders = []string{"Authorization", "Content-Type"}
	r.Use(cors.New(config))
	routes.SetRoutes(r)
	r.Run(":8080")
}
