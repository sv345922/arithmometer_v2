package giveAnswer

import (
	"context"
	"encoding/json"
	"github.com/sv345922/arithmometer_v2/internal/entities"
	"github.com/sv345922/arithmometer_v2/internal/wSpace"
	"log"
	"net/http"
)

// var ZeroDiv = errors.New("zero division")

// Обработчик, принимает от вычислителя ответ
func GiveAnswer(ctx context.Context, ws *wSpace.WorkingSpace) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// Проверить что это метод POST
		if r.Method != http.MethodPost {
			log.Println("метод не POST")
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte("требуется метод POST"))
			return
		}

		// Читаем тело запроса, в котором записан ответ
		defer r.Body.Close()
		var container entities.MessageResult
		err := json.NewDecoder(r.Body).Decode(&container)
		if err != nil {
			log.Println("ошибка json при обработке ответа вычислителя")
			return
		}
		log.Printf("Получен ответ от вычислителя %f, ошибка %s\n",
			container.Result,
			container.Err,
		)
		// получаем id задачи
		id := container.Id
		// Если деление на ноль
		if container.Err == "zero division" {
			// найти корень выражения и внести ошибку деления на ноль
			root, expression, err := ws.GetRoot(id)
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
			}
			expression.Status = "zero division"
			// обновить в БД выражение
			err = UpdateExpressionResult(ctx, ws.DB, expression)
			if err != nil {
				log.Println(err)
			}

			// получить все узлы выражения
			nodesId := make([]uint64, 0)
			ws.GetExpressionNodesID(root.Id, &nodesId)

			// удалить все узлы дерева выражения из allNodes
			// и из queue
			for _, nodeId := range nodesId {
				// Удаляем из БД записи узлов
				err = DeleteNodeTask(ctx, ws.DB, nodeId)
				if err != nil {
					log.Println(err)
					w.WriteHeader(http.StatusInternalServerError)
				}
				expression.RootId = 0
				ws.Mu.Lock()
				delete(ws.AllNodes, nodeId)
				ws.Mu.Unlock()
				ws.Queue.RemoveTask(nodeId)
			}
			w.WriteHeader(http.StatusOK)
			return
		}

		// Обновляем очередь задач с учетом выполненной задачи и заносим результат вычисления
		rootFlag, err := ws.Queue.AddAnswer(id, container.Result)
		if err != nil {
			log.Println("ошибка обновления задач:", err)
		}
		// Если вычислен корень выражения
		// обновим выражение и запишем результат
		if rootFlag {
			//
			// удаляем все узлы выражения из очереди и из allNodes
			// обновляем статус выражения
			node := ws.AllNodes[id]
			if node == nil {
				log.Println("no node in ws")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			expression := ws.Expressions[node.ExpressionId]
			if expression == nil {
				log.Println("no expression in ws")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			expression.Status = "done"
			expression.ResultExpr = container.Result
			// Обновить expression в БД
			err = UpdateExpressionResult(ctx, ws.DB, expression)
			if err != nil {
				log.Println(err)
			}

			// получить все узлы выражения
			nodesId := make([]uint64, 0)
			ws.GetExpressionNodesID(id, &nodesId)

			// удалить все узлы дерева выражения
			// и из очереди задач
			for _, nodeId := range nodesId {
				// Удаляем из БД записи узлов
				err = DeleteNodeTask(ctx, ws.DB, nodeId)
				if err != nil {
					log.Println(err)
					w.WriteHeader(http.StatusInternalServerError)
				}
				expression.RootId = 0

				ws.Mu.Lock()
				delete(ws.AllNodes, nodeId)
				ws.Mu.Unlock()
				ws.Queue.RemoveTask(nodeId)
			}

		} else {
			// Удаляем решенную задачу из БД
			err = DeleteNodeTask(ctx, ws.DB, id)
			if err != nil {
				log.Println(err)
			}
			// и из queue
			ws.Queue.RemoveTask(id)
		}
		w.WriteHeader(http.StatusOK)
	}
}
