package parsing

import (
	"errors"

	"parsing-service/internal/models"
	"parsing-service/internal/tools/config"

	"github.com/gocolly/colly"
)

// Parser - структура парсера.
type parser struct {
	collector *colly.Collector
	cfg       config.EventsURIs
}

// NewParsing - конструктор парсера.
func NewParsing(cfg config.EventsURIs) Parsing {
	return parser{
		collector: colly.NewCollector(),
		cfg:       cfg,
	}
}

// allocationEvent - количество выделяемого места для среза событий.
var allocationEvent = 20

// ParseEvent - парсит события.
func (p parser) ParseEvent(place *models.Place) ([]*models.Event, error) {
	p.collector.AllowURLRevisit = true
	p.collector.OnRequest(func(r *colly.Request) {
		r.Headers.Set("User-Agent", "Mozilla/5.0 (Windows NT 6.1; Win64; x64; rv:47.0) Gecko/20100101 Firefox/47.0")
	})

	switch place.Country {
	case "Russia":
		var urlPlace string

		switch place.City {
		case "Krasnodar":
			urlPlace = p.cfg.GetRussiaKrasnodar()
		default:
			return nil, models.ErrBadRequest
		}

		events := make([]*models.Event, 0, allocationEvent)
		err := p.parseEventRussia(urlPlace, &events)
		if err != nil {
			return nil, err
		}

		return events, nil
	default:
		return nil, errors.New("country not found")
	}
}
