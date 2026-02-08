package service

import (
	"time"
)

// ShouldResetTask checks if tasks should be reset for a new day
// Returns true if the last reset was on a different day
func ShouldResetTask(last time.Time) bool {
	now := time.Now()
	return now.Format("2006-01-02") != last.Format("2006-01-02")
}

// DayKey returns the current day name for streak tracking
// Maps time.Weekday to lowercase day names used in Streak model
func DayKey() string {
	switch time.Now().Weekday() {
	case time.Monday:
		return "mon"
	case time.Tuesday:
		return "tue"
	case time.Wednesday:
		return "wed"
	case time.Thursday:
		return "thu"
	case time.Friday:
		return "fri"
	case time.Saturday:
		return "sat"
	case time.Sunday:
		return "sun"
	default:
		return ""
	}
}
