package models

import "gorm.io/gorm"

type Todo struct {
	gorm.Model
	Title     string
	Content   string
	Author_id string
}
type User struct {
	gorm.Model
	Name     string
	Username string
	Email    string
	Password string
}
