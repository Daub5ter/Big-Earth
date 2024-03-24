package viewer

import "net/http"

// Viewer - API для работы с получение параметров
// в ссылке запроса типа /{param}
type Viewer interface {
	ViewParam(r *http.Request, param string) string
}
