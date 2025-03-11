package main

import (
	"fmt"
	"github.com/Sunshine9d/go-inventory/config"
	"log"
)

func main() {
	app, err := config.InitializeApp(config.DbUser, config.DbPassword, config.DbName)
	if err != nil {
		log.Fatal("Failed to initialize app:", err)
	}
	fmt.Println("Server is running on", config.DbHost)
	config.RunApp(app, ":8088")
}
