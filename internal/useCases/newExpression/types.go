package newExpression

import (
	"github.com/sv345922/arithmometer_v2/internal/entities"
)

// Для получения выражения от клиента
type NewExpr struct {
	Expression string            `json:"expression"`
	Timings    *entities.Timings `json:"timings"`
}
