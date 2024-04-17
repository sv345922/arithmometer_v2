package newExpression

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/sv345922/arithmometer_v2/internal/wSpace"
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
		newExpression, status, err := ProcessExpression(ctx, ws, &newClientExpression)
		// todo доработать при передаче только таймингов!
		w.WriteHeader(status)
		if err != nil {
			log.Println(err)
			w.Write([]byte(err.Error()))
			return
		}
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
