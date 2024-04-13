package handlers

import (
	"context"
	"errors"
	"fmt"
	log "log/slog"

	"parsing-service/internal/models"
	parsinggrpc "parsing-service/internal/tools/grpc/parsing"
	"parsing-service/internal/tools/parsing"
	"parsing-service/internal/tools/postgres"
)

// события - kassir.ru
// фото - wikipedia.org
// видео - youtube.com
// текст - tripadvisor.ru

// Parser - структура для работы с другими микросервисами.
type Parser struct {
	parsinggrpc.UnimplementedParsingServer
	database postgres.Database
	parsing  parsing.Parsing
}

// NewParser создает парсера.
func NewParser(database postgres.Database, parsing parsing.Parsing) *Parser {
	return &Parser{
		database: database,
		parsing:  parsing,
	}
}

// Parse получает данные из парсера и базы данных.
func (p *Parser) Parse(ctx context.Context, place *parsinggrpc.Place) (*parsinggrpc.PlaceInformation, error) {
	if place.Country == "" || place.City == "" {
		return nil, models.ErrEmptyData
	}

	pl := models.Place{
		Country: place.Country,
		City:    place.City,
	}

	eventsLink, err := p.database.GetEventsLink(&pl)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return nil, models.ErrNotFound
		}
		if errors.Is(err, models.ErrEmptyData) {
			return nil, models.ErrEmptyData
		}

		log.Error(fmt.Sprintf("ошибка на стороне сервера: %v", err))
		return nil, models.ErrServer
	}

	type eventsWithError struct {
		events []*models.Event
		err    error
	}

	events := make(chan eventsWithError)

	go func() {
		defer close(events)

		e, errParseEvent := p.parsing.ParseEvents(&pl, eventsLink)
		log.Debug("e:", e, "errParseEvent", errParseEvent)
		if errParseEvent != nil {
			events <- eventsWithError{events: nil, err: errParseEvent}
			return
		}

		events <- eventsWithError{events: e, err: nil}
	}()

	placeInformationWithoutEvents, err := p.database.GetPlaceInformation(&pl)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return nil, models.ErrNotFound
		}
		if errors.Is(err, models.ErrEmptyData) {
			return nil, models.ErrEmptyData
		}

		log.Error(fmt.Sprintf("ошибка на стороне сервера: %v", err))
		return nil, models.ErrServer
	}

	for {
		select {
		case <-ctx.Done():
			log.Warn("отмена контекста: перестаем ждать ответа от events")
			return nil, nil
		case e := <-events:
			var evs []*parsinggrpc.Event
			placeInformation := parsinggrpc.PlaceInformation{
				Text:   placeInformationWithoutEvents.Text,
				Photos: placeInformationWithoutEvents.Photos,
				Videos: placeInformationWithoutEvents.Videos,
				Events: evs,
			}

			if e.err != nil {
				log.Warn(fmt.Sprintf("не пришли события места: %v", e.err))
				return &placeInformation, nil
			}

			for _, event := range e.events {
				evs = append(evs,
					&parsinggrpc.Event{
						Name:  event.Name,
						Image: event.Image,
						Link:  event.Link,
					},
				)
			}

			placeInformation.Events = evs

			return &placeInformation, nil
		}
	}
}
