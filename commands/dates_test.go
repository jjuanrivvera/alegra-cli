package commands

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseDateExpr(t *testing.T) {
	now := time.Date(2026, 6, 15, 10, 0, 0, 0, time.UTC)
	cases := map[string]string{
		"2026-01-02":   "2026-01-02",
		"today":        "2026-06-15",
		"yesterday":    "2026-06-14",
		"tomorrow":     "2026-06-16",
		"this-month":   "2026-06-01",
		"last-month":   "2026-05-01",
		"this-year":    "2026-01-01",
		"last-year":    "2025-01-01",
		"this-quarter": "2026-04-01",
		"7d":           "2026-06-08",
		"2w":           "2026-06-01",
		"3m":           "2026-03-15",
		"1y":           "2025-06-15",
	}
	for expr, want := range cases {
		t.Run(expr, func(t *testing.T) {
			got, err := parseDateExpr(expr, now)
			require.NoError(t, err)
			assert.Equal(t, want, got)
		})
	}
}

func TestParseDateExpr_Invalid(t *testing.T) {
	now := time.Date(2026, 6, 15, 0, 0, 0, 0, time.UTC)
	for _, bad := range []string{"", "nonsense", "5x", "2026/01/01"} {
		_, err := parseDateExpr(bad, now)
		assert.Error(t, err, bad)
	}
}
