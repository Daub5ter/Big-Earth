package postgres

// createPlace, createTexts, createPhotos, createVideos
// - команды создания таблиц базы данных, если они отсутсвуют.
const (
	createPlace = `CREATE TABLE IF NOT EXISTS place
					(
						id serial PRIMARY KEY,
						country VARCHAR(255),
						city VARCHAR(255)
					)`

	createTexts = `CREATE TABLE IF NOT EXISTS texts
					(
						place_id INT REFERENCES place(id),
						text_place TEXT
					)`

	createPhotos = `CREATE TABLE IF NOT EXISTS photos
						(
							place_id INT REFERENCES place(id),
							link VARCHAR(550)
						)`

	createVideos = `CREATE TABLE IF NOT EXISTS videos
						(
							place_id INT REFERENCES place(id),
							link VARCHAR(255)
						);`

	createEvents = `CREATE TABLE IF NOT EXISTS events
						(
							place_id INT REFERENCES place(id),
							link VARCHAR(255)
						);`
)

// indexTexts, indexPhotos, indexVideos
// - команды для создания индексов в таблице.
const (
	indexTexts = `CREATE INDEX IF NOT EXISTS idx_texts 
					ON texts 
					USING HASH (place_id)`

	indexPhotos = `CREATE INDEX IF NOT EXISTS idx_photos 
					ON photos 
					USING HASH (place_id)`

	indexVideos = `CREATE INDEX IF NOT EXISTS idx_videos 
					ON videos 
					USING HASH (place_id)`

	indexEvents = `CREATE INDEX IF NOT EXISTS idx_events 
					ON events 
					USING HASH (place_id)`
)

// selectPlaceID, selectText, selectPhotos, selectVideos
// - запросы на получение данных из базы данных.
const (
	selectPlaceID = `SELECT id
						FROM place
						WHERE country = $1 AND city = $2`

	selectText = `SELECT text_place
					FROM texts
					WHERE place_id = $1`

	selectPhotos = `SELECT link
						FROM photos
						WHERE place_id = $1`

	selectVideos = `SELECT link
						FROM videos
						WHERE place_id = $1`

	selectEventsLink = `SELECT link
							FROM events
							WHERE place_id = $1`
)
