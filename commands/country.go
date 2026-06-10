package commands

import (
	"context"
	"strings"

	"github.com/jjuanrivvera/alegra-cli/internal/api"
	"github.com/jjuanrivvera/alegra-cli/internal/config"
)

// detectCountry reads the account's localization from GET /company. The Alegra
// API is a single endpoint whose behavior is localized per country, exposed as
// the company's applicationVersion (e.g. "colombia", "mexico", "costaRica").
// Returns the lowercased value, or "" if it can't be determined.
func detectCountry(ctx context.Context, client *api.Client) string {
	var company map[string]any
	if err := client.GetInto(ctx, "company", nil, &company); err != nil {
		return ""
	}
	if v, ok := company["applicationVersion"].(string); ok {
		return strings.ToLower(strings.TrimSpace(v))
	}
	return ""
}

// cacheCountry persists the detected country on the profile so later commands
// (validation) know the version without another /company round-trip. Best
// effort: a save error is non-fatal and silently ignored.
func cacheCountry(cfg *config.Config, profile, country string) {
	if country == "" {
		return
	}
	p := cfg.Profile(profile)
	if p.Country == country {
		return // already current; avoid a needless write
	}
	p.Country = country
	cfg.SetProfile(p)
	// Best-effort by design: the cache only saves a future detection call, so
	// a failed write must never break the command that triggered it.
	_ = cfg.Save()
}
