package main

import "github.com/jinzhu/gorm"

type ToDo struct {
	gorm.Model
	Title       string
	Description string
}
