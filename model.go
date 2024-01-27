package main

import (
	grom "gorm.io/gorm"
)

type Book struct {
	grom.Model
	Name        string
	Author      string
	Description string
	Peice       int
}
