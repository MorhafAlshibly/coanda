package model

import (
	"context"
	"database/sql"
	"time"

	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/mysql"
)

var gq = goqu.Dialect("mysql")

type GetItemsParams struct {
	Type   sql.NullString `db:"type"`
	Limit  uint64
	Offset uint64
}

func filterGetItemsParams(arg GetItemsParams) goqu.Expression {
	expressions := []goqu.Expression{}
	if arg.Type.Valid {
		expressions = append(expressions, goqu.C("type").Eq(arg.Type))
	}
	return goqu.And(expressions...)
}

func (q *Queries) GetItems(ctx context.Context, arg GetItemsParams) ([]Item, error) {
	item := gq.Select("id", "type", "data", "expires_at", "created_at", "updated_at").From("item").Prepared(true)
	query, args, err := item.Where(filterGetItemsParams(arg),
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
	var items []Item
	for rows.Next() {
		var i Item
		if err := rows.Scan(
			&i.ID,
			&i.Type,
			&i.Data,
			&i.ExpiresAt,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	return items, nil
}
