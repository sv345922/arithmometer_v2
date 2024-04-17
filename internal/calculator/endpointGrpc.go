package calculator

import (
	"context"
	"errors"
	"github.com/sv345922/arithmometer_v2/internal/entities"
	"github.com/sv345922/arithmometer_v2/internal/grps/client"
	"github.com/sv345922/arithmometer_v2/internal/proto"
	"log"
)

func (c *Calculator) GetTaskGrpc(ctx context.Context) (bool, error) {
	client, conn, err := client.StartGrpcClient()
	defer conn.Close()
	if err != nil {
		log.Println("could not connect to grpc server: ", err)
		return false, err
	}
	id := &proto.CalculatorID{Id: c.Id}
	container, err := client.GetTask(ctx, id)
	//if err != nil {
	//	log.Println("could not get task: ", err)
	//	return false, err
	//}
	if container == nil || container.Id == 0 {
		return false, nil
	}
	task := &entities.MessageTask{
		Id: container.Id,
		X:  container.X,
		Y:  container.Y,
		Op: container.Op,
		Timings: &entities.Timings{
			Plus:  int(container.Plus),
			Minus: int(container.Minus),
			Mult:  int(container.Mult),
			Div:   int(container.Div),
		},
	}
	c.Task = task
	return true, nil
}

func (c *Calculator) SendAnswerGrpc(ctx context.Context, container *entities.MessageResult) error {
	cl, conn, err := client.StartGrpcClient()
	defer conn.Close()
	if err != nil {
		log.Println("could not connect to grpc server: ", err)
		return err
	}
	answer := &proto.MessageResult{
		Id:     int64(container.Id),
		Result: container.Result,
		Err:    container.Err,
	}
	reciept, err := cl.SendAnswer(ctx, answer)
	//if err != nil {
	//	log.Println("could not send answer: ", err)
	//	return err
	//}
	if reciept == nil {
		return errors.New("could not send answer")
	}
	ok := reciept.Ok
	if !ok {
		return errors.New("could not send answer")
	}
	return nil
}
