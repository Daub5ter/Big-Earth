// Пакет для записи и чтения данных формата JSON.

package code

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

// JSONResponse структура формата JSON.
type JSONResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

// ReadJSON декодирует JSON в любой тип данных.
// Данные следует передавать в виде адреса памяти (&database-data).
func ReadJSON(bodyR io.ReadCloser, data any) error {
	// Считывание тела запроса.
	body, err := io.ReadAll(bodyR)
	if err != nil {
		return err
	}

	// Декодирование JSON в данные.
	err = json.Unmarshal(body, data)
	if err != nil {
		return err
	}

	return nil
}

// WriteJSON кодирует любой тип данных в JSON и пишет его в ResponseWriter.
func WriteJSON(w http.ResponseWriter, status int, data any) {
	// Кодирование данных в JSON.
	payload, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		log.Println(err)
		return
	}

	// Запись данных.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_, err = w.Write(payload)
	if err != nil {
		log.Println(err)
		return
	}
}

// ErrorJSON оборачивает ошибку в JSON и пишет ее в ResponseWriter.
func ErrorJSON(w http.ResponseWriter, status int, err error) {
	jsonResponse := JSONResponse{
		Error:   true,
		Message: err.Error(),
	}

	WriteJSON(w, status, jsonResponse)
}
