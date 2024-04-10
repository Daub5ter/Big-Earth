package grpc

import (
	"context"
	"fmt"
	"net"
	"testing"
	"time"

	"api-gateway/internal/tools/grpc/parsing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestGRPC(t *testing.T) {
	listenGRPC, err := net.Listen("tcp",
		fmt.Sprintf("%s:%s", "", "50001"))
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

	conn, err := grpc.Dial("localhost:50001", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Errorf("не создается клиент grpc: %v", err)
		return
	}
	defer func() { _ = conn.Close() }()

	c := parsing.NewParsingClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

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

		t.Errorf("данные не совпадают")
		return
	}
}

type parserMock struct {
	parsing.UnimplementedParsingServer
}

var placeInformation = parsing.PlaceInformation{
	Text:   "text",
	Photos: []string{"photo"},
	Videos: []string{"video"},
	Events: []*parsing.Event{
		{
			Name:  "name",
			Image: "image",
			Link:  "link",
		}},
}

func (p parserMock) Parse(ctx context.Context, place *parsing.Place) (*parsing.PlaceInformation, error) {
	return &placeInformation, nil
}