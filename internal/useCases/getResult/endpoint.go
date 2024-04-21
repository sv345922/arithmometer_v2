package getResult

import (
	"context"
	"fmt"
	"github.com/sv345922/arithmometer_v2/internal/wSpace"
	"log"
	"net/http"
	"strconv"
)

// Обрабатывает запросы клиента о проверке результата вычислений
func GetResult(ctx context.Context, ws *wSpace.WorkingSpace) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		// Проверить метод
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("требуется метод Get"))
			return
		}
		// Читаем id из параметров запроса
		id := r.URL.Query().Get("id")
		user, err := strconv.ParseUint(r.Header["X-Username"][0], 10, 64)
		if err != nil {
			log.Printf("Error parsing user id: %v", err)
		}
		// при пустом ID возвращать все выражения (пользователя)
		if id == "" {
			ws.Mu.Lock()
			result, err := GetExpressions(user, ws.Expressions)
			ws.Mu.Unlock()
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(result))
			log.Printf("отправлены все выражения пользователя %d", 0)
			return
		}

		// преобразуем id в число
		idInt, _ := strconv.ParseUint(id, 10, 64)
		// Поиск выражения в списке выражений
		log.Printf("Выражение %d запрошено клиентом", idInt)
		ws.Mu.RLock()
		expression, ok := ws.Expressions[idInt]

		// Проверка на наличие выражения с таким ID
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("id not found"))
			return
		}
		// и совпадение с ID пользователя
		if expression.UserId != user {
			w.WriteHeader(http.StatusNotAcceptable)
			w.Write([]byte("id for user not found"))
			return
		}

		ws.Mu.RUnlock()
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
