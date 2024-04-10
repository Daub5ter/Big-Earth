package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	log "log/slog"
	"testing"

	"parsing-service/internal/models"
	"parsing-service/internal/tools/config"
	"parsing-service/internal/tools/postgres"
	"parsing-service/pkg/dbconn/pgsql"
	"parsing-service/pkg/logger"

	"github.com/brianvoe/gofakeit/v7"
	tc "github.com/testcontainers/testcontainers-go/modules/compose"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestPostgres(t *testing.T) {
	logger.SetLogger("debug")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	log.Debug("запуск контейнеров...")
	var cs containers
	compose, err := cs.runCompose(ctx)
	if err != nil {
		t.Errorf("ошибка запуска контейнеров: %v", err)
		return
	}

	defer func() {
		log.Debug("завершение работы контейеров...")
		err = cs.downCompose(context.Background(), compose)
		if err != nil {
			t.Errorf("ошибка завершения работы контейнеров: %v", err)
			return
		}
	}()

	log.Debug("настройка конфигурации...")
	t.Setenv("DSN", "host=localhost port=5432 user=some_user password=some_password dbname=some_db sslmode=disable timezone=UTC connect_timeout=5")
	cfg, err := config.NewConfig("postgres_test_config.yaml")
	if err != nil {
		t.Errorf("ошибка прочтения файла конфигруаций: %v", err)
		return
	}

	log.Debug("подключение к базе данных postgresql...")
	conn := pgsql.ConnectToDB(cfg.(config.DatabaseConfig).DSN())
	if conn == nil {
		t.Errorf("ошибка подключения к Postgres")
		return
	}

	pgdb, err := postgres.NewDB(conn, cfg)
	if err != nil {
		t.Errorf("ошибка подключения к базе данных: %v", err)
		return
	}

	db := database{
		ctx:  ctx,
		conn: conn,
	}

	log.Debug("сохранение корректных данных...")
	err = db.saveCorrectData()
	if err != nil {
		t.Errorf("ошибка сохранения корректных данных: %v", err)
		return
	}

	log.Debug("проверка корректных данных...")
	err = db.readCorrectData(pgdb)
	if err != nil {
		t.Errorf("ошибка проверки корректных данных: %v", err)
		return
	}

	log.Debug("проверка на обработку неккоректных данных...")
	err = db.checkUnCorrectData(pgdb)
	if err != nil {
		t.Errorf("ошибка проверки на обработку неккоретных данных: %v", err)
		return
	}

	log.Debug("удаление данных из таблиц...")
	err = db.deleteAllData()
	if err != nil {
		t.Errorf("ошибка удаления данных из таблиц: %v", err)
		return
	}

	log.Debug("сохранение большого количества данных...")
	err = db.saveManyData()
	if err != nil {
		t.Errorf("ошибка сохранения большого количества данных: %v", err)
		return
	}

	log.Debug("проверка корректности большого количества данных...")
	err = db.countManyData()
	if err != nil {
		t.Errorf("ошибка прооверки большого количества данных: %v", err)
		return
	}

	err = pgsql.CloseConnection(conn)
	if err != nil {
		t.Errorf("ошибка при завершении работы базы данных %v", err)
		return
	}
}

// containers - структура для работы с запуском и завершенем контейнера.
type containers struct{}

// runCompose запускает docker compose clickhouse
func (cs containers) runCompose(ctx context.Context) (tc.ComposeStack, error) {
	compose, err := tc.NewDockerComposeWith(tc.WithStackFiles("postgres_test.yml"),
		tc.StackIdentifier("postgres_test_identifier"))
	if err != nil {
		return nil, err
	}

	err = compose.Up(ctx, tc.Wait(true))
	if err != nil {
		return nil, err
	}

	compose.WaitForService("clickhouse_test", wait.ForLog("Awaiting socket connections on 0.0.0.0:5432"))

	return compose, nil
}

