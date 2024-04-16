package taskQueue

import (
	"context"
	"database/sql"
	"github.com/sv345922/arithmometer_v2/internal/dataBase/query"
)

// Вставляет задачу со всеми полями в БД
func InsertTask(ctx context.Context, tx *sql.Tx, task *Task) error {
	q := query.InsertTask
	_, err := tx.ExecContext(ctx, q,
		task.Node.Id,
		task.X,
		task.XReady,
		task.Y,
		task.YReady,
		task.CalcId,
		task.Deadline,
		task.Duration,
	)
	if err != nil {
		return err
	}

	return nil
}
