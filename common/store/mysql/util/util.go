package util

import (
	"database/sql"
	"time"
)

func SQLTimeToUnix(t sql.NullTime) int64 {
	if t.Valid {
		return t.Time.Unix()
	}

	return time.Time{}.Unix()
}
