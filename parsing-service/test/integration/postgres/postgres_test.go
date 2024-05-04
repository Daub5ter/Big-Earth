package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"testing"

	"parsing-service/internal/models"
	"parsing-service/internal/tools/config"
	"parsing-service/internal/tools/postgres"
	"parsing-service/pkg/dbconn/pgsql"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/suite"
	tc "github.com/testcontainers/testcontainers-go/modules/compose"
	"github.com/testcontainers/testcontainers-go/wait"
)

// ClientSuite - структура для тестов.
type ClientSuite struct {
	suite.Suite
	compose  tc.ComposeStack
	ctx      context.Context
	cfg      config.DatabaseConfig
	conn     *sql.DB
	database postgres.Database
}

// SetupSuite настраивает тесты
// (включается перед тестами) .
func (c *ClientSuite) SetupSuite() {
	var err error
	c.compose, err = tc.NewDockerComposeWith(tc.WithStackFiles("postgres_test.yml"),
		tc.StackIdentifier("postgres_test_identifier"))
	c.NoError(err)

	c.ctx = context.Background()

	err = c.compose.Up(c.ctx, tc.Wait(true))
	c.NoError(err)

	c.compose.WaitForService("postgres_test", wait.ForLog("Awaiting socket connections on 0.0.0.0:5432"))

	c.Suite.T().Setenv("DSN", "host=localhost port=5432 user=some_user password=some_password dbname=some_db sslmode=disable timezone=UTC connect_timeout=5")

	c.cfg, err = config.NewConfig("postgres_test_config.yaml")
	c.NoError(err)

	c.conn = pgsql.ConnectToDB(c.cfg.DSN())
	if c.conn == nil {
		c.NoError(fmt.Errorf("ошибка подключения к Postgres"))
	}

	c.database, err = postgres.NewDB(c.conn, c.cfg)
	c.NoError(err)

	c.T().Log("Postgres запущен")
}

// TearDownSuite завершает тесты
// (включается после тестами) .
func (c *ClientSuite) TearDownSuite() {
	c.NoError(c.conn.Close())
	c.NoError(c.compose.Down(c.ctx, tc.RemoveImagesAll, tc.RemoveVolumes(true)))
}

// TestPostgres запускает тесты.
func TestPostgres(t *testing.T) {
	suite.Run(t, new(ClientSuite))
}

// TestA_CorrectData проверяет работу с корректными сообщениями.
func (c *ClientSuite) TestA_CorrectData() {
	var placeID int

	query := `insert into place(
                country,
                city) 
				values ($1, $2)
				returning id`
	err := c.conn.QueryRowContext(c.ctx, query,
		country,
		city,
	).Scan(&placeID)
	c.NoError(err)

	query = `insert into texts(
        		place_id,
                text_place) 
				values ($1, $2)`
	_, err = c.conn.ExecContext(c.ctx, query,
		placeID,
		text)
	c.NoError(err)

	query = `insert into photos(
        		place_id,
                link) 
				values ($1, $2)`
	_, err = c.conn.ExecContext(c.ctx, query,
		placeID,
		photoLink)
	c.NoError(err)

	query = `insert into videos(
        		place_id,
                link) 
				values ($1, $2)`
	_, err = c.conn.ExecContext(c.ctx, query,
		placeID,
		videoLink)
	c.NoError(err)

	query = `insert into events(
        		place_id,
                link) 
				values ($1, $2)`
	_, err = c.conn.ExecContext(c.ctx, query,
		placeID,
		eventsLink)
	c.NoError(err)

	place := models.Place{
		Country: country,
		City:    city,
	}

	placeInformation, err := c.database.GetPlaceInformation(&place)
	c.NoError(err)

	if !c.Equal(text, placeInformation.Text) {
		c.NoError(fmt.Errorf("неверные данные таблицы texts получены из базы данных; ожидалось: %s, получено: %s",
			text, placeInformation.Text))
	}

	if !c.Equal(photoLink, placeInformation.Photos[0]) {
		c.NoError(fmt.Errorf("неверные данные таблицы photos получены из базы данных; ожидалось: %s, получено: %s",
			photoLink, placeInformation.Photos[0]))
	}

	if !c.Equal(videoLink, placeInformation.Videos[0]) {
		c.NoError(fmt.Errorf("неверные данные таблицы videos получены из базы данных; ожидалось: %s, получено: %s",
			videoLink, placeInformation.Videos[0]))
	}

	if len(placeInformation.Photos) == 0 {
		c.NoError(fmt.Errorf("нет данных, ожидалась минимум 1 из таблицы photos"))
	}

	if len(placeInformation.Videos) == 0 {
		c.NoError(fmt.Errorf("нет данных, ожидалась минимум 1 из таблицы videos"))
	}

	eventsLinkFromDB, err := c.database.GetEventsLink(&place)
	c.NoError(err)

	if !c.Equal(eventsLink, eventsLinkFromDB) {
		c.NoError(fmt.Errorf("неверные данные таблицы events получены из базы данных; ожидалось: %s, получено: %s",
			eventsLink, eventsLinkFromDB))
	}

	c.T().Log("корректные данные записываются и читаются верно")

	_, err = c.conn.ExecContext(c.ctx, "TRUNCATE TABLE texts CASCADE")
	c.NoError(err)

	_, err = c.conn.ExecContext(c.ctx, "TRUNCATE TABLE photos CASCADE")
	c.NoError(err)

	_, err = c.conn.ExecContext(c.ctx, "TRUNCATE TABLE videos CASCADE")
	c.NoError(err)

	_, err = c.conn.ExecContext(c.ctx, "TRUNCATE TABLE events CASCADE")
	c.NoError(err)

	_, err = c.conn.ExecContext(c.ctx, "TRUNCATE TABLE place CASCADE")
	c.NoError(err)
}

