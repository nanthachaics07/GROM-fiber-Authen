package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "myuser"
	password = "mypassword"
	dbname   = "mydatabase"
)

func main() {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second, // Slow SQL threshold
			LogLevel:      logger.Info, // Log level
			Colorful:      true,        // Enable color
		},
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		log.Fatal(err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}
	defer sqlDB.Close()

	app := fiber.New()

	db.AutoMigrate(&Book{})
	fmt.Println("Migrated! Successfully!")

	app.Get("/books", func(c *fiber.Ctx) error {
		return c.JSON(GetAllBooks(db))
	})

	app.Listen(":8080")

	// CreateBook(db, &Book{
	// 	Name:        "Book 1",
	// 	Author:      "Author 1",
	// 	Description: "Description 1",
	// 	Peice:       100,
	// })

	// changeBook := GetBook(db, 1)

	// changeBook.Name = "Book 1 Updated"
	// changeBook.Author = "Author 1 Updated"
	// UpdateBook(db, changeBook)

	books := SearchBook(db, "Book")
	for _, book := range books {
		fmt.Println(book)
	}

	// DeleteBook(db, 3)

}
