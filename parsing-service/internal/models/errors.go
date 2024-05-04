package models

import "errors"

// ErrNotFound, ErrEmptyData, ErrServer
// - это ошибки, которые обрабатываются в функциях.
var (
	ErrNotFound  = errors.New("ничего не найдено")
	ErrEmptyData = errors.New("некоторые данные из запроса постые")
	ErrServer    = errors.New("ошибка на стороне сервера")
)
