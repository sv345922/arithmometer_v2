package getTask

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/sv345922/arithmometer_v2/internal/taskQueue"
)

func UpdateInGetTask(ctx context.Context, db *sql.DB, task *taskQueue.Task) error {
	_, err := db.ExecContext(ctx,
		`UPDATE tasks SET calcID = $1, deadline = $2, duration = $3 WHERE nodeId = $4`,
		task.CalcId, task.Deadline, task.Duration, task.Node.Id)
	if err != nil {
		return fmt.Errorf("UpdateInGetTask: %w", err)
	}
	return nil
}