// downCompose выключает docker compose и удаляет образы.
func (cs containers) downCompose(ctx context.Context, compose tc.ComposeStack) error {
	return compose.Down(ctx, tc.RemoveImagesAll, tc.RemoveVolumes(true))
}

// database - структура для работы с базой данных.
type database struct {
	ctx  context.Context
	conn *sql.DB
}

// saveCorrectData сохраняет корректные данные.
func (db database) saveCorrectData() error {
	var placeID int

	query := `insert into place(
                country,
                city) 
				values ($1, $2)
				returning id`
	err := db.conn.QueryRowContext(db.ctx, query,
		country,
		city,
	).Scan(&placeID)
	if err != nil {
		return err
	}

	query = `insert into texts(
        		place_id,
                text_place) 
				values ($1, $2)`
	_, err = db.conn.ExecContext(db.ctx, query,
		placeID,
		text)
	if err != nil {
		return err
	}

	query = `insert into photos(
        		place_id,
                link) 
				values ($1, $2)`
	_, err = db.conn.ExecContext(db.ctx, query,
		placeID,
		photoLink)
	if err != nil {
		return err
	}

	query = `insert into videos(
        		place_id,
                link) 
				values ($1, $2)`
	_, err = db.conn.ExecContext(db.ctx, query,
		placeID,
		videoLink)
	if err != nil {
		return err
	}

	query = `insert into events(
        		place_id,
                link) 
				values ($1, $2)`
	_, err = db.conn.ExecContext(db.ctx, query,
		placeID,
		eventsLink)
	if err != nil {
		return err
	}

	return nil
}

// readCorrectData читает корректные данные.
func (db database) readCorrectData(pgdb postgres.Database) error {
	place := models.Place{
		Country: country,
		City:    city,
	}

	placeInformation, err := pgdb.GetPlaceInformation(&place)

	if err != nil {
		return err
	}

	if text != placeInformation.Text {
		return fmt.Errorf("неверные данные таблицы texts получены из базы данных; ожидалось: %s, получено: %s",
			text, placeInformation.Text)
	}

	if len(placeInformation.Photos) == 0 {
		return fmt.Errorf("нет данных, ожидалась минимум 1 из таблицы photos")
	}
	if photoLink != placeInformation.Photos[0] {
		return fmt.Errorf("неверные данные таблицы photos получены из базы данных; ожидалось: %s, получено: %s",
			photoLink, placeInformation.Photos[0])
	}

	if len(placeInformation.Videos) == 0 {
		return fmt.Errorf("нет данных, ожидалась минимум 1 из таблицы videos")
	}
	if videoLink != placeInformation.Videos[0] {
		return fmt.Errorf("неверные данные таблицы videos получены из базы данных; ожидалось: %s, получено: %s",
			videoLink, placeInformation.Videos[0])
	}

	eventsLinkFromDB, err := pgdb.GetEventsLink(&place)
	if err != nil {
		return err
	}

	if eventsLink != eventsLinkFromDB {
		return fmt.Errorf("неверные данные таблицы events получены из базы данных; ожидалось: %s, получено: %s",
			eventsLink, eventsLinkFromDB)
	}

	return nil
}

// checkUnCorrectData обрабатывает некорректные данные (проверяет на ошибки).
func (db database) checkUnCorrectData(pgdb postgres.Database) error {
	place := models.Place{}

	if _, err := pgdb.GetPlaceInformation(&place); err != nil {
		if !errors.Is(err, models.ErrEmptyData) {
			return err
		}
	}

	if _, err := pgdb.GetEventsLink(&place); err != nil {
		if !errors.Is(err, models.ErrEmptyData) {
			return err
		}
	}

	return nil
}

