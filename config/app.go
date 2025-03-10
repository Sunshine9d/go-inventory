package config

import (
	"database/sql"
	"encoding/json"
	"fmt"
	db "github.com/Sunshine9d/go-inventory/db/mysql"
	sv "github.com/Sunshine9d/go-inventory/services"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func (app *sv.App) Initialize(user, password, dbname string) error {
	connectionString := fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, password, dbname)
	var err error
	app.DB, err = sql.Open("mysql", connectionString)
	if err != nil {
		log.Fatal(err)
	}

	app.Router = mux.NewRouter()
	return nil
}

func (app *App) Run(addr string) {
	app.initializeRoutes()
	err := http.ListenAndServe(addr, app.Router)
	if err != nil {
		log.Fatal(err)
	}
}

//func sendResponse(writer http.ResponseWriter, status int, data interface{}) {
//	response, _ := json.Marshal(data)
//	writer.Header().Set("Content-Type", "application/json")
//	writer.WriteHeader(status)
//	writer.Write(response)
//}
//
//func sendError(writer http.ResponseWriter, status int, message string) {
//	errorMessage := map[string]string{"error": message}
//	sendResponse(writer, status, errorMessage)
//}

func (app *App) initializeRoutes() {
	app.Router.HandleFunc("/", app.getHome).Methods("GET")
	app.Router.HandleFunc("/products", app.getProducts).Methods("GET")
	app.Router.HandleFunc("/product", app.createProduct).Methods("POST")
	app.Router.HandleFunc("/product/{id:[0-9]+}", app.getProduct).Methods("GET")
	app.Router.HandleFunc("/product/{id:[0-9]+}", app.updateProduct).Methods("PUT")
	app.Router.HandleFunc("/product/{id:[0-9]+}", app.deleteProduct).Methods("DELETE")
}

//func (app *App) getProducts(writer http.ResponseWriter, request *http.Request) {
//	limit := 10
//	offset := 0
//	name := ""
//	MysqlDB := db.MysqlDB{DB: app.DB}
//	products, err := MysqlDB.GetProducts(limit, offset, name)
//	if err != nil {
//		sendError(writer, http.StatusInternalServerError, err.Error())
//		return
//	}
//
//	sendResponse(writer, http.StatusOK, products)
//}
//
//func (app *App) createProduct(writer http.ResponseWriter, request *http.Request) {
//	fmt.Println(writer, "not implemented yet")
//}
//
//func (app *App) getProduct(writer http.ResponseWriter, request *http.Request) {
//	fmt.Println(writer, "not implemented yet")
//}
//
//func (app *App) updateProduct(writer http.ResponseWriter, request *http.Request) {
//	fmt.Println(writer, "not implemented yet")
//}
//
//func (app *App) deleteProduct(writer http.ResponseWriter, request *http.Request) {
//	fmt.Println(writer, "not implemented yet")
//}
//
//func (app *App) getHome(writer http.ResponseWriter, request *http.Request) {
//	sendResponse(writer, http.StatusOK, map[string]string{"message": "Welcome to the Inventory Management System"})
//}
