package wSpace

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/sv345922/arithmometer_v2/internal/dataBase/query"
	"github.com/sv345922/arithmometer_v2/internal/entities"
)

func InsertNode(ctx context.Context, tx *sql.Tx, node *entities.Node) (uint64, error) {
	q := query.InsertNode
	result, err := tx.ExecContext(ctx, q,
		node.ExpressionId,
		node.Op,
		node.X,
		node.Y,
		node.Val,
		node.Sheet,
		node.Calculated,
		node.Parent,
	)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	node.Id = uint64(id)
	return uint64(id), nil
}

func UpdateNode(ctx context.Context, tx *sql.Tx, node *entities.Node) error {
	var q = `UPDATE allNodes SET x = $1, y = $2, parent = $3 WHERE id = $4;`
	_, err := tx.ExecContext(ctx, q, node.X, node.Y, node.Parent, node.Id)
	if err != nil {
		return fmt.Errorf("node %d: %w", node.Id, err)
	}
	return nil
}
func InsertExpression(ctx context.Context,
	tx *sql.Tx,
	expression *entities.Expression,
) (uint64, error) {
	q := query.InsertExpression

	// записать выражение в БД
	result, err := tx.ExecContext(ctx, q,
		expression.UserId,
		expression.UserTask,
		expression.ResultExpr,
		expression.Status,
		expression.RootId,
	)

	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	expression.Id = uint64(id)
	return uint64(id), nil
}
