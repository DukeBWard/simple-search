package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/joho/godotenv"
)

func main() {
	env := godotenv.Load()
	if env != nil {
		panic("no env variables")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = ":8080"
	} else {
		port = ":" + port
	}

	// init fiber app
	app := fiber.New(fiber.Config{
		IdleTimeout: 10 * time.Second,
	})

	// add in some middleware
	app.Use(compress.New())

	// start server on its own go routine
	go func() {
		if err := app.Listen(port); err != nil {
			log.Panic(err)
		}
	}()

	// chan to listen to go routine
	c := make(chan os.Signal, 1)
	// listens for interrupt signals
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	// block main thread until interrupt
	<-c
	app.Shutdown()
	fmt.Print("Server shutting down")

}
