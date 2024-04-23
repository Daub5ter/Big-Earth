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
	"github.com/stretchr/testify/suite"
)

// ClientSuite - структура для тестов.
type ClientSuite struct {
	suite.Suite
	ctx context.Context
}

// SetupSuite настраивает тесты
// (включается перед тестами) .
func (c *ClientSuite) SetupSuite() {
	c.ctx = context.Background()
}

// TestHandlers запускает тесты.
func TestHandlers(t *testing.T) {
	suite.Run(t, new(ClientSuite))
}

// parseHandler - запуск функций обработчика parse.
func (c *ClientSuite) TestParseHandler() {
	var db databaseMock
	var parse parsingMock

	parser := handlers.NewParser(db, parse)

	_, err := parser.Parse(c.ctx, &parsinggrpc.Place{})
	if err != nil {
		if !errors.Is(err, models.ErrEmptyData) {
			c.NoError(fmt.Errorf("ошибка при отправке нккоректного запроса: %v", err))
		}
	}

	_, err = parser.Parse(c.ctx,
		&parsinggrpc.Place{
			Country: gofakeit.Country(),
			City:    gofakeit.City(),
		})
	c.NoError(err)

	parse.returnErr = true
	_, err = parser.Parse(c.ctx,
		&parsinggrpc.Place{
			Country: gofakeit.Country(),
			City:    gofakeit.City(),
		})
	if err != nil {
		if !errors.Is(err, someErr) {
			c.NoError(fmt.Errorf("ошибка при получении ошибки; моя ошибка: %v, полученная: %v", someErr, err))
		}
	}

	c.T().Log("обработчик parse работает")
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
