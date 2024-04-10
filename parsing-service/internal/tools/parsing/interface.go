package parsing

import "parsing-service/internal/models"

// Parsing - API для работы с парсером.
type Parsing interface {
	ParseEvent(place *models.Place, link string) ([]*models.Event, error)
}
