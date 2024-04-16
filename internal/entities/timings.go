package entities

import (
	"fmt"
	"github.com/sv345922/arithmometer_v2/internal/configs"
)

var tConst = configs.TConst

// Timings Тайминги для операторов
type Timings struct {
	Plus  int `json:"plus"`
	Minus int `json:"minus"`
	Mult  int `json:"mult"`
	Div   int `json:"div"`
}

// Стрингер
func (t *Timings) String() string {
	return fmt.Sprintf("+: %ds, -: %ds, *: %ds, /: %ds", t.Plus, t.Minus, t.Mult, t.Div)
}

// GetDuration Возвращает время выполнения конкретной операции
// Если оператор неизвестен, возвращает нулевую длительность
func (t *Timings) GetDuration(op string) int {
	switch op {
	case "+":
		return t.Plus
	case "-":
		return t.Minus
	case "*":
		return t.Mult
	case "/":
		return t.Div
	default:
		return 0
	}
}
