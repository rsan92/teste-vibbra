package http_responses

import (
	"encoding/json"
	"log"
	"net/http"
)

func JSON(w http.ResponseWriter, statusCode int, info interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if info != nil {
		if err := json.NewEncoder(w).Encode(info); err != nil {
			log.Fatal(err)
		}
	}
}

func Error(w http.ResponseWriter, statusCode int, err error) {
	JSON(w, statusCode, struct {
		Error string `json:"erro"`
	}{
		Error: err.Error(),
	})
}
