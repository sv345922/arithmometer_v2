package giveAnswer

import (
	"context"
	"github.com/sv345922/arithmometer_v2/internal/entities"
	"github.com/sv345922/arithmometer_v2/internal/wSpace"
	"log"
	"net/http"
)

func AcceptAnswer(ctx context.Context, ws *wSpace.WorkingSpace, container entities.MessageResult) (status int, err error) {
	status = http.StatusOK
	// получаем id задачи
	id := container.Id
	// Если деление на ноль
	if container.Err == "zero division" {
		// найти корень выражения и внести ошибку деления на ноль
		root, expression, err := ws.GetRoot(id)
		if err != nil {
			log.Println(err)
			status = http.StatusInternalServerError
			return status, err
		}
		expression.SetStatus("zero division")
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
				status = http.StatusInternalServerError
			}
			expression.SetRoot(0)
			ws.Mu.Lock()
			delete(ws.AllNodes, nodeId)
			ws.Mu.Unlock()
			ws.Queue.RemoveTask(nodeId)
		}
		return status, nil
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
			status = http.StatusInternalServerError
			return status, nil
		}
		expression := ws.Expressions[node.ExpressionId]
		if expression == nil {
			log.Println("no expression in ws")
			status = http.StatusInternalServerError
			return status, nil
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
				status = http.StatusInternalServerError
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
	return status, nil
}
