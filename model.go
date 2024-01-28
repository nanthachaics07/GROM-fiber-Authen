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

func CreateBook(db *gorm.DB, book *Book) {
	result := db.Create(book)
	if result.Error != nil {
		log.Fatalf("Error creating book: %v", result.Error)
	}
	fmt.Println("Book created successfully")
}

func GetBook(db *gorm.DB, id uint) *Book {
	var book Book
	result := db.First(&book, id)
	if result.Error != nil {
		log.Fatalf("Error finding book: %v", result.Error)
	}
	return &book
}

func GetAllBooks(db *gorm.DB) []Book {
	var books []Book
	result := db.Find(&books)
	if result.Error != nil {
		log.Fatalf("Error finding books: %v", result.Error)
	}
	return books
}

func UpdateBook(db *gorm.DB, book *Book) {
	result := db.Save(book)
	if result.Error != nil {
		log.Fatalf("Error updating book: %v", result.Error)
	}
	fmt.Println("Book updated successfully")
}

func DeleteBook(db *gorm.DB, id uint) {
	var book Book
	// result := db.Unscoped().Delete(&book, id)
	result := db.Delete(&book, id)
	if result.Error != nil {
		log.Fatalf("Error deleting book: %v", result.Error)
	}
	fmt.Println("Book deleted successfully")
}

func SearchBook(db *gorm.DB, name string) []Book {
	var books []Book
	result := db.Where("name like ?", "%"+name+"%").Order("id asc").Find(&books)
	if result.Error != nil {
		log.Fatalf("Error searching book: %v", result.Error)
	}
	return books
}