// TestB_UnCorrectData проверяет работу с некорректными сообщениями.
func (c *ClientSuite) TestB_UnCorrectData() {
	place := models.Place{}

	if _, err := c.database.GetPlaceInformation(&place); err != nil {
		if !errors.Is(err, models.ErrEmptyData) {
			c.NoError(err)
		}
	}

	if _, err := c.database.GetEventsLink(&place); err != nil {
		if !errors.Is(err, models.ErrEmptyData) {
			c.NoError(err)
		}
	}

	c.T().Log("некорректные данные обрабатываются верно")
}

// TestC_ManyData проверяет работу с большим количеством сообщений.
func (c *ClientSuite) TestC_ManyData() {
	var count = 5000

	for i := 1; i <= count; i++ {
		var placeID int

		query := `insert into place(
                country,
                city) 
				values ($1, $2)
				returning id`
		err := c.conn.QueryRowContext(c.ctx, query,
			gofakeit.Country(),
			gofakeit.City(),
		).Scan(&placeID)
		c.NoError(err)

		query = `insert into texts(
        		place_id,
                text_place) 
				values ($1, $2)`
		_, err = c.conn.ExecContext(c.ctx, query,
			placeID,
			gofakeit.Word())
		c.NoError(err)

		query = `insert into photos(
        		place_id,
                link) 
				values ($1, $2)`
		_, err = c.conn.ExecContext(c.ctx, query,
			placeID,
			gofakeit.URL())
		c.NoError(err)

		query = `insert into videos(
        		place_id,
                link) 
				values ($1, $2)`
		_, err = c.conn.ExecContext(c.ctx, query,
			placeID,
			gofakeit.URL())
		c.NoError(err)

		query = `insert into events(
        		place_id,
                link) 
				values ($1, $2)`
		_, err = c.conn.ExecContext(c.ctx, query,
			placeID,
			gofakeit.URL())
		c.NoError(err)
	}

	var countFromDB int

	row := c.conn.QueryRowContext(c.ctx, "SELECT COUNT(*) FROM place")
	err := row.Scan(&countFromDB)
	c.NoError(err)
	if count != countFromDB {
		c.NoError(fmt.Errorf("пришло неверное количество данных таблицы place, требовалось: %v, получили: %v",
			count, countFromDB))
	}

	row = c.conn.QueryRowContext(c.ctx, "SELECT COUNT(*) FROM texts")
	err = row.Scan(&countFromDB)
	c.NoError(err)
	if count != countFromDB {
		c.NoError(fmt.Errorf("пришло неверное количество данных таблицы texts, требовалось: %v, получили: %v",
			count, countFromDB))
	}

	row = c.conn.QueryRowContext(c.ctx, "SELECT COUNT(*) FROM photos")
	err = row.Scan(&countFromDB)
	c.NoError(err)
	if count != countFromDB {
		c.NoError(fmt.Errorf("пришло неверное количество данных таблицы photos, требовалось: %v, получили: %v",
			count, countFromDB))
	}

	row = c.conn.QueryRowContext(c.ctx, "SELECT COUNT(*) FROM videos")
	err = row.Scan(&countFromDB)
	c.NoError(err)
	if count != countFromDB {
		c.NoError(fmt.Errorf("пришло неверное количество данных таблицы videos, требовалось: %v, получили: %v",
			count, countFromDB))
	}

	row = c.conn.QueryRowContext(c.ctx, "SELECT COUNT(*) FROM events")
	err = row.Scan(&countFromDB)
	c.NoError(err)
	if count != countFromDB {
		c.NoError(fmt.Errorf("пришло неверное количество данных таблицы events, требовалось: %v, получили: %v",
			count, countFromDB))
	}

	c.T().Log("большое количество данных верно сохраняются и читаются")
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
