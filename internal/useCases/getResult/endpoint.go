package getResult

import (
	"fmt"
	"github.com/sv345922/arithmometer_v2/internal/wSpace"
	"log"
	"net/http"
	"strconv"
)

// Обрабатывает запросы клиента о проверке результата вычислений
func GetResult(ws *wSpace.WorkingSpace) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		// Проверить метод
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("требуется метод Get"))
			return
		}
		// Читаем id из параметров запроса
		id := r.URL.Query().Get("id")

		// при пустом ID возвращать все выражения (пользователя)
		if id == "" {
			// TODO пока нет USERS используется userID=0
			ws.Mu.Lock()
			result, err := GetExpressions(0, ws.Expressions)
			ws.Mu.Unlock()
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(result))
			log.Printf("отправлены все выражения пользователя %d", 0)
			return
		}

		// TODO реализовать проверку пользователя и его выражений

		// преобразуем id в число
		idInt, _ := strconv.ParseUint(id, 10, 64)
		// Поиск выражения в списке выражений
		log.Printf("Выражени %d запрошено клиентом", idInt)
		ws.Mu.RLock()
		expression, ok := ws.Expressions[idInt]
		ws.Mu.RUnlock()
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("id not found"))
			return
		}

		w.WriteHeader(http.StatusOK)
		switch expression.Status {
		case "done":
			w.Write([]byte(fmt.Sprintf("%s = %f", expression.UserTask, expression.ResultExpr)))
			return
		case "zero division":
			w.Write([]byte(fmt.Sprintf("%s = zero division error", expression.UserTask)))
			return
		default:
			w.Write([]byte(fmt.Sprintf("%s ... calculating", expression.UserTask)))
		}
	}
}
