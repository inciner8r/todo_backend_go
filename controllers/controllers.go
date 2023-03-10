package controllers

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/inciner8r/todo_backend_go/models"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func initDb() *gorm.DB {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	username := os.Getenv("MYSQL_USER")
	password := os.Getenv("MYSQL_PASSWORD")
	dbName := os.Getenv("MYSQL_DATABASE")
	dsn := username + ":" + password + "@tcp(127.0.0.1:3306)/" + dbName + "?parseTime=true"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalln(err)
	}
	db.AutoMigrate(&models.Todo{})
	db.AutoMigrate(&models.User{})
	return db
}

var db = initDb()

func Home(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"data": "ok"})
}
func AddTodo(c *gin.Context) {
	var todo models.Todo
	if err := c.BindJSON(&todo); err != nil {
		log.Fatal(err.Error())
	}
	if err := db.Create(&todo).Error; err != nil {
		log.Fatal(err.Error())
	}
	c.JSON(http.StatusOK, gin.H{"data": todo})
}
func MyTodos(c *gin.Context) {
	auth_id, _ := c.Params.Get("user")
	var todos []models.Todo
	if err := db.Find(&todos, "Author_id = ?", auth_id).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
}
func AllTodos(c *gin.Context) {
	var todos []models.Todo
	if err := db.Find(&todos).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": todos})
}
