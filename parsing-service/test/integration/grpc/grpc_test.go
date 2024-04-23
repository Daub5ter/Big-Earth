package grpc

import (
	"context"
	"fmt"
	"net"
	"testing"
	"time"

	"parsing-service/internal/tools/grpc/parsing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// ClientSuite - структура для тестов.
type ClientSuite struct {
	suite.Suite
}

// TestGRPC запускает тесты.
func TestGRPC(t *testing.T) {
	suite.Run(t, new(ClientSuite))
}

// TestServer тестирует gRPC сервер.
func (c *ClientSuite) TestServer() {
	listenGRPC, err := net.Listen("tcp", fmt.Sprintf("%s:%s", "", "50001"))
	c.NoError(err)

	grpcSrv := grpc.NewServer()
	parsing.RegisterParsingServer(grpcSrv, parserMock{})

	go func() {
		err = grpcSrv.Serve(listenGRPC)
		c.NoError(err)
	}()

	time.Sleep(5 * time.Second)

	conn, err := grpc.Dial("localhost:50001", grpc.WithTransportCredentials(insecure.NewCredentials()))
	c.NoError(err)
	defer func() { _ = conn.Close() }()

	pc := parsing.NewParsingClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	placeInformationFromServer, err := pc.Parse(ctx, &parsing.Place{})
	c.NoError(err)

	if len(placeInformationFromServer.Events) == 0 ||
		len(placeInformationFromServer.Videos) == 0 || len(placeInformationFromServer.Photos) == 0 {

		c.NoError(fmt.Errorf("данные не полные"))
		return
	}

	if placeInformation.Text != placeInformationFromServer.Text ||
		placeInformation.Videos[0] != placeInformationFromServer.Videos[0] ||
		placeInformation.Photos[0] != placeInformationFromServer.Photos[0] ||
		placeInformation.Events[0].Link != placeInformationFromServer.Events[0].Link ||
		placeInformation.Events[0].Image != placeInformationFromServer.Events[0].Image ||
		placeInformation.Events[0].Name != placeInformationFromServer.Events[0].Name {

		c.NoError(fmt.Errorf(fmt.Sprint("данные не совпадают, ожидалось:", &placeInformation, "получено:", placeInformationFromServer)))
		return
	}

	c.T().Log("gRPC сервер работает")
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
