package controllers

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/inciner8r/todo_backend_go/models"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
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
func AddUser(c *gin.Context) {
	var user models.User

	if err := c.BindJSON(&user); err != nil {
		log.Fatal(err.Error())
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println("hashing error")
		return
	}

	user.Password = string(hash)

	if err := db.Create(&user).Error; err != nil {
		log.Fatal(err.Error())
	}
	c.JSON(http.StatusOK, gin.H{"data": user, "message": "registered"})
}

var key = []byte("key")

func Login(c *gin.Context) {
	var user models.User
	var expected models.User

	if err := c.BindJSON(&user); err != nil {
		log.Fatal(err.Error())
	}

	if err := db.Table("users").Where("username = ?", user.Username).First(&expected).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"data": err})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(expected.Password), []byte(user.Password)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"data": err.Error()})
		return
	}
	claims := models.Claims{
		Id: user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(key)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"data": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"jwt": signedToken, "message": "logged in"})
}

func ValidateJwt(c *gin.Context) {
	var tokenString models.TokenString
	if err := c.BindJSON(&tokenString); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "cookie not found"})
		return
	}
	token, err := jwt.ParseWithClaims(tokenString.JWTString, &models.Claims{}, func(t *jwt.Token) (interface{}, error) {
		return key, nil
	})
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
	}
	claims := token.Claims.(*models.Claims)
	c.JSON(http.StatusOK, gin.H{"Id": claims.Id})
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
func DeleteTodo(c *gin.Context) {

	var todo models.Todo

	if err := db.Where("name = ?", c.Param("id")).First(&todo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "record not found"})
		return
	}
	db.Delete(&todo)
	c.JSON(http.StatusOK, gin.H{"data": todo})
}
func EditTodo(c *gin.Context) {

	var todo models.Todo

	if err := c.BindHeader(&todo); err != nil {
		log.Fatal(err.Error())
	}
	//db.Where()
}
