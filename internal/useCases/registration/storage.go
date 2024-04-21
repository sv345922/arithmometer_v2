package registration

import (
	"context"
	"database/sql"
	"log"

	"github.com/sv345922/arithmometer_v2/internal/entities"
)

// Возвращаем true если пользователь добавлен, false если не добавлен
// по любым причинам, если при этом нет ошибки, значит пользователь уже есть в БД
func checkAndInsertUser(ctx context.Context, db *sql.DB, userData entities.UserData) (bool, error) {
	var user struct {
		id       int
		username string
		password string
	}
	var q = "SELECT * FROM users WHERE userName = $1"
	err := db.QueryRowContext(ctx, q, userData.Name).Scan(&user)
	switch {
	case err == sql.ErrNoRows:
		// Добавляем запись в БД
		q = "INSERT INTO users(username, password) VALUES($1, $2)"
		_, err := db.ExecContext(ctx, q, userData.Name, userData.Password)
		if err != nil {
			return false, err
		}
		return true, nil
	case err != nil:
		return false, err
	default:
		// имя занято
		log.Printf("username '%s' is busy", userData.Name)
		return false, nil
	}
}
