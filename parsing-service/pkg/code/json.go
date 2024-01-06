package code

import (
	"encoding/json"
	"log"
	"net/http"
)

// code JSON it`s json help structs & functions to write, read and error json.

// JSONResponse is JSON struct.
type JSONResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

// ReadJSON is decode json to struct.
func ReadJSON(r *http.Request, data any) error {
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(data)
	if err != nil {
		return err
	}

	return nil
}

// WriteJSON write struct to json and send it.
func WriteJSON(w http.ResponseWriter, status int, data any) {
	payload, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		log.Println(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_, err = w.Write(payload)
	if err != nil {
		log.Println(err)
		return
	}
}

func ErrorJSON(w http.ResponseWriter, status int, err error) {
	jsonResponse := JSONResponse{
		Error:   true,
		Message: err.Error(),
	}

	WriteJSON(w, status, jsonResponse)
}
