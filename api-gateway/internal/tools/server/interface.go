package server

import (
	"net/http"
)

// Server - API сервера, его обработчиков.
type Server interface {
	Routes() http.Handler
}
