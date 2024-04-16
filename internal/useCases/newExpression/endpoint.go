package newExpression

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/sv345922/arithmometer_v2/internal/parser"
	"github.com/sv345922/arithmometer_v2/internal/wSpace"
	"log"
	"net/http"
)

// Длительность по умолчанию

// Обработчик создания нового выражения
func NewExpression(ctx context.Context, ws *wSpace.WorkingSpace) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// Проверить что это запрос POST
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("требуется метод POST"))
			return
		}
		// Читаем тело запроса, в котором записано выражение и тайминги операций
		var newClientExpression NewExpr
		err := json.NewDecoder(r.Body).Decode(&newClientExpression)
		defer r.Body.Close()
		if err != nil {
			log.Println("ошибка POST запроса")
			return
		}
		//Если тайминги не передаются, тогда они ставятся по умолчанию
		setTimingsWhileZero(&newClientExpression)

		// Обновляем полученные тайминги
		ws.Timings = newClientExpression.Timings
		// Парсим выражение, и проверяем его
		// Предполагается, что если парсинг с ошибкой, значит невалидное выражение
		newExpression := parser.NewExpression()

		err = newExpression.Parse(newClientExpression.Expression, *newClientExpression.Timings)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("invalid expression"))
			return
		}
		// Добавляем новое выражение в WorkingSpace
		// создаем транзакцию
		tx, err := ws.DB.BeginTx(ctx, nil)
		if err != nil {
			panic(err)
		}
		// Добавляем в AllNodes
		nodes, err, tx := ws.InsertToAllNodes(ctx, tx, newExpression.Nodes)
		if err != nil {
			log.Printf("add to all nodes failed: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("expression does not accepted"))
			_ = tx.Rollback()
			return
		}
		// Добавляем в Expressions
		newExpression.Id, err = ws.AddToExpressions(ctx, tx, TransformParseExpression(newExpression))
		if err != nil {
			log.Printf("add to expressions failed: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("expression does not accepted"))
			tx.Rollback()
			return
		}
		// Установить id выражения в узлы
		for _, node := range nodes {
			node.ExpressionId = newExpression.Id
		}

		//// Получаем список преобразованных узлов
		//nodes := make([]*entities.Node, 0, len(newExpression.Nodes))
		//for _, node := range newExpression.Nodes {
		//	nodes = append(nodes, parser.TransformNode(node))
		//}

		// Добавляем в таблицу queue БД
		err = ws.Queue.AddExpressionNodes(ctx, tx, nodes)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("expression does not accepted"))
			tx.Rollback()
			return
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
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("expression does not accepted"))
				tx.Rollback()
				return
			}
		}

		// обновляем тайминги
		err = updateTimings(ctx, tx, newClientExpression.Timings)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("expression does not accepted"))
			tx.Rollback()
			return
		}

		// Завершаем транзакцию
		tx.Commit()
		// Записываем в тело ответа id выражения
		w.Write([]byte(fmt.Sprintf("%d", newExpression.Id)))

		log.Printf("Method: %s, Expression: %s, Timings: %s, id: %d",
			r.Method,
			newExpression.UserTask,
			newExpression.Times.String(),
			newExpression.Id,
		)
	}
}
