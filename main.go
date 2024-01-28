package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
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

func authenticationRequired(c *fiber.Ctx) error {
	cookie := c.Cookies("jwt")
	token, err := jwt.ParseWithClaims(cookie, &jwt.MapClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET_KEY")), nil
		})
	if err != nil || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}
	claims, ok := token.Claims.(*jwt.MapClaims)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}
	fmt.Println(claims)
	return c.Next()
}

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

	db.AutoMigrate(&Book{}, &User{})
	fmt.Println("Migrated! Successfully!")

	app := fiber.New()

	app.Use("/books", authenticationRequired)
	app.Use("/books/:id", authenticationRequired)
	app.Use("/search-books", authenticationRequired)
	app.Use("/update-books/:id", authenticationRequired)
	app.Use("/update-date-books/:id", authenticationRequired)
	app.Use("/delete-books/:id", authenticationRequired)
	app.Use("/force-delete-books/:id", authenticationRequired)

	app.Get("/books", func(c *fiber.Ctx) error {
		return c.JSON(GetAllBooks(db))
	})

	app.Get("/books/:id", func(c *fiber.Ctx) error {
		id, err := strconv.ParseUint(c.Params("id"), 10, 64)
		if err != nil {
			return c.Status(500).SendString(err.Error())
		}
		return c.JSON(GetBook(db, uint(id)))
	})

	app.Get("/search-books", func(c *fiber.Ctx) error {
		name := c.Query("name")
		return c.JSON(SearchBook(db, name))
	})

	app.Post("/update-books", func(c *fiber.Ctx) error {
		book := new(Book)
		if err := c.BodyParser(book); err != nil {
			return c.Status(500).SendString(err.Error())
		}
		err := CreateBook(db, book)
		if err != nil {
			return c.Status(500).SendString(err.Error())
		}
		return c.JSON(book)
	})

	app.Put("/update-books/:id", func(c *fiber.Ctx) error {
		id, err := strconv.ParseUint(c.Params("id"), 10, 64)
		if err != nil {
			return c.Status(500).SendString(err.Error())
		}
		book := new(Book)
		if err := c.BodyParser(book); err != nil {
			return c.Status(500).SendString(err.Error())
		}
		book.ID = uint(id)
		err = UpdateBook(db, book)
		if err != nil {
			return c.Status(500).SendString(err.Error())
		}
		return c.SendStatus(200)
	})

	app.Put("/update-date-books/:id", func(c *fiber.Ctx) error {
		id, err := strconv.ParseUint(c.Params("id"), 10, 64)
		if err != nil {
			return c.Status(500).SendString(err.Error())
		}
		book := new(Book)
		if err := c.BodyParser(book); err != nil {
			return c.Status(500).SendString(err.Error())
		}
		book.ID = uint(id)
		err = UpdateDateBook(db, book)
		if err != nil {
			return c.Status(500).SendString(err.Error())
		}
		return c.SendStatus(200)
	})

	app.Delete("/delete-books/:id", func(c *fiber.Ctx) error {
		id, err := strconv.ParseUint(c.Params("id"), 10, 64)
		if err != nil {
			return c.Status(500).SendString(err.Error())
		}
		DeleteBook(db, uint(id))
		return c.SendStatus(200)
	})

	app.Delete("/force-delete-books/:id", func(c *fiber.Ctx) error {
		id, err := strconv.ParseUint(c.Params("id"), 10, 64)
		if err != nil {
			return c.Status(500).SendString(err.Error())
		}
		ForceDeleteBook(db, uint(id))
		return c.SendStatus(200)
	})

	// User APIs
	app.Post("/register", func(c *fiber.Ctx) error {
		user := new(User)
		if err := c.BodyParser(user); err != nil {
			return c.Status(500).SendString(err.Error())
		}
		err := createUser(db, user)
		if err != nil {
			return c.Status(500).SendString(err.Error())
		}
		fmt.Println("User created successfully")
		return c.JSON(user)
	})

	app.Post("/login", func(c *fiber.Ctx) error {
		user := new(User)
		if err := c.BodyParser(user); err != nil {
			return c.Status(500).SendString(err.Error())
		}
		token, err := LoginUser(db, user)
		if err != nil {
			return c.Status(500).SendString(err.Error())
		}

		c.Cookie(&fiber.Cookie{
			Name:     "jwt",
			Value:    token,
			Expires:  time.Now().Add(time.Hour * 24),
			HTTPOnly: true,
		})
		// fmt.Println("User logged in successfully")
		return c.JSON(map[string]string{
			"message": "User logged in successfully",
			"token":   "Your token is: " + token,
		})
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