// deleteAllData удаляет все данные из таблиц.
func (db database) deleteAllData() error {
	if _, err := db.conn.ExecContext(db.ctx, "TRUNCATE TABLE texts CASCADE"); err != nil {
		return err
	}

	if _, err := db.conn.ExecContext(db.ctx, "TRUNCATE TABLE photos CASCADE"); err != nil {
		return err
	}

	if _, err := db.conn.ExecContext(db.ctx, "TRUNCATE TABLE videos CASCADE"); err != nil {
		return err
	}

	if _, err := db.conn.ExecContext(db.ctx, "TRUNCATE TABLE events CASCADE"); err != nil {
		return err
	}

	if _, err := db.conn.ExecContext(db.ctx, "TRUNCATE TABLE place CASCADE"); err != nil {
		return err
	}

	return nil
}

var count = 5000

// saveManyData сохраняет много значений.
func (db database) saveManyData() error {
	for i := 1; i <= count; i++ {
		var placeID int

		query := `insert into place(
                country,
                city) 
				values ($1, $2)
				returning id`
		err := db.conn.QueryRowContext(db.ctx, query,
			gofakeit.Country(),
			gofakeit.City(),
		).Scan(&placeID)
		if err != nil {
			return err
		}

		query = `insert into texts(
        		place_id,
                text_place) 
				values ($1, $2)`
		_, err = db.conn.ExecContext(db.ctx, query,
			placeID,
			gofakeit.Word())
		if err != nil {
			return err
		}

		query = `insert into photos(
        		place_id,
                link) 
				values ($1, $2)`
		_, err = db.conn.ExecContext(db.ctx, query,
			placeID,
			gofakeit.URL())
		if err != nil {
			return err
		}

		query = `insert into videos(
        		place_id,
                link) 
				values ($1, $2)`
		_, err = db.conn.ExecContext(db.ctx, query,
			placeID,
			gofakeit.URL())
		if err != nil {
			return err
		}

		query = `insert into events(
        		place_id,
                link) 
				values ($1, $2)`
		_, err = db.conn.ExecContext(db.ctx, query,
			placeID,
			gofakeit.URL())
		if err != nil {
			return err
		}
	}

	return nil
}

// countManyData считает много значений; если не совпадает, выдает ошикую
func (db database) countManyData() error {
	var countFromDB int

	row := db.conn.QueryRowContext(db.ctx, "SELECT COUNT(*) FROM place")
	err := row.Scan(&countFromDB)
	if err != nil {
		return err
	}
	if count != countFromDB {
		return fmt.Errorf("пришло неверное количество данных таблицы place, требовалось: %v, получили: %v",
			count, countFromDB)
	}

	row = db.conn.QueryRowContext(db.ctx, "SELECT COUNT(*) FROM texts")
	err = row.Scan(&countFromDB)
	if err != nil {
		return err
	}
	if count != countFromDB {
		return fmt.Errorf("пришло неверное количество данных таблицы texts, требовалось: %v, получили: %v",
			count, countFromDB)
	}

	row = db.conn.QueryRowContext(db.ctx, "SELECT COUNT(*) FROM photos")
	err = row.Scan(&countFromDB)
	if err != nil {
		return err
	}
	if count != countFromDB {
		return fmt.Errorf("пришло неверное количество данных таблицы photos, требовалось: %v, получили: %v",
			count, countFromDB)
	}

	row = db.conn.QueryRowContext(db.ctx, "SELECT COUNT(*) FROM videos")
	err = row.Scan(&countFromDB)
	if err != nil {
		return err
	}
	if count != countFromDB {
		return fmt.Errorf("пришло неверное количество данных таблицы videos, требовалось: %v, получили: %v",
			count, countFromDB)
	}

	row = db.conn.QueryRowContext(db.ctx, "SELECT COUNT(*) FROM events")
	err = row.Scan(&countFromDB)
	if err != nil {
		return err
	}
	if count != countFromDB {
		return fmt.Errorf("пришло неверное количество данных таблицы events, требовалось: %v, получили: %v",
			count, countFromDB)
	}

	return nil
}

// country, city, text, photoLink, videoLink, eventsLink
// - данные для тестирования корректности базы данных.
var (
	country = "some country"
	city    = "some city"

	text = "some text"

	photoLink = "some photo link"

	videoLink = "some video link"

	eventsLink = "some events link"
)
