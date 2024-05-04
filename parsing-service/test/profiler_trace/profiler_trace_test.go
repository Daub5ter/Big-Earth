package profiler_trace

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/brianvoe/gofakeit/v7"
	"net"
	"os"
	"parsing-service/internal/handlers"
	"parsing-service/internal/tools/config"
	parsinggrpc "parsing-service/internal/tools/grpc/parsing"
	"parsing-service/internal/tools/parsing"
	"parsing-service/internal/tools/postgres"
	"runtime/pprof"
	"runtime/trace"
	"testing"

	tc "github.com/testcontainers/testcontainers-go/modules/compose"
	"github.com/testcontainers/testcontainers-go/wait"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"parsing-service/pkg/dbconn/pgsql"
)

// BenchmarkPprofTrace - основная функция, записывающая данные
// в файлы профайлера и трэйса.
func BenchmarkPprofTrace(b *testing.B) {
	ctx := context.Background()

	b.Log("запуск контейнера...")
	cs := containers{
		ctx: ctx,
	}
	compose, err := cs.runCompose()
	if err != nil {
		b.Errorf("ошибка запуска контейнера: %v", err)
		return
	}

	defer func() {
		b.Log("завершение работы контейеров...")
		err = cs.downCompose(compose)
		if err != nil {
			b.Errorf("ошибка остовновки контейнера: %v", err)
			return
		}
	}()

	b.Log("настройка конфигурации...")
	b.Setenv("DSN", "host=localhost port=5432 user=some_user password=some_password dbname=some_db sslmode=disable timezone=UTC connect_timeout=5")

	cfg, err := config.NewConfig("profiler_trace_test_config.yaml")
	if err != nil {
		b.Errorf("ошибка прочтения файла конфигруаций: %v", err)
		return
	}

	b.Log("запуск основной функции main...")
	conn, err := main(cfg)
	if err != nil {
		b.Errorf("ошибка запуска main функции: %v", err)
		return
	}

	db := database{
		ctx:  ctx,
		conn: conn,
	}

	err = db.saveData()
	if err != nil {
		b.Errorf("ошибка сохранения корректных данных: %v", err)
		return
	}

	b.Log("создание и запуск cpu профайлера...")
	cpuProf, err := os.Create("cpu.prof")
	if err != nil {
		b.Errorf("ошибка создания cpu профайлера: %v", err)
		return
	}
	defer func() { _ = cpuProf.Close() }()

	err = pprof.StartCPUProfile(cpuProf)
	if err != nil {
		b.Errorf("ошибка запуска cpu профайлера: %v", err)
		return
	}
	defer pprof.StopCPUProfile()

	b.Log("создание и запусе heap профайлера...")
	heapProf, err := os.Create("heap.prof")
	if err != nil {
		b.Errorf("ошибка создания heap профайлера: %v", err)
		return
	}
	defer func() { _ = heapProf.Close() }()

	err = pprof.WriteHeapProfile(heapProf)
	if err != nil {
		b.Errorf("ошибка запуска heap профайлера: %v", err)
		return
	}

	b.Log("создание и запусе trace...")
	traceFile, err := os.Create("trace.out")
	if err != nil {
		b.Errorf("ошибка создания trace: %v", err)
		return
	}
	defer func() { _ = traceFile.Close() }()

	err = trace.Start(traceFile)
	if err != nil {
		b.Errorf("ошибка запуска trace: %v", err)
		return
	}
	defer trace.Stop()

	b.Log("отправка запросов...")
	err = sendRequest(ctx)
	if err != nil {
		b.Errorf("ошибка отправки запросов: %v", err)
		return
	}

	pprof.StopCPUProfile()
	_ = heapProf.Close()
	trace.Stop()

	err = conn.Close()
	if err != nil {
		b.Errorf("ошибка при завершении работы базы данных %v", err)
		return
	}
}

