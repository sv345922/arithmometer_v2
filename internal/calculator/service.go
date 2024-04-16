package calculator

import (
	"fmt"
	"github.com/sv345922/arithmometer_v2/internal/configs"
	"github.com/sv345922/arithmometer_v2/internal/entities"
	"log"
	"time"
)

func NewCalculator() *Calculator {
	return &Calculator{
		Id:   entities.GetDelta(5),
		Task: new(entities.MessageTask),
		Ch:   make(chan entities.MessageResult),
	}
}
func (c *Calculator) Do() error {
	log.Printf("calculator %d is runing", c.Id)
	for {
		ok, err := c.GetTask()
		if err != nil {
			log.Println("ошибка получения задачи", err)
		}
		// Окрестратор не дал задание
		if !ok {
			time.Sleep(5 * time.Second)
			log.Println("tik ")
			time.Sleep(5 * time.Second)
			continue
		}
		// запускаем задачу в горутине
		go func() {
			res, err := c.Calculate()
			fmt.Printf("res: %f, time: %v\n", res, c.Task.Timings)
			errString := "nil"
			if err != nil {
				errString = err.Error()
			}
			c.Ch <- entities.MessageResult{
				Id:     c.Task.Id,
				Result: res,
				Err:    errString,
			}
		}()
		answer := <-c.Ch
		log.Printf("%.3f%s%.3f=%f, ошибка %s\n",
			c.Task.X,
			c.Task.Op,
			c.Task.Y,
			answer.Result,
			answer.Err,
		)
		// отправляем ответ, до тех пор пока он не будет принят
		for {
			err = c.SendAnswer(answer)
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
	log.Println("задача получена", x, Op, y)
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
