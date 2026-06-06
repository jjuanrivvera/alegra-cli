package commands

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// parseDateExpr resolves a human date expression to a YYYY-MM-DD string,
// relative to now. Supported forms:
//
//	2026-06-01            absolute date (passed through)
//	today | yesterday | tomorrow
//	this-month | last-month | this-year | last-year | this-quarter
//	7d | 2w | 3m | 1y    N days/weeks/months/years ago
func parseDateExpr(expr string, now time.Time) (string, error) {
	e := strings.ToLower(strings.TrimSpace(expr))
	if e == "" {
		return "", fmt.Errorf("empty date expression")
	}
	const layout = "2006-01-02"

	// Absolute date.
	if t, err := time.Parse(layout, e); err == nil {
		return t.Format(layout), nil
	}

	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	switch e {
	case "today":
		return today.Format(layout), nil
	case "yesterday":
		return today.AddDate(0, 0, -1).Format(layout), nil
	case "tomorrow":
		return today.AddDate(0, 0, 1).Format(layout), nil
	case "this-month":
		return time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC).Format(layout), nil
	case "last-month":
		return time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC).AddDate(0, -1, 0).Format(layout), nil
	case "this-year":
		return time.Date(now.Year(), 1, 1, 0, 0, 0, 0, time.UTC).Format(layout), nil
	case "last-year":
		return time.Date(now.Year()-1, 1, 1, 0, 0, 0, 0, time.UTC).Format(layout), nil
	case "this-quarter":
		q := (int(now.Month()) - 1) / 3
		return time.Date(now.Year(), time.Month(q*3+1), 1, 0, 0, 0, 0, time.UTC).Format(layout), nil
	}

	// Relative: <N><unit>, e.g. 7d, 2w, 3m, 1y.
	if len(e) >= 2 {
		unit := e[len(e)-1]
		if n, err := strconv.Atoi(e[:len(e)-1]); err == nil && n >= 0 {
			switch unit {
			case 'd':
				return today.AddDate(0, 0, -n).Format(layout), nil
			case 'w':
				return today.AddDate(0, 0, -7*n).Format(layout), nil
			case 'm':
				return today.AddDate(0, -n, 0).Format(layout), nil
			case 'y':
				return today.AddDate(-n, 0, 0).Format(layout), nil
			}
		}
	}

	return "", fmt.Errorf("unrecognized date %q (use YYYY-MM-DD, today, yesterday, this-month, last-month, this-year, this-quarter, or Nd/Nw/Nm/Ny)", expr)
}
