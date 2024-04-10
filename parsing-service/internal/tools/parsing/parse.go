package parsing

import (
	"errors"

	"parsing-service/internal/models"

	"github.com/gocolly/colly"
)

// Parser - структура парсера.
type parser struct {
	collector *colly.Collector
}

// NewParsing - конструктор парсера.
func NewParsing() Parsing {
	collector := colly.NewCollector()

	collector.AllowURLRevisit = true
	collector.OnRequest(func(r *colly.Request) {
		r.Headers.Set("User-Agent", "Mozilla/5.0 (Windows NT 6.1; Win64; x64; rv:47.0) Gecko/20100101 Firefox/47.0")
	})

	return parser{
		collector: collector,
	}
}

// allocationEvent - количество выделяемого места для среза событий.
var allocationEvent = 20

// ParseEvent - парсит события.
func (p parser) ParseEvent(place *models.Place, link string) ([]*models.Event, error) {
	p.collector.AllowURLRevisit = true
	p.collector.OnRequest(func(r *colly.Request) {
		r.Headers.Set("User-Agent", "Mozilla/5.0 (Windows NT 6.1; Win64; x64; rv:47.0) Gecko/20100101 Firefox/47.0")
	})

	switch place.Country {
	case "Russia":
		events := make([]*models.Event, 0, allocationEvent)
		err := p.parseRussia(link, &events)
		if err != nil {
			return nil, err
		}

		return events, nil
	default:
		return nil, errors.New("country not found")
	}
}
