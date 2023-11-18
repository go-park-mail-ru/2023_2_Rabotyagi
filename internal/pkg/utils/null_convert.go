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

func UnsafeStringToNull(str *string) sql.NullString {
	if str == nil {
		return sql.NullString{Valid: false} //nolint:exhaustruct
	}

	return sql.NullString{String: *str, Valid: true}
}
