package main

import (
	// "fmt"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New(fiber.Config{
		IdleTimeout: time.Second * 5,
		WriteTimeout: time.Second * 5,
		ReadTimeout: time.Second * 5,
		// Prefork: true,
	})

	// Akan berjalan di endpoint yang ada /api nya
	app.Use("/api",func (ctx *fiber.Ctx) error  {
		fmt.Println("I'm middleware before processing request")
		err := ctx.Next()
		fmt.Println("I'm middleware after processing request")
		return err
	})


	// if fiber.IsChild() {
	// 	fmt.Println("I'm a child")
	// } else {
	// 	fmt.Println("I'm a parent")
	// }

	app.Get("/api/hello", func(c *fiber.Ctx) error {
        return c.SendString("Hello, World 👋!")
    })

	err := app.Listen("Localhost:3000")
	if err != nil {
		panic(err)
	}
}