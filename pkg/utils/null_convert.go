package utils

import (
	"database/sql"
	"time"
)

func NullStringToUnsafe(nullStr sql.NullString) *string {
	if nullStr.Valid {
		return &nullStr.String
	}

	return nil
}

func UnsafeStringToNull(str *string) sql.NullString {
	if str == nil {
		return sql.NullString{Valid: false} //nolint:exhaustruct
	}

	return sql.NullString{String: *str, Valid: true}
}

func NullTimeToUnsafe(nullTime sql.NullTime) *time.Time {
	if nullTime.Valid {
		return &nullTime.Time
	}

	return nil
}

func UnsafeTimeToNull(time *time.Time) sql.NullTime {
	if time == nil {
		return sql.NullTime{Valid: false} //nolint:exhaustruct
	}

	return sql.NullTime{Time: *time, Valid: true}
}

func NullInt64ToUnsafeUint(nullInt64 sql.NullInt64) *uint64 {
	if nullInt64.Valid {
		innerUint64 := uint64(nullInt64.Int64)

		return &(innerUint64)
	}

	return nil
}

func UnsafeUint64ToNullInt(unsafeInt64 *uint64) sql.NullInt64 {
	if unsafeInt64 == nil {
		return sql.NullInt64{Valid: false} //nolint:exhaustruct
	}

	return sql.NullInt64{Int64: int64(*unsafeInt64), Valid: true}
}

func NullFloat64ToUnsafeFloat(nullFloat64 sql.NullFloat64) *float64 {
	if nullFloat64.Valid {
		innerFloat64 := nullFloat64.Float64

		return &(innerFloat64)
	}

	return nil
}

func UnsafeFloat64ToNullFloat(unsafeFloat64 *float64) sql.NullFloat64 {
	if unsafeFloat64 == nil {
		return sql.NullFloat64{Valid: false} //nolint:exhaustruct
	}

	return sql.NullFloat64{Float64: *unsafeFloat64, Valid: true}
}
