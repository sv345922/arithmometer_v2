package getTask

import (
	"context"
	"fmt"
	"github.com/sv345922/arithmometer_v2/internal/entities"
	"github.com/sv345922/arithmometer_v2/internal/wSpace"
	"log"
	"net/http"
)

func PrepareTask(ctx context.Context, ws *wSpace.WorkingSpace, calcId uint64) (*entities.MessageTask, int, error) {
	ws.Mu.Lock()
	defer ws.Mu.Unlock()
	// TODO (элементы очереди) Обновляем очередь задач, чтобы обработать просроченные
	// ws.Queue.CheckDeadlines() повторяется в GetTask
	// TODO (элементы очереди) Получаем задачу из очереди
	task := ws.Queue.GetTask()

	if task == nil {
		// Если активных задач нет
		return nil, http.StatusNoContent, nil
	}
	// Устанавливаем id калькулятора
	task.SetCalc(calcId)
	// Устанавливаем длительность операции
	task.SetDuration(ws.Timings.GetDuration(task.Node.Op))
	// Устанавливаем дедлайн для задачи
	timeout := task.Duration * 15 / 10
	task.SetDeadline(timeout)
	// обновляем поля в БД
	err := UpdateInGetTask(ctx, ws.DB, task)
	if err != nil {
		log.Printf("%v", err)
		return nil, http.StatusInternalServerError, fmt.Errorf("UpdateInGetTask: %v", err)
	}
	container := entities.MessageTask{
		Id:      task.Node.Id,
		X:       task.X,
		Y:       task.Y,
		Op:      task.Node.Op,
		Timings: ws.Timings,
	}

	return &container, http.StatusOK, nil
}
