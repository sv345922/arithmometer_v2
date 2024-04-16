package giveAnswer

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/sv345922/arithmometer_v2/internal/entities"
)

// Удаляет узел и соответсвующую задачу из БД
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

// Обновляет в таблице Expressions поля вычисленного выражения:
// результат, статус, удаляет ссылку на корневой узел
func UpdateExpressionResult(ctx context.Context, db *sql.DB, expr *entities.Expression) error {
	_, err := db.ExecContext(ctx,
		`UPDATE expressions SET resultExpr = $1, status = $2, rootId = 0 WHERE id = $3`,
		expr.ResultExpr,
		expr.Status,
		expr.Id)
	if err != nil {
		return fmt.Errorf("cannot update expression result: %w", err)
	}
	return nil
}
