package conversion

import (
	"database/sql"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func StringToSqlNullString(s *string) sql.NullString {
	if s == nil {
		return sql.NullString{
			String: "",
			Valid:  false,
		}
	}
	return sql.NullString{
		String: *s,
		Valid:  true,
	}
}

func Int64ToSqlNullInt64(i *int64) sql.NullInt64 {
	if i == nil {
		return sql.NullInt64{
			Int64: 0,
			Valid: false,
		}
	}
	return sql.NullInt64{
		Int64: *i,
		Valid: true,
	}
}

func Uint64ToSqlNullInt64(i *uint64) sql.NullInt64 {
	if i == nil {
		return sql.NullInt64{
			Int64: 0,
			Valid: false,
		}
	}
	return sql.NullInt64{
		Int64: int64(*i),
		Valid: true,
	}
}

func Uint32ToSqlNullInt32(i *uint32) sql.NullInt32 {
	if i == nil {
		return sql.NullInt32{
			Int32: 0,
			Valid: false,
		}
	}
	return sql.NullInt32{
		Int32: int32(*i),
		Valid: true,
	}
}

func BoolToSqlNullBool(b *bool) sql.NullBool {
	if b == nil {
		return sql.NullBool{
			Bool:  false,
			Valid: false,
		}
	}
	return sql.NullBool{
		Bool:  *b,
		Valid: true,
	}
}

func Float64ToSqlNullFloat64(f *float64) sql.NullFloat64 {
	if f == nil {
		return sql.NullFloat64{
			Float64: 0,
			Valid:   false,
		}
	}
	return sql.NullFloat64{
		Float64: *f,
		Valid:   true,
	}
}

func TimeToSqlNullTime(t *time.Time) sql.NullTime {
	if t == nil {
		return sql.NullTime{
			Time:  time.Time{},
			Valid: false,
		}
	}
	return sql.NullTime{
		Time:  *t,
		Valid: true,
	}
}

func TimestampToSqlNullTime(t *timestamppb.Timestamp) sql.NullTime {
	if t == nil {
		return sql.NullTime{
			Time:  time.Time{},
			Valid: false,
		}
	}
	return sql.NullTime{
		Time:  t.AsTime(),
		Valid: true,
	}
}

func SqlNullInt64ToInt64(i sql.NullInt64) *int64 {
	if i.Valid {
		return &i.Int64
	}
	return nil
}

func SqlNullInt64ToUint64(i sql.NullInt64) *uint64 {
	if i.Valid {
		value := uint64(i.Int64)
		return &value
	}
	return nil
}

func SqlNullStringToString(s sql.NullString) *string {
	if s.Valid {
		return &s.String
	}
	return nil
}

func SqlNullTimeToTimestamp(t sql.NullTime) *timestamppb.Timestamp {
	if t.Valid {
		return TimeToTimestamppb(&t.Time)
	}
	return nil
}
