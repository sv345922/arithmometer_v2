package newExpression

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/sv345922/arithmometer_v2/internal/entities"
)

func updateTimings(ctx context.Context, tx *sql.Tx, t *entities.Timings) error {
	_, err := tx.ExecContext(ctx,
		`UPDATE timings SET plus = $1, minus = $2, mult = $3, div = $4  WHERE id = 1`,
		t.Plus, t.Minus, t.Mult, t.Div)
	if err != nil {
		return fmt.Errorf("updating timings: %w", err)
	}
	return nil
}
