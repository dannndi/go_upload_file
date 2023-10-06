package main

import (
	"os"

	"github.com/dannndi/go_upload_file/core/utils"
	"github.com/dannndi/go_upload_file/module/file"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	config := fiber.Config{
		// 100MB
		BodyLimit: 100 * 1024 * 1024,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			if err == nil {
				return nil
			}

			code := fiber.ErrInternalServerError.Code
			message := fiber.ErrInternalServerError.Message

			if value, ok := err.(*fiber.Error); ok {
				code = value.Code
				message = value.Message
			}

			return c.Status(code).JSON(utils.BaseResponse{
				Message: message,
				Error:   nil,
			})
		},
	}

	app := fiber.New(config)
	app.Use(cors.New())
	app.Static("/", "./public")
	app.Static("/uploads", "./public/uploads")

	v1 := app.Group("/api/v1")
	v1.Route("/file", file.Route)

	port := ":" + os.Getenv("PORT")
	app.Listen(port)
}
