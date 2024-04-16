package getTask

import (
	"context"
	"encoding/json"
	"github.com/sv345922/arithmometer_v2/internal/entities"
	"github.com/sv345922/arithmometer_v2/internal/wSpace"
	"log"
	"net/http"
	"strconv"
)

// Даёт задачу калькулятору
func GetTask(ctx context.Context, ws *wSpace.WorkingSpace) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// Проверить метод
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("требуется метод Get"))
			return
		}
		// Читаем id вычислителя из параметров запроса
		id := r.URL.Query().Get("id")
		if id == "" {
			log.Println("не найден id в запросе вычислителя")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		log.Printf("calc %s tik", id)
		calcId, err := strconv.Atoi(id)
		if err != nil {
			log.Println(id, "id вычислителя не число", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		// TODO (элементы очереди) Обновляем очередь задач, чтобы обработать просроченные
		ws.Queue.CheckDeadlines()
		// TODO (элементы очереди) Получаем задачу из очереди
		task := ws.Queue.GetTask()

		if task == nil {
			// Если активных задач нет
			log.Println("tok")
			w.WriteHeader(http.StatusNoContent)
			return
		}
		// Устанавливаем id калькулятора
		task.CalcId = uint64(calcId)
		// Устанавливаем длительность операции
		task.Duration = ws.Timings.GetDuration(task.Node.Op)
		// Устанавливаем дедлайн для задачи
		timeout := task.Duration * 15 / 10
		task.SetDeadline(timeout)
		// обновляем поля в БД
		err = UpdateInGetTask(ctx, ws.DB, task)
		if err != nil {
			log.Printf("%v", err)
		}

		// Создаем структуру для передачи вычислителю
		container := entities.MessageTask{
			Id:      task.Node.Id,
			X:       task.X,
			Y:       task.Y,
			Op:      task.Node.Op,
			Timings: ws.Timings,
		}
		// Маршалим её
		data, _ := json.Marshal(&container) //ошибку пропускаем
		// и записываем в ответ вычислителю
		w.Write(data)
		log.Printf("calc %d, задача %.3f%s%.3f", task.CalcId, task.X, task.Node.Op, task.Y)
	}
}
