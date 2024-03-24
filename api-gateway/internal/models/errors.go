package models

import "errors"

// ErrBadRequest - это ошибка, которая обрабатываются в функциях.
var (
	ErrBadRequest = errors.New("плохо сформулирован запрос")
)
