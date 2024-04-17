package getTask

import (
	"context"
	"encoding/json"
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

		container, httpStatus, err := PrepareTask(ctx, ws, uint64(calcId))
		if err != nil {
			log.Println(err)
			w.WriteHeader(httpStatus)
			w.Write([]byte(err.Error()))
		}
		if container == nil {
			w.WriteHeader(httpStatus)
			return
		}
		// Маршалим её
		data, _ := json.Marshal(&container) //ошибку пропускаем
		// и записываем в ответ вычислителю
		w.WriteHeader(httpStatus)
		w.Write(data)

		log.Printf("calc %d: задача %.3f%s%.3f", calcId, container.X, container.Op, container.Y)
	}
}
