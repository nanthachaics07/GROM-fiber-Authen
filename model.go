package main

import (
	"fmt"
	"log"

	"gorm.io/gorm"
)

type Book struct {
	gorm.Model
	Name        string
	Author      string
	Description string
	Peice       int
}

func CreateBook(db *gorm.DB, book *Book) error {
	result := db.Create(book)
	if result.Error != nil {
		log.Fatalf("Failed to create book: %v", result.Error)
	}
	fmt.Println("Book created successfully")
	return nil
}

func GetBook(db *gorm.DB, id uint) (*Book, error) {
	var books Book
	result := db.First(&books, id)
	if result.Error != nil {
		log.Fatalf("Failed to get book: %v", result.Error)
	}
	return &books, nil
}

func UpdateBook(db *gorm.DB, book *Book) error {
	result := db.Save(book)
	if result.Error != nil {
		log.Fatalf("Failed to update book: %v", result.Error)
	}
	fmt.Println("Book updated successfully")
	return nil
}
