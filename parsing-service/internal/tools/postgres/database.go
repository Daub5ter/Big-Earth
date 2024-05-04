package postgres

import (
	"database/sql"
	log "log/slog"
	"time"

	"parsing-service/internal/tools/config"
)

// DB - это представление БД.
type db struct {
	conn    *sql.DB
	timeout time.Duration
}

// NewDB создает новую структуру DB.
func NewDB(conn *sql.DB, cfg config.DatabaseConfig) (Database, error) {
	dbase := db{
		conn:    conn,
		timeout: cfg.DBTimeout(),
	}

	conn.SetMaxOpenConns(100)
	conn.SetMaxIdleConns(10)
	conn.SetConnMaxLifetime(time.Hour)

	if err := dbase.checkConnection(); err != nil {
		return nil, err
	}

	if err := dbase.createTables(); err != nil {
		return nil, err
	}

	if err := dbase.createIndexes(); err != nil {
		return nil, err
	}

	return dbase, nil
}

// checkConnection проверяет соединение с базой данных.
func (db db) checkConnection() error {
	err := db.conn.Ping()
	if err != nil {
		return err
	}
	log.Info("подключено к базе данных")

	return nil
}

// createTables создает таблицы, если они не созданы.
func (db db) createTables() error {
	if _, err := db.conn.Exec(createPlace); err != nil {
		return err
	}

	if _, err := db.conn.Exec(createPhotos); err != nil {
		return err
	}

	if _, err := db.conn.Exec(createTexts); err != nil {
		return err
	}

	if _, err := db.conn.Exec(createVideos); err != nil {
		return err
	}

	if _, err := db.conn.Exec(createEvents); err != nil {
		return err
	}

	return nil
}

// createIndexes создает индексы на требуемых таблицах,
// если они еще не созданы.
func (db db) createIndexes() error {
	if _, err := db.conn.Exec(indexTexts); err != nil {
		return err
	}

	if _, err := db.conn.Exec(indexPhotos); err != nil {
		return err
	}

	if _, err := db.conn.Exec(indexVideos); err != nil {
		return err
	}

	if _, err := db.conn.Exec(indexEvents); err != nil {
		return err
	}

	return nil
}
