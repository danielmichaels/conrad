package repository

import "time"

const SQLiteCurrentTimestamp = "2006-01-02 15:04:05"

func NextRefreshTimestamp(rr int64) string {
	n := time.Now().UTC()
	nt := n.Add(time.Duration(rr) * time.Minute)
	return nt.Format(SQLiteCurrentTimestamp)
}
