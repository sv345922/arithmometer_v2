package app

import (
	"context"
	"database/sql"
	"github.com/sv345922/arithmometer_v2/internal/wSpace"
)

// Создает рабочее пространство из сохраненной базы данных
func RunTasker(ctx context.Context, db *sql.DB) (*wSpace.WorkingSpace, error) {
	// Восстанавливаем выражения и задачи из базы данных
	// Загрузка сохраненной БД
	ws, err := wSpace.LoadDB(ctx, db)
	if err != nil {
		// log.Println("ошибка загрузки БД", err)
		return ws, err
	}

	return ws, err
}
