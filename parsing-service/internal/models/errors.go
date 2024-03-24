package models

import "errors"

// ErrNotFound ErrBadRequest
// - это ошибки, которые обрабатываются в функциях.
var (
	ErrNotFound   = errors.New("ничего не найдено")
	ErrBadRequest = errors.New("плохо сформулирован запрос")
)
