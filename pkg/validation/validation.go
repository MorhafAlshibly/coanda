package validation

import (
	"database/sql"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func ValidateMaxPageLength(max *uint32, defaultMaxPageLength uint8, maxMaxPageLength uint8) uint8 {
	if max == nil {
		return defaultMaxPageLength
	}
	if *max > uint32(maxMaxPageLength) {
		return maxMaxPageLength
	}
	return uint8(*max)
}

func ValidateAnSqlNullString(s *string) sql.NullString {
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

func ValidateAnSqlNullInt64(i *int64) sql.NullInt64 {
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

func ValidateAUint64ToSqlNullInt64(i *uint64) sql.NullInt64 {
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

func ValidateAnSqlNullBool(b *bool) sql.NullBool {
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

func ValidateAnSqlNullFloat64(f *float64) sql.NullFloat64 {
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

func ValidateAnSqlNullTime(t *time.Time) sql.NullTime {
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

func ValidateATimestampToSqlNullTime(t *timestamppb.Timestamp) sql.NullTime {
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
