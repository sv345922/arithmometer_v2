package calculator

import (
	"context"
	"github.com/sv345922/arithmometer_v2/internal/configs"
	"github.com/sv345922/arithmometer_v2/internal/entities"
	"log"
	"time"
)

func NewCalculator() *Calculator {
	return &Calculator{
		Id:   uint64(entities.GetDelta(5)),
		Task: new(entities.MessageTask),
	}
}
func (c *Calculator) Do(ctx context.Context) error {
	log.Printf("calculator %d is runing", c.Id)
	for {
		ok, err := c.GetTaskGrpc(ctx)
		if err != nil {
			log.Println("ошибка получения задачи", err)
		}
		// Окрестратор не дал задание
		if !ok {
			//log.Println("tik ")
			time.Sleep(5 * time.Second)
			continue
		}
		res, err := c.Calculate()
		// fmt.Printf("res: %f, time: %v\n", res, c.Task.Timings) // todo delete
		errString := "nil"
		if err != nil {
			errString = err.Error()
		}

		answer := &entities.MessageResult{
			Id:     c.Task.Id,
			Result: res,
			Err:    errString,
		}
		log.Printf("%.3f%s%.3f=%f, ошибка %s\n",
			c.Task.X,
			c.Task.Op,
			c.Task.Y,
			answer.Result,
			answer.Err,
		)
		// отправляем ответ, до тех пор пока он не будет принят
		for {
			err = c.SendAnswerGrpc(ctx, answer)
			if err == nil {
				break
			}
		}
	}
}
func (c *Calculator) Calculate() (float64, error) {
	task := c.Task
	timings := task.Timings

	Op := task.Op
	x := task.X
	y := task.Y
	log.Printf("задача %f%s%f вычисляется, id: %d", x, Op, y, c.Task.Id)
	switch Op {
	case "+":
		t := time.Duration(timings.Plus)
		time.Sleep(configs.TConst * t)
		return x + y, nil
	case "-":
		t := time.Duration(timings.Minus)
		time.Sleep(configs.TConst * t)
		return x - y, nil
	case "*":
		t := time.Duration(timings.Mult)
		time.Sleep(configs.TConst * t)
		return x * y, nil
	case "/":
		if y == 0 {
			return 0, entities.ZeroDiv
		}
		t := time.Duration(timings.Div)
		time.Sleep(configs.TConst * t)
		return x / y, nil
	default:
		return 0, entities.OpEr
	}
}
