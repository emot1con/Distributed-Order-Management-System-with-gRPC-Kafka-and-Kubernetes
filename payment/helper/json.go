package helper

import (
	"encoding/json"
	"net/http"
)

func ReadJSON(r http.Request, data any) error {
	return json.NewDecoder(r.Body).Decode(&data)
}

func WriteJSON(w http.ResponseWriter, data any, status int) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}

func ErrorJSON(w http.ResponseWriter, err error, status int) {
	WriteJSON(w, err, status)
}
