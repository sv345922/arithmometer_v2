package giveAnswer

import (
	"context"
	"database/sql"
	"fmt"
)

func DeleteNodeTask(ctx context.Context, db *sql.DB, id uint64) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("cannot begin transaction: %w", err)
	}
	_, err = tx.ExecContext(ctx,
		`DELETE FROM allNodes WHERE id = $1`, id)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("cannot delete node with id %d: %w", id, err)
	}
	_, err = tx.ExecContext(ctx,
		`DELETE FROM tasks WHERE nodeId = $1`, id)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("cannot delete task with id %d: %w", id, err)
	}
	tx.Commit()
	return nil
}

func DeleteNode(ctx context.Context, db *sql.DB, id uint64) error {
	_, err := db.ExecContext(ctx,
		`DELETE FROM allNodes WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("cannot delete node with id %d: %w", id, err)
	}
	return nil
}
