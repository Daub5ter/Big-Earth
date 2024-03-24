package viewer

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

// paramsViewer - структура для получения
// параметров в uri запросе.
type paramsViewer struct{}

// NewParamsViewer - конструктор paramsViewer
func NewParamsViewer() Viewer {
	return paramsViewer{}
}

// ViewParam получает парметр из uri запроса.
func (cp paramsViewer) ViewParam(r *http.Request, param string) string {
	return chi.URLParam(r, param)
}
