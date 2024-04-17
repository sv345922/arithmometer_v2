package server

import (
	"context"
	"fmt"
	"github.com/sv345922/arithmometer_v2/internal/configs"
	"github.com/sv345922/arithmometer_v2/internal/entities"
	"github.com/sv345922/arithmometer_v2/internal/useCases/getTask"
	"github.com/sv345922/arithmometer_v2/internal/useCases/giveAnswer"
	"github.com/sv345922/arithmometer_v2/internal/wSpace"
	"google.golang.org/grpc"
	"log"
	"net"

	pb "github.com/sv345922/arithmometer_v2/internal/proto"
)

// pb "github.com/sv345922/arithmometer_v2/internal/proto"

type Server struct {
	ctx context.Context
	ws  *wSpace.WorkingSpace
	pb.CalcServiceServer
}

func NewServer(ctx context.Context, ws *wSpace.WorkingSpace) *Server {
	return &Server{ctx: ctx, ws: ws}
}

func StartGRPCServer(ctx context.Context, ws *wSpace.WorkingSpace) error {
	host := configs.Host

	addr := fmt.Sprintf("%s:%s", host, configs.GRPCPort)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Printf("failed to listen: %v", err)
		return err
	}
	log.Printf("gRPC listening on %s", addr)
	grpcServer := grpc.NewServer()
	pb.RegisterCalcServiceServer(grpcServer, NewServer(ctx, ws))
	if err := grpcServer.Serve(lis); err != nil {
		log.Printf("error serving grpc: %v", err)
		return fmt.Errorf("error serving grpc: %v", err)
	}
	return nil
}

func (s *Server) GetTask(ctx context.Context, in *pb.CalculatorID) (*pb.MessageTask, error) {
	container, _, err := getTask.PrepareTask(ctx, s.ws, in.Id)
	if err != nil {
		log.Printf("error preparing task: %v", err)
		return nil, err
	}
	if container == nil {
		return nil, nil
	}
	result := &pb.MessageTask{
		Id:    container.Id,
		X:     container.X,
		Y:     container.Y,
		Op:    container.Op,
		Plus:  int64(container.Timings.Plus),
		Minus: int64(container.Timings.Minus),
		Mult:  int64(container.Timings.Mult),
		Div:   int64(container.Timings.Div),
	}
	return result, nil
}
func (s *Server) SendAnswer(ctx context.Context, in *pb.MessageResult) (*pb.Receipt, error) {
	container := entities.MessageResult{
		Id:     uint64(in.Id),
		Result: in.Result,
		Err:    in.Err,
	}
	result := &pb.Receipt{Ok: false}
	status, err := giveAnswer.AcceptAnswer(ctx, s.ws, container)
	if status != 200 || err != nil {
		return result, err
	}
	result.Ok = true
	return result, nil
}
