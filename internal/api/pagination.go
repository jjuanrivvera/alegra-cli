package api

import (
	"net/url"
	"strconv"
)

// ListParams are the query parameters common to Alegra list endpoints.
//
// Alegra uses offset pagination via start/limit (limit max 30). Extra holds any
// resource-specific filters (e.g. status, date, client_id) so a single struct
// serves every list endpoint.
type ListParams struct {
	Start          int
	Limit          int
	OrderField     string
	OrderDirection string // ASC or DESC
	Query          string // free-text search (most resources support ?query=)
	Fields         string // comma-separated field projection
	Extra          url.Values
}

// values renders the params as a url.Values, applying defaults.
func (p ListParams) values(defaultLimit int) url.Values {
	v := url.Values{}
	for k, vals := range p.Extra {
		for _, val := range vals {
			v.Add(k, val)
		}
	}
	limit := p.Limit
	if limit <= 0 {
		limit = defaultLimit
	}
	if limit > MaxPageLimit {
		limit = MaxPageLimit
	}
	v.Set("start", strconv.Itoa(p.Start))
	v.Set("limit", strconv.Itoa(limit))
	if p.OrderField != "" {
		v.Set("order_field", p.OrderField)
	}
	if p.OrderDirection != "" {
		v.Set("order_direction", p.OrderDirection)
	}
	if p.Query != "" {
		v.Set("query", p.Query)
	}
	if p.Fields != "" {
		v.Set("fields", p.Fields)
	}
	return v
}

// effectiveLimit reports the page size that values() would apply.
func (p ListParams) effectiveLimit(defaultLimit int) int {
	limit := p.Limit
	if limit <= 0 {
		limit = defaultLimit
	}
	if limit > MaxPageLimit {
		limit = MaxPageLimit
	}
	return limit
}
