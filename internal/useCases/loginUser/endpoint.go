package loginUser

import (
	"context"
	"encoding/json"
	"github.com/sv345922/arithmometer_v2/internal/entities"
	"net/http"

	"github.com/sv345922/arithmometer_v2/internal/wSpace"
)

func Login(ctx context.Context, ws *wSpace.WorkingSpace) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// Проверить что это запрос Get
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write([]byte("требуется метод Get"))
			return
		}
		var data entities.UserData
		err := json.NewDecoder(r.Body).Decode(&data)
		defer r.Body.Close()
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		// Создать токен и вернуть его в теле ответа
		token, err := getAccess(ctx, ws, data)
		if err != nil {
			w.WriteHeader(http.StatusNotAcceptable)
			w.Write([]byte(err.Error()))
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(token))
	}
}
