package main

import (
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


	// if fiber.IsChild() {
	// 	fmt.Println("I'm a child")
	// } else {
	// 	fmt.Println("I'm a parent")
	// }

	app.Get("/", func(c *fiber.Ctx) error {
        return c.SendString("Hello, World 👋!")
    })

	err := app.Listen("Localhost:3000")
	if err != nil {
		panic(err)
	}
}