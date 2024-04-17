package client

import (
	"fmt"
	"github.com/sv345922/arithmometer_v2/internal/configs"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"

	pb "github.com/sv345922/arithmometer_v2/internal/proto"
)

// pb "github.com/sv345922/arithmometer_v2/internal/proto"
// "google.golang.org/grpc"
//	"google.golang.org/grpc/credentials/insecure"

// Запускает клиент для gRPC сервера и возвращает:
// клиента - для использования методов
// соединение - чтобы закрыть после вызова метода
// ошибку
func StartGrpcClient() (pb.CalcServiceClient, *grpc.ClientConn, error) {
	addr := fmt.Sprintf("%s:%s", configs.Host, configs.GRPCPort)

	// установим соединение
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Println("could not connect to grpc server: ", err)
		return nil, nil, err
	}
	// закроем соединение, когда выйдем из функции
	client := pb.NewCalcServiceClient(conn)
	return client, conn, nil
}
