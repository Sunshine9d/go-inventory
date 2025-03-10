package services

import (
	"fmt"
	db "github.com/Sunshine9d/go-inventory/db/mysql"
	"net/http"
)

func (app *App) getProducts(writer http.ResponseWriter, request *http.Request) {
	limit := 10
	offset := 0
	name := ""
	MysqlDB := db.MysqlDB{DB: app.DB}
	products, err := MysqlDB.GetProducts(limit, offset, name)
	if err != nil {
		sendError(writer, http.StatusInternalServerError, err.Error())
		return
	}

	sendResponse(writer, http.StatusOK, products)
}

func (app *App) createProduct(writer http.ResponseWriter, request *http.Request) {
	fmt.Println(writer, "not implemented yet")
}

func (app *App) getProduct(writer http.ResponseWriter, request *http.Request) {
	fmt.Println(writer, "not implemented yet")
}

func (app *App) updateProduct(writer http.ResponseWriter, request *http.Request) {
	fmt.Println(writer, "not implemented yet")
}

func (app *App) deleteProduct(writer http.ResponseWriter, request *http.Request) {
	fmt.Println(writer, "not implemented yet")
}

func (app *App) getHome(writer http.ResponseWriter, request *http.Request) {
	sendResponse(writer, http.StatusOK, map[string]string{"message": "Welcome to the Inventory Management System"})
}
