package handlers

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"parsing-service/internal/handlers"
	"parsing-service/internal/models"
	parsinggrpc "parsing-service/internal/tools/grpc/parsing"

	"github.com/brianvoe/gofakeit/v7"
)

// TestHandler - основная функция для теста обработчиков.
func TestHandler(t *testing.T) {
	err := parseHandler()
	if err != nil {
		t.Error(err)
		return
	}
}

// parseHandler - запуск функций обработчика parse.
func parseHandler() error {
	ctx := context.Background()
	var db databaseMock
	var parse parsingMock

	parser := handlers.NewParser(db, parse)

	fmt.Println("отправляем неккоректные данные...")
	_, err := parser.Parse(ctx, &parsinggrpc.Place{})
	if err != nil {
		if !errors.Is(err, models.ErrEmptyData) {
			return fmt.Errorf("ошибка при отправке нккоректного запроса: %v", err)
		}
	}

	fmt.Println("отправляем корректные данные...")
	_, err = parser.Parse(ctx,
		&parsinggrpc.Place{
			Country: gofakeit.Country(),
			City:    gofakeit.City(),
		})
	if err != nil {
		return fmt.Errorf("ошибка при отправке корректных данных: %v", err)
	}

	fmt.Println("проверка на получение ошибки от событий места...")
	parse.returnErr = true
	_, err = parser.Parse(ctx,
		&parsinggrpc.Place{
			Country: gofakeit.Country(),
			City:    gofakeit.City(),
		})
	if err != nil {
		if !errors.Is(err, someErr) {
			return fmt.Errorf("ошибка при получении ошибки; моя ошибка: %v, полученная: %v", someErr, err)
		}
	}

	return nil
}

// databaseMock - мок базы данных.
type databaseMock struct{}

// GetPlaceInformation - мок функция для получения информации о месте.
func (d databaseMock) GetPlaceInformation(place *models.Place) (*models.PlaceInformation, error) {
	return &models.PlaceInformation{
		Text:   gofakeit.Word(),
		Photos: []string{gofakeit.URL()},
		Videos: []string{gofakeit.URL()},
	}, nil
}

// GetEventsLink - мок функция для получения ссылки на события.
func (d databaseMock) GetEventsLink(place *models.Place) (string, error) {
	return gofakeit.URL(), nil
}

// someErr - некоторая ошибка для теста.
var someErr = gofakeit.Error()

// parsingMock - мок парсинга.
type parsingMock struct {
	returnErr bool
}

// ParseEvents - мок функция для парсинга события
func (p parsingMock) ParseEvents(place *models.Place, link string) ([]*models.Event, error) {
	if p.returnErr {
		return nil, someErr
	}

	return []*models.Event{
		{
			Name:  gofakeit.Name(),
			Image: gofakeit.URL(),
			Link:  gofakeit.URL(),
		}}, nil
}
