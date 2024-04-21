package loginUser

import (
	"context"
	"database/sql"
)

func findUserPass(ctx context.Context, db *sql.DB, username string) (string, error) {
	var password string
	var q = "SELECT password FROM users WHERE username=$1;"
	err := db.QueryRowContext(ctx, q, username).Scan(&password)
	if err != nil {
		return "", err
	}
	return password, nil
}
