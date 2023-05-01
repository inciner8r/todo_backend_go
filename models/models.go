package models

import (
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

type Todo struct {
	gorm.Model
	Content   string
	Author_id uint
	Completed bool
}
type User struct {
	gorm.Model
	Name     string
	Username string
	Email    string
	Password string
}

type Claims struct {
	Id uint
	jwt.RegisteredClaims
}
type TokenString struct {
	JWTString string
}
