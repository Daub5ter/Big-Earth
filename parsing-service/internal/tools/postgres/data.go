package postgres

import (
	"context"
	"database/sql"
	"errors"

	"parsing-service/internal/models"
)

// GetPlaceInformation получает данные по месту.
func (db db) GetPlaceInformation(place *models.Place) (*models.PlaceInformation, error) {
	if place.Country == "" || place.City == "" {
		return nil, models.ErrEmptyData
	}

	ctx, cancel := context.WithTimeout(context.Background(), db.timeout)
	defer cancel()

	var placeInformation models.PlaceInformation
	var placeID uint64

	row := db.conn.QueryRowContext(ctx, selectPlaceID, place.Country, place.City)

	err := row.Scan(&placeID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNotFound
		}
		return nil, err
	}

	row = db.conn.QueryRowContext(ctx, selectText, placeID)

	err = row.Scan(&placeInformation.Text)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNotFound
		}
		return nil, err
	}

	rows, err := db.conn.QueryContext(ctx, selectPhotos, placeID)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var photos []string
	var tempPhoto string

	for rows.Next() {
		err = rows.Scan(
			&tempPhoto,
		)
		if err != nil {
			return nil, err
		}

		photos = append(photos, tempPhoto)
	}

	placeInformation.Photos = photos

	rows, err = db.conn.QueryContext(ctx, selectVideos, placeID)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var videos []string
	var tempVideo string

	for rows.Next() {

		err = rows.Scan(
			&tempVideo,
		)
		if err != nil {
			return nil, err
		}

		videos = append(videos, tempVideo)
	}

	placeInformation.Videos = videos

	return &placeInformation, nil
}

// GetEventsLink получает ссылку на события по месту.
func (db db) GetEventsLink(place *models.Place) (string, error) {
	if place.Country == "" || place.City == "" {
		return "", models.ErrEmptyData
	}

	ctx, cancel := context.WithTimeout(context.Background(), db.timeout)
	defer cancel()

	var placeID uint64
	row := db.conn.QueryRowContext(ctx, selectPlaceID, place.Country, place.City)

	err := row.Scan(&placeID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", models.ErrNotFound
		}
		return "", err
	}

	var link string
	row = db.conn.QueryRowContext(ctx, selectEventsLink, placeID)

	err = row.Scan(&link)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", models.ErrNotFound
		}
		return "", err
	}

	return link, nil
}
