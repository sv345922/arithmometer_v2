package newExpression

import (
	"context"
	"fmt"
	"github.com/sv345922/arithmometer_v2/internal/configs"
	"github.com/sv345922/arithmometer_v2/internal/entities"
	"github.com/sv345922/arithmometer_v2/internal/parser"
	"github.com/sv345922/arithmometer_v2/internal/wSpace"
	"log"
	"net/http"
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

func ProcessExpression(ctx context.Context,
	ws *wSpace.WorkingSpace,
	newClientExpression *NewExpr) (newExpression *parser.Expression, status int, err error) {
	status = http.StatusOK

	//Если тайминги не передаются, тогда они ставятся по умолчанию
	setTimingsWhileZero(newClientExpression)
	ws.Mu.Lock()
	defer ws.Mu.Unlock()
	// Обновляем полученные тайминги
	ws.Timings = newClientExpression.Timings

	newExpression = parser.NewExpression()

	tx, err := ws.DB.BeginTx(ctx, nil)
	if err != nil {
		panic(err)
	}
	// обновляем тайминги
	err = updateTimings(ctx, tx, newClientExpression.Timings)
	if err != nil {
		log.Println(err)
		status = http.StatusInternalServerError
		tx.Rollback()
		err = fmt.Errorf("timings does not accepted: %w", err)
	}
	tx.Commit()
	// Если передавалить только тайминги без выражения, то выходим
	if newClientExpression.Expression == "" {
		return newExpression, status, nil
	}

	// Парсим выражение, и проверяем его
	// Предполагается, что если парсинг с ошибкой, значит невалидное выражение
	err = newExpression.Parse(newClientExpression.Expression, *newClientExpression.Timings)
	if err != nil {
		log.Println(err)
		status = http.StatusBadRequest
		err = fmt.Errorf("invalid expression: %w", err)
	}
	// Добавляем новое выражение в WorkingSpace
	// создаем транзакцию
	tx, err = ws.DB.BeginTx(ctx, nil)
	if err != nil {
		panic(err)
	}
	// Добавляем в AllNodes
	nodes, err, tx := ws.InsertToAllNodes(ctx, tx, newExpression.Nodes)
	if err != nil {
		log.Printf("add to all nodes failed: %v", err)
		status = http.StatusInternalServerError
		tx.Rollback()
		return nil, status, fmt.Errorf("expression does not accepted")
	}
	// Добавляем в Expressions
	newExpression.Id, err = ws.AddToExpressions(ctx, tx, TransformParseExpression(newExpression))
	if err != nil {
		log.Printf("add to expressions failed: %v", err)
		status = http.StatusInternalServerError
		tx.Rollback()
		return nil, status, fmt.Errorf("expression does not accepted")
	}
	// Установить id выражения в узлы
	for _, node := range nodes {
		node.ExpressionId = newExpression.Id
	}

	// Добавляем в таблицу queue БД
	err = ws.Queue.AddExpressionNodes(ctx, tx, nodes)
	if err != nil {
		log.Println(err)
		status = http.StatusInternalServerError
		tx.Rollback()
		return nil, status, fmt.Errorf("expression does not accepted")
	}

	// Обновляем поле expressionID в node в таблице AllNodes
	for _, node := range nodes {
		_, err = tx.ExecContext(ctx,
			`UPDATE allNodes SET expressionId = $1 WHERE id = $2;`,
			node.ExpressionId,
			node.Id,
		)
		if err != nil {
			log.Println(err)
			status = http.StatusInternalServerError
			tx.Rollback()
			return nil, status, fmt.Errorf("expression does not accepted")
		}
	}

	// Завершаем транзакцию
	tx.Commit()
	return newExpression, status, nil
}
