package grpc

import (
	"context"
	"fmt"
	"net"
	"testing"
	"time"

	"parsing-service/internal/tools/grpc/parsing"

	"github.com/brianvoe/gofakeit/v7"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// TestGRPC - основная функция теста gRPC сервера.
func TestGRPC(t *testing.T) {
	t.Log("запуск gRPC сервера...")
	listenGRPC, err := net.Listen("tcp", fmt.Sprintf("%s:%s", "", "50001"))
	if err != nil {
		t.Errorf("ошибка прослушивания порта gRPC: %v", err)
		return
	}

	grpcSrv := grpc.NewServer()
	parsing.RegisterParsingServer(grpcSrv, parserMock{})

	go func() {
		err = grpcSrv.Serve(listenGRPC)
		if err != nil {
			t.Errorf("ошибка запуска сервера gRPC: %v", err)
			return
		}
	}()

	time.Sleep(5 * time.Second)

	t.Log("подключение к серверу...")
	conn, err := grpc.Dial("localhost:50001", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Errorf("не создается клиент grpc: %v", err)
		return
	}
	defer func() { _ = conn.Close() }()

	c := parsing.NewParsingClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	t.Log("отправка и проверка запросов...")
	placeInformationFromServer, err := c.Parse(ctx, &parsing.Place{})
	if err != nil {
		t.Errorf("не отправляется запрос: %v", err)
		return
	}

	if len(placeInformationFromServer.Events) == 0 ||
		len(placeInformationFromServer.Videos) == 0 || len(placeInformationFromServer.Photos) == 0 {

		t.Errorf("данные не полные")
		return
	}

	if placeInformation.Text != placeInformationFromServer.Text ||
		placeInformation.Videos[0] != placeInformationFromServer.Videos[0] ||
		placeInformation.Photos[0] != placeInformationFromServer.Photos[0] ||
		placeInformation.Events[0].Link != placeInformationFromServer.Events[0].Link ||
		placeInformation.Events[0].Image != placeInformationFromServer.Events[0].Image ||
		placeInformation.Events[0].Name != placeInformationFromServer.Events[0].Name {

		t.Error("данные не совпадают, ожидалось:", &placeInformation, "получено:", placeInformationFromServer)
		return
	}
}

// parserMock - мок парсера.
type parserMock struct {
	parsing.UnimplementedParsingServer
}

// placeInformation - модель информации парсинга.
var placeInformation = parsing.PlaceInformation{
	Text:   gofakeit.Word(),
	Photos: []string{gofakeit.URL()},
	Videos: []string{gofakeit.URL()},
	Events: []*parsing.Event{
		{
			Name:  gofakeit.Name(),
			Image: gofakeit.URL(),
			Link:  gofakeit.URL(),
		}},
}

// Parse - мок функции парсинга данных.
func (p parserMock) Parse(ctx context.Context, place *parsing.Place) (*parsing.PlaceInformation, error) {
	return &placeInformation, nil
}
