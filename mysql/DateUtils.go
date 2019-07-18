package mysql

import (
	"github.com/go-sql-driver/mysql"
	"time"
)

func GetMaxTime(t1 mysql.NullTime, t2 mysql.NullTime, t3 mysql.NullTime, defaultTime time.Time) time.Time {
	max12 := GetMaxTime2(t1, t2)
	max3 := GetMaxTime2(max12, t3)
	if max3.Valid {
		return max3.Time
	}
	return defaultTime
}

func GetMaxTime2(t1 mysql.NullTime, t2 mysql.NullTime) mysql.NullTime {
	if t1.Valid && t2.Valid {
		if t1.Time.After(t2.Time) {
			return t1
		}
		return t2
	}
	if t1.Valid {
		return t1
	}
	if t2.Valid {
		return t2
	}
	return mysql.NullTime{Valid: false}
}

