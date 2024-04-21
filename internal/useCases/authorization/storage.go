package authorization

import (
	"context"
	"database/sql"
	"fmt"
)

func SelectUserID(ctx context.Context, db *sql.DB, userName string) (string, error) {
	var userId int
	var q = "SELECT id FROM users WHERE userName = $1"
	err := db.QueryRowContext(ctx, q, userName).Scan(&userId)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%d", userId), nil
}
