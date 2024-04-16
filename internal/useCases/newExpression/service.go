package newExpression

import (
	"github.com/sv345922/arithmometer_v2/internal/configs"
	"github.com/sv345922/arithmometer_v2/internal/entities"
	"github.com/sv345922/arithmometer_v2/internal/parser"
)

var t = configs.DefaultTimings

// Преобразует тип *parser.Expression в *entities.Expression
func TransformParseExpression(expression *parser.Expression) *entities.Expression {
	return &entities.Expression{
		Id:         expression.Id,
		UserTask:   expression.UserTask,
		ResultExpr: expression.ResultExpr,
		Status:     expression.Status,
		RootId:     expression.Root.NodeId,
	}
}

func setTimingsWhileZero(n *NewExpr) {
	if n.Timings == nil {
		n.Timings = &entities.Timings{
			Plus:  t,
			Minus: t,
			Mult:  t,
			Div:   t,
		}
	}
	times := n.Timings
	if times.Plus == 0 {
		times.Plus = t
	}
	if times.Minus == 0 {
		times.Minus = t
	}
	if times.Mult == 0 {
		times.Mult = t
	}
	if times.Div == 0 {
		times.Div = t
	}
}
