package calculator

import (
	"github.com/sv345922/arithmometer_v2/internal/configs"
	"github.com/sv345922/arithmometer_v2/internal/entities"
	"log"
	"time"
)

var zeroDiv = entities.ZeroDiv
var OpEr = entities.OpEr

// Выбирает оператор вычисления, производит вычисления с учетом тайминга
// и возвращает результат с ошибкой
func Do(c *entities.MessageTask) (float64, error) {
	timings := c.Timings
	Op := c.Op
	x := c.X
	y := c.Y
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
			return 0, zeroDiv
		}
		t := time.Duration(timings.Div)
		time.Sleep(configs.TConst * t)
		return x / y, nil
	default:
		return 0, OpEr
	}
}
