package parsing

import "parsing-service/internal/models"

// Parsing - API для работы с парсером.
type Parsing interface {
	ParseEvents(place *models.Place, link string) ([]*models.Event, error)
}
