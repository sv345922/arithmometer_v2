package dataBase

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/sv345922/arithmometer_v2/internal/configs"
	"github.com/sv345922/arithmometer_v2/internal/dataBase/query"
	"github.com/sv345922/arithmometer_v2/internal/entities"
)

// TODO
func CreateTables(ctx context.Context, db *sql.DB) error {

	if _, err := db.ExecContext(ctx, query.Q_Tasks); err != nil {
		return fmt.Errorf("Q_Tasks %w", err)
	}
	if _, err := db.ExecContext(ctx, query.Q_ReadyToCalc); err != nil {
		return fmt.Errorf("Q_ReadyToCalc %w", err)
	}
	if _, err := db.ExecContext(ctx, query.Q_Working); err != nil {
		return fmt.Errorf("Q_Working %w", err)
	}
	if _, err := db.ExecContext(ctx, query.Q_NotReady); err != nil {
		return fmt.Errorf("Q_NotReady %w", err)
	}
	if _, err := db.ExecContext(ctx, query.Expressions); err != nil {
		return fmt.Errorf("Expressions %w", err)
	}
	if _, err := db.ExecContext(ctx, query.AllNodes); err != nil {
		return fmt.Errorf("AllNodes %w", err)
	}
	if _, err := db.ExecContext(ctx, query.Timings); err != nil {
		return fmt.Errorf("Timings %w", err)
	}
	if _, err := db.ExecContext(ctx, query.Users); err != nil {
		return fmt.Errorf("Users %w", err)
	}
	return nil
}

// Загружает в пустую структуру DataBase  данные из sqlite
func (db *DataBase) Load(ctx context.Context, dataBase *sql.DB) error {
	// TODO таблицы очередей и user
	// заполняем поля
	expressions, err := GetExpressions(ctx, dataBase)
	if err != nil {
		return fmt.Errorf("GetExpressions %w", err)
	}
	db.Expressions = expressions
	allNodes, err := GetNodes(ctx, dataBase)
	if err != nil {
		return fmt.Errorf("GetNodes %w", err)
	}
	db.AllNodes = allNodes
	allTasks, err := GetTasks(ctx, dataBase)
	if err != nil {
		return fmt.Errorf("GetTasks %w", err)
	}
	db.Tasks = allTasks
	timings, err := GetTimings(ctx, dataBase)
	if err != nil {
		return fmt.Errorf("GetTimings %w", err)
	}
	db.Timings = &timings
	users, err := GetUsers(ctx, dataBase)
	if err != nil {
		return fmt.Errorf("GetUsers %w", err)
	}
	db.Users = users
	return nil
}

func GetUsers(ctx context.Context, db *sql.DB) ([]*entities.User, error) {
	var users []*entities.User
	var q = "SELECT * FROM users;"
	rows, err := db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		u := &entities.User{}
		err = rows.Scan(&u.ID, &u.Username, &u.Password)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

func GetTimings(ctx context.Context, db *sql.DB) (entities.Timings, error) {
	var timings entities.Timings
	var q = "SELECT plus, minus, mult, div FROM timings WHERE id = 1"
	err := db.QueryRowContext(ctx, q).Scan(&timings.Plus, &timings.Minus, &timings.Mult, &timings.Div)
	if err != nil {
		if err == sql.ErrNoRows {
			// log.Print("установлены по умолчанию тайминги БД")
			timings.Plus = configs.DefaultTimings
			timings.Minus = configs.DefaultTimings
			timings.Mult = configs.DefaultTimings
			timings.Div = configs.DefaultTimings
			var q = "INSERT INTO timings (plus, minus, mult, div) VALUES ($1, $1, $1, $1)"
			_, err = db.ExecContext(ctx, q, configs.DefaultTimings)
			if err != nil {
				return timings, err
			}
			return timings, nil
		} else {
			return timings, fmt.Errorf("%w", err)
		}
	}
	return timings, nil
}

func GetTasks(ctx context.Context, db *sql.DB) ([]*entities.Task, error) {
	var tasks []*entities.Task
	rows, err := db.QueryContext(ctx, query.SelectAllTasks)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		t := entities.Task{}
		err = rows.Scan(
			&t.NodeId,
			&t.X,
			&t.XReady,
			&t.Y,
			&t.YReady,
			&t.CalcId,
			&t.Deadline,
			&t.Duration,
		)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, &t)
	}
	return tasks, nil
}

func GetNodes(ctx context.Context, db *sql.DB) ([]*entities.Node, error) {
	var nodes []*entities.Node

	rows, err := db.QueryContext(ctx, query.SelectAllNodes)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		n := entities.Node{}
		err = rows.Scan(
			&n.Id,
			&n.ExpressionId,
			&n.Op,
			&n.X,
			&n.Y,
			&n.Val,
			&n.Sheet,
			&n.Calculated,
			&n.Parent,
		)
		if err != nil {
			return nil, err
		}
		nodes = append(nodes, &n)
	}
	return nodes, nil
}

func InsertNode(ctx context.Context, db *sql.DB, node *entities.Node) (uint64, error) {
	q := query.InsertNode
	result, err := db.ExecContext(ctx, q,
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
func InsertExpression(ctx context.Context,
	db *sql.DB,
	expression *entities.Expression,
) (uint64, error) {
	q := query.InsertExpression
	result, err := db.ExecContext(ctx, q,
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

func GetExpressions(ctx context.Context, db *sql.DB) ([]*entities.Expression, error) {
	var expressions []*entities.Expression

	rows, err := db.QueryContext(ctx, query.SelectAllExpressions)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		e := entities.Expression{}
		err = rows.Scan(
			&e.Id,
			&e.UserId,
			&e.UserTask,
			&e.ResultExpr,
			&e.Status,
			&e.RootId)
		if err != nil {
			return nil, err
		}
		expressions = append(expressions, &e)
	}
	return expressions, nil
}

// Функциия для тестирования доступности БД
func Testlocked(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, `
	INSERT INTO timings (plus, minus, mult, div) VALUES (5, 5, 5, 5)`)
	return err
}
