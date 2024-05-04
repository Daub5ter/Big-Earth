package postgres

import "parsing-service/internal/models"

// Database - это абстракция базы данных.
type Database interface {
	GetPlaceInformation(place *models.Place) (*models.PlaceInformation, error)
	GetEventsLink(place *models.Place) (string, error)
}
