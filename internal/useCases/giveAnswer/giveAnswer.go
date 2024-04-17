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
		status, err := AcceptAnswer(ctx, ws, container)
		if err != nil {
			log.Println(err)
		}

		w.WriteHeader(status)
	}
}
