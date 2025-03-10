package services

import (
	"database/sql"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

type App struct {
	Router *mux.Router
	DB     *sql.DB
}

func sendResponse(writer http.ResponseWriter, status int, data interface{}) {
	response, _ := json.Marshal(data)
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(status)
	writer.Write(response)
}

func sendError(writer http.ResponseWriter, status int, message string) {
	errorMessage := map[string]string{"error": message}
	sendResponse(writer, status, errorMessage)
}
