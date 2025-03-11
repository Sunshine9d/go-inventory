package config

import (
	"database/sql"
	"fmt"
	"github.com/Sunshine9d/go-inventory/services"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func InitializeApp(user, password, dbname string) (*services.App, error) {
	connectionString := fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, password, dbname)
	var err error
	// Initialize App
	app := &services.App{
		Router: mux.NewRouter(),
		DB:     nil,
	}
	app.DB, err = sql.Open("mysql", connectionString)
	if err != nil {
		log.Fatal(err)
	}
	app.Router = mux.NewRouter()
	app.InitRoutes()
	return app, nil
}

func RunApp(app *services.App, addr string) {
	err := http.ListenAndServe(addr, app.Router)
	if err != nil {
		log.Fatal(err)
	}
}
