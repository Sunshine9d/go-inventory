package main

import (
	"fmt"
	"github.com/Sunshine9d/go-inventory/config"
)

func main() {
	app := config.App{}
	app.Initialize(config.DbUser, config.DbPassword, config.DbName)
	fmt.Println("Server is running on", config.DbHost)
	app.Run(config.DbHost)
}
