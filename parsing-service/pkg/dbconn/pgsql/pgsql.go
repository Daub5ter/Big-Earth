package pgsql

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

// openDB открывает соединение с pgsql.
func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

// ConnectToDB подключается к pgsql.
func ConnectToDB(dsn string) *sql.DB {
	var counts int

	for {
		connection, err := openDB(dsn)
		if err != nil {
			log.Println("попытка подключение к Postgres...")
			counts++
		} else {
			log.Println("подключено к Postgres!")
			return connection
		}

		if counts > 10 {
			log.Println(err)
			return nil
		}

		log.Println("ожидание 2 секунды...")
		time.Sleep(2 * time.Second)
		continue
	}
}
