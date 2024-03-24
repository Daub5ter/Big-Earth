package handlers

import (
	"context"
	"errors"

	"parsing-service/internal/models"
	"parsing-service/internal/tools/database"
	"parsing-service/internal/tools/parsing"
	"parsing-service/internal/tools/rpc/grpcparsing"
)

// kassir.ru
// wikipedia.org
// youtube.com
// tripadvisor.ru

// Parser - структура для работы с другими микросервисами.
type Parser struct {
	grpcparsing.UnsafeParsingServer
	database database.Database
	parsing  parsing.Parsing
}

// NewParser создает парсера.
func NewParser(database database.Database, parsing parsing.Parsing) *Parser {
	return &Parser{
		database: database,
		parsing:  parsing,
	}
}

// Parse получает данные из парсера и базы данных.
func (p *Parser) Parse(ctx context.Context, place *grpcparsing.Place) (*grpcparsing.PlaceInformation, error) {
	pl := models.Place{
		Country: place.Country,
		City:    place.City,
	}

	type eventsWithError struct {
		events []*models.Event
		err    error
	}

	events := make(chan eventsWithError)

	go func() {
		defer close(events)

		e, err := p.parsing.ParseEvent(&pl)
		if err != nil {
			events <- eventsWithError{events: nil, err: err}
			return
		}

		events <- eventsWithError{events: e, err: nil}
	}()

	data, err := p.database.GetPlaceInformation(&pl)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return nil, models.ErrNotFound
		}
		if errors.Is(err, models.ErrBadRequest) {
			return nil, models.ErrBadRequest
		}

		return nil, err
	}

	// TODO: использовать context для проверки на ошибку и возвращаться в таком случае
	e := <-events
	if e.err != nil {
		return nil, err
	}

	var evnts []*grpcparsing.Event
	for _, event := range e.events {
		evnts = append(evnts,
			&grpcparsing.Event{
				Name:  event.Name,
				Image: event.Image,
				Link:  event.Link,
			},
		)
	}

	return &grpcparsing.PlaceInformation{
		Text:   data.Text,
		Photos: data.Photos,
		Videos: data.Videos,
		Events: evnts,
	}, nil
}
