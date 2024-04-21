package registration

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/sv345922/arithmometer_v2/internal/entities"
	"github.com/sv345922/arithmometer_v2/internal/wSpace"
	"log"
	"net/http"
)

func Registration(ctx context.Context, ws *wSpace.WorkingSpace) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// Проверить что это запрос Post
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("требуется метод POST"))
			return
		}
		// Читаем тело запроса, в котором записан логин и пароль
		var data entities.UserData
		err := json.NewDecoder(r.Body).Decode(&data)
		defer r.Body.Close()
		if err != nil {
			log.Println("Переданы неправильные данные")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		ok, err := checkAndInsertUser(ctx, ws.DB, data)
		if err != nil {
			log.Printf("check and insert user failed: %v\n", err)
		}
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "пользователь с имененм %s зарегистрирован ранее", data.Name)
			return
		}
		w.WriteHeader(http.StatusCreated)
	}
}