// main - функция main из корня cmd, переписанная для бенчмарка
// (отсутсвует graceful shutdown), контекст передается в main, как и конфиг,
// логер не настраивается, функция возвращает соединение с бд, ошибку).
func main(cfg config.Config) (*sql.DB, error) {
	// Создание парсера.
	parse := parsing.NewParsing(cfg.(config.ServerConfig).Timeout())

	// Соединение с БД.
	conn := pgsql.ConnectToDB(cfg.(config.DatabaseConfig).DSN())
	if conn == nil {
		return nil, fmt.Errorf("ошибка подключения к Postgres")
	}

	db, err := postgres.NewDB(conn, cfg)
	if err != nil {
		return nil, fmt.Errorf("ошибка подключения к базе данных: %v", err)
	}

	// Настройка конфигурации сервера.
	listenGRPC, err := net.Listen("tcp",
		fmt.Sprintf("%s:%s", cfg.(config.ServerConfig).Domain(), cfg.(config.ServerConfig).Port()))
	if err != nil {
		return nil, fmt.Errorf("ошибка прослушивания порта gRPC: %v", err)
	}

	grpcSrv := grpc.NewServer()
	parsinggrpc.RegisterParsingServer(grpcSrv, handlers.NewParser(db, parse))

	go func() {
		err = grpcSrv.Serve(listenGRPC)
		if err != nil {
			return
		}
	}()

	return conn, nil
}

// containers - структура для работы с запуском и завершенем контейнера.
type containers struct {
	ctx context.Context
}

// runCompose запускает docker compose kafka + zookeeper
func (cs containers) runCompose() (tc.ComposeStack, error) {
	compose, err := tc.NewDockerComposeWith(tc.WithStackFiles("profiler_trace_test.yml"),
		tc.StackIdentifier("profiler_trace_test_identifier"))
	if err != nil {
		return nil, err
	}

	err = compose.Up(cs.ctx, tc.Wait(true))
	if err != nil {
		return nil, err
	}

	compose.WaitForService("postgres_test", wait.ForLog("Awaiting socket connections on 0.0.0.0:5432"))

	return compose, nil
}

// downCompose выключает docker compose и удаляет образы.
func (cs containers) downCompose(compose tc.ComposeStack) error {
	return compose.Down(cs.ctx, tc.RemoveImagesAll, tc.RemoveVolumes(true))
}

// database - структура для работы с базой данных.
type database struct {
	ctx  context.Context
	conn *sql.DB
}

// saveData сохраняет данные в бд.
func (db database) saveData() error {
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

// count - число отправляемых запросов.
var count = 50

// sendRequest отправляет запросы на сервер и получает данные.
func sendRequest(ctx context.Context) error {
	for i := 1; i < count; i++ {
		conn, err := grpc.Dial("localhost:50001", grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			return fmt.Errorf("не создается клиент grpc: %v", err)
		}

		c := parsinggrpc.NewParsingClient(conn)

		placeInformationFromServer, err := c.Parse(ctx,
			&parsinggrpc.Place{
				Country: country,
				City:    city,
			})
		if err != nil {
			_ = conn.Close()
			return fmt.Errorf("не отправляется запрос: %v", err)
		}

		if len(placeInformationFromServer.Events) == 0 || placeInformationFromServer.Videos[0] != videoLink ||
			placeInformationFromServer.Photos[0] != photoLink || placeInformationFromServer.Text != text {

			_ = conn.Close()
			return fmt.Errorf("ошибка, не полныее данные: %s", fmt.Sprint(placeInformationFromServer))
		}

		_ = conn.Close()
	}

	return nil
}

// country, city, text, photoLink, videoLink, eventsLink
// - данные для тестирования корректности базы данных.
var (
	country = "Russia"
	city    = "Moscow"

	text = "Moscow is a capital of Russia"

	photoLink = gofakeit.URL()

	videoLink = gofakeit.URL()

	eventsLink = "https://msk.kassir.ru/"
)
