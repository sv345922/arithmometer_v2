package authorization

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/sv345922/arithmometer_v2/internal/jwToken"
	"log"
	"net/http"
)

// MiddleWare проверка токена и запись в заголовок имени пользователя
func Authorization(ctx context.Context, db *sql.DB, next http.HandlerFunc) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		header := r.Header
		jwt, ok := header["Authorization"]

		if ok != true {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, " Unauthorized ")
			return
		}
		// проверка jwt и авторизация пользователя
		userName, ok, err := jwToken.CheckJWT(jwt[0][7:])
		if err != nil || ok != true {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, " Unauthorized ")
			return
		}
		// проверка в БД и получение id пользователя
		fmt.Println("запрос в БД") // todo delete
		userId, err := SelectUserID(ctx, db, userName)
		if err != nil {
			log.Printf("Select User Error: %v\n", err)
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, " Unauthorized ")
			return
		}
		// Запись заголовка "X-username"
		r.Header.Add("X-username", userId)

		next.ServeHTTP(w, r)

	}
}

// MiddleWare проверка токена и запись в заголовок имени пользователя
func FakeAuthorization(ctx context.Context, db *sql.DB, next http.HandlerFunc) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		header := r.Header
		jwt, ok := header["Authorization"]
		fmt.Println(jwt) // todo delete
		if ok != true {

			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, " Unauthorized ")
		} else {
			userId := "0"
			// Запись заголовка "X-username"
			r.Header.Add("X-username", userId)
			next.ServeHTTP(w, r)
		}
	}
}
