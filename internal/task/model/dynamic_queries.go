package model

import (
	"context"
	"database/sql"
	"time"

	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/mysql"
)

var gq = goqu.Dialect("mysql")

type GetTasksParams struct {
	Type      sql.NullString `db:"type"`
	Completed sql.NullBool   `db:"completed"`
	Limit     uint64
	Offset    uint64
}

func filterGetTasksParams(arg GetTasksParams) goqu.Ex {
	expressions := goqu.Ex{}
	if arg.Type.Valid {
		expressions["type"] = arg.Type
	}
	if arg.Completed.Valid {
		condition := goqu.Or(
			goqu.Ex{"completed_at": nil},
			goqu.Ex{"completed_at": goqu.Op{"<": time.Now()}},
		)
		if arg.Completed.Bool {
			expressions["completed_at"] = goqu.Op{"not": condition}
		} else {
			expressions["completed_at"] = condition
		}
	}
	return expressions
}

func (q *Queries) GetTasks(ctx context.Context, arg GetTasksParams) ([]Task, error) {
	task := gq.From("task").Prepared(true).Select(
		"id",
		"type",
		"data",
		"expires_at",
		"completed_at",
		"created_at",
		"updated_at",
	)
	query, args, err := task.Where(filterGetTasksParams(arg),
		goqu.Or(
			goqu.C("expires_at").IsNull(),
			goqu.C("expires_at").Gt(time.Now()),
		),
	).Limit(uint(arg.Limit)).Offset(uint(arg.Offset)).ToSQL()
	if err != nil {
		return nil, err
	}
	rows, err := q.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Task
	for rows.Next() {
		var i Task
		if err := rows.Scan(
			&i.ID,
			&i.Type,
			&i.Data,
			&i.ExpiresAt,
			&i.CompletedAt,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	return items, nil
}
