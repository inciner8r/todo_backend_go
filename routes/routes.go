package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/inciner8r/todo_backend_go/controllers"
)

func SetRoutes(r *gin.Engine) {
	r.GET("/", controllers.Home)
}
