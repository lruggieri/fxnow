package util

import (
	"database/sql"
)

func SQLTimeToUnix(t sql.NullTime) int64 {
	if t.Valid {
		return t.Time.Unix()
	}

	return 0
}
