package model

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/mysql"
)

var gq = goqu.Dialect("mysql")

type GetRecordParams struct {
	Id     sql.NullInt64 `db:"id"`
	Name   sql.NullString
	UserID sql.NullInt64 `db:"user_id"`
}

func filterGetRecordParams(arg GetRecordParams) goqu.Ex {
	expressions := goqu.Ex{}
	if arg.Id.Valid {
		expressions["id"] = arg.Id
	}
	if arg.Name.Valid {
		expressions["name"] = arg.Name
	}
	if arg.UserID.Valid {
		expressions["user_id"] = arg.UserID
	}
	return expressions
}

func (q *Queries) GetRecord(ctx context.Context, arg GetRecordParams) (RankedRecord, error) {
	record := gq.From("ranked_record").Prepared(true).Select(
		"id",
		"name",
		"user_id",
		"record",
		"ranking",
		"data",
		"created_at",
		"updated_at",
	)
	query, args, err := record.Where(filterGetRecordParams(arg)).Limit(1).ToSQL()
	if err != nil {
		return RankedRecord{}, err
	}
	var i RankedRecord
	err = q.db.QueryRowContext(ctx, query, args...).Scan(
		&i.ID,
		&i.Name,
		&i.UserID,
		&i.Record,
		&i.Ranking,
		&i.Data,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

type GetRecordsParams struct {
	Name   sql.NullString `db:"name"`
	UserId sql.NullInt64  `db:"user_id"`
	Limit  uint64
	Offset uint64
}

func filterGetRecordsParams(arg GetRecordsParams) goqu.Ex {
	expressions := goqu.Ex{}
	if arg.Name.Valid {
		expressions["name"] = arg.Name
	}
	if arg.UserId.Valid {
		expressions["user_id"] = arg.UserId
	}
	return expressions
}

func (q *Queries) GetRecords(ctx context.Context, arg GetRecordsParams) ([]RankedRecord, error) {
	records := gq.From("ranked_record").Prepared(true).Select(
		"id",
		"name",
		"user_id",
		"record",
		"ranking",
		"data",
		"created_at",
		"updated_at",
	)
	query, args, err := records.Where(filterGetRecordsParams(arg)).Limit(uint(arg.Limit)).Offset(uint(arg.Offset)).ToSQL()
	if err != nil {
		return nil, err
	}
	rows, err := q.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []RankedRecord
	for rows.Next() {
		var i RankedRecord
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.UserID,
			&i.Record,
			&i.Ranking,
			&i.Data,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	return items, nil
}

type UpdateRecordParams struct {
	GetRecordParams GetRecordParams
	Record          sql.NullInt64   `db:"record"`
	Data            json.RawMessage `db:"data"`
}

func (q *Queries) UpdateRecord(ctx context.Context, arg UpdateRecordParams) (sql.Result, error) {
	record := gq.Update("record").Prepared(true)
	var updates = goqu.Record{}
	if arg.Record.Valid {
		updates["record"] = arg.Record.Int64
	}
	if arg.Data != nil {
		updates["data"] = []byte(arg.Data)
	}
	query, args, err := record.Where(filterGetRecordParams(arg.GetRecordParams)).Set(updates).ToSQL()
	if err != nil {
		return nil, err
	}
	return q.db.ExecContext(ctx, query, args...)
}

func (q *Queries) DeleteRecord(ctx context.Context, arg GetRecordParams) (sql.Result, error) {
	record := gq.Delete("record").Prepared(true)
	query, args, err := record.Where(filterGetRecordParams(arg)).Limit(1).ToSQL()
	if err != nil {
		return nil, err
	}
	return q.db.ExecContext(ctx, query, args...)
}
