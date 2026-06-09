package api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// Resource is a generic, typed handle to one Alegra REST resource collection
// (e.g. contacts, invoices). It centralizes list/get/create/update/delete and
// custom sub-actions so each concrete resource is a one-line accessor on Client
// plus its data type. T is the resource's record type.
type Resource[T any] struct {
	client *Client
	// path is the collection path segment, e.g. "contacts" or "credit-notes".
	path string
}

// NewResource constructs a typed resource handle. Concrete resources expose a
// Client accessor (e.g. c.Contacts()) that calls this.
func NewResource[T any](c *Client, path string) *Resource[T] {
	return &Resource[T]{client: c, path: strings.Trim(path, "/")}
}

// Path returns the collection path segment.
func (r *Resource[T]) Path() string { return r.path }

// List fetches a single page of records.
//
// Alegra is inconsistent about list response shapes: most endpoints return a
// bare JSON array, but some wrap it as {"data":[...]} or {"total":N,"results":
// [...]} (and {"metadata":...,"data":[...]} when metadata=true). decodeList
// normalizes all of these.
func (r *Resource[T]) List(ctx context.Context, params ListParams) ([]T, error) {
	q := params.values(r.client.defaultLimit)
	raw, err := r.client.do(ctx, http.MethodGet, r.path, q, nil)
	if err != nil {
		if errors.Is(err, errDryRun) {
			return nil, nil
		}
		return nil, err
	}
	return decodeList[T](raw)
}

// decodeList parses a list response that may be a bare array or an object that
// wraps the array under "data", "results", or "rows".
func decodeList[T any](raw []byte) ([]T, error) {
	trimmed := bytes.TrimSpace(raw)
	if len(trimmed) == 0 {
		return nil, nil
	}
	if trimmed[0] == '[' {
		var out []T
		if err := json.Unmarshal(trimmed, &out); err != nil {
			return nil, fmt.Errorf("alegra: decoding list response: %w", err)
		}
		return out, nil
	}
	var wrapper struct {
		Data          json.RawMessage `json:"data"`
		Results       json.RawMessage `json:"results"`
		Rows          json.RawMessage `json:"rows"`
		Subscriptions json.RawMessage `json:"subscriptions"` // /webhooks/subscriptions
	}
	if err := json.Unmarshal(trimmed, &wrapper); err == nil {
		for _, arr := range []json.RawMessage{wrapper.Data, wrapper.Results, wrapper.Rows, wrapper.Subscriptions} {
			if len(bytes.TrimSpace(arr)) == 0 {
				continue
			}
			var out []T
			if err := json.Unmarshal(arr, &out); err != nil {
				return nil, fmt.Errorf("alegra: decoding wrapped list response: %w", err)
			}
			return out, nil
		}
	}
	return nil, fmt.Errorf("alegra: unexpected list response shape (not an array or a {data|results|rows} wrapper)")
}

// ListAll auto-paginates, walking start/limit pages until a short page signals
// the end. It is bounded by maxPages to avoid runaway loops; pass 0 for the
// default cap.
func (r *Resource[T]) ListAll(ctx context.Context, params ListParams, maxPages int) ([]T, error) {
	if maxPages <= 0 {
		maxPages = 1000
	}
	limit := params.effectiveLimit(r.client.defaultLimit)
	var all []T
	start := params.Start
	for page := 0; page < maxPages; page++ {
		p := params
		p.Start = start
		p.Limit = limit
		batch, err := r.List(ctx, p)
		if err != nil {
			return all, err
		}
		all = append(all, batch...)
		if len(batch) < limit {
			return all, nil
		}
		start += limit
	}
	// Reached the page cap with a full final page: more records may exist.
	r.client.logger.Warn("alegra: --all stopped at the page cap; results may be truncated",
		"resource", r.path, "fetched", len(all), "max_pages", maxPages)
	return all, nil
}

// Count returns the total number of records matching params, using Alegra's
// metadata=true response wrapper ({"metadata":{"total":N},"data":[...]}).
// It requests a single record to minimize payload.
func (r *Resource[T]) Count(ctx context.Context, params ListParams) (int64, error) {
	p := params
	if p.Extra == nil {
		p.Extra = url.Values{}
	}
	p.Extra.Set("metadata", "true")
	p.Limit = 1
	q := p.values(r.client.defaultLimit)
	raw, err := r.client.do(ctx, http.MethodGet, r.path, q, nil)
	if err != nil {
		if errors.Is(err, errDryRun) {
			return 0, nil
		}
		return 0, err
	}
	// Total may live under metadata.total (e.g. contacts) or at the top level
	// (e.g. taxes: {"total":"4","results":[...]}).
	var wrapper struct {
		Total    Int `json:"total"`
		Metadata struct {
			Total Int `json:"total"`
		} `json:"metadata"`
	}
	if err := json.Unmarshal(bytes.TrimSpace(raw), &wrapper); err == nil {
		if wrapper.Metadata.Total > 0 {
			return int64(wrapper.Metadata.Total), nil
		}
		if wrapper.Total > 0 {
			return int64(wrapper.Total), nil
		}
	}
	// Fallback: the endpoint may not report a total; count the returned array.
	items, derr := decodeList[T](raw)
	if derr != nil {
		return 0, fmt.Errorf("alegra: endpoint does not report a total count")
	}
	return int64(len(items)), nil
}

// Get fetches a single record by ID.
func (r *Resource[T]) Get(ctx context.Context, id string) (*T, error) {
	var out T
	path := r.path + "/" + url.PathEscape(id)
	if err := r.client.doJSON(ctx, http.MethodGet, path, nil, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Create posts a new record. body may be a typed struct, a map, or
// json.RawMessage.
func (r *Resource[T]) Create(ctx context.Context, body any) (*T, error) {
	var out T
	if err := r.client.doJSON(ctx, http.MethodPost, r.path, nil, body, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Update edits an existing record by ID (HTTP PUT).
func (r *Resource[T]) Update(ctx context.Context, id string, body any) (*T, error) {
	var out T
	path := r.path + "/" + url.PathEscape(id)
	if err := r.client.doJSON(ctx, http.MethodPut, path, nil, body, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Delete removes a record by ID.
func (r *Resource[T]) Delete(ctx context.Context, id string) error {
	path := r.path + "/" + url.PathEscape(id)
	return r.client.doJSON(ctx, http.MethodDelete, path, nil, nil, nil)
}

// Action performs a custom POST sub-action on a record, e.g.
// POST /invoices/{id}/void or POST /invoices/{id}/email. The decoded response
// (if any) is unmarshaled into out when out is non-nil.
func (r *Resource[T]) Action(ctx context.Context, id, action string, body, out any) error {
	path := r.path + "/" + url.PathEscape(id) + "/" + strings.Trim(action, "/")
	return r.client.doJSON(ctx, http.MethodPost, path, nil, body, out)
}

// CollectionAction performs a custom POST on the collection itself, e.g.
// POST /invoices/preview or POST /bills/import-by-cufe.
func (r *Resource[T]) CollectionAction(ctx context.Context, action string, body, out any) error {
	path := r.path + "/" + strings.Trim(action, "/")
	return r.client.doJSON(ctx, http.MethodPost, path, nil, body, out)
}

// Raw exposes the underlying client for endpoints that do not fit CRUD (e.g.
// singletons like /company, or report endpoints under /reports/...).
func (r *Resource[T]) Raw() *Client { return r.client }

// GetInto fetches an arbitrary GET path under the API root into out. Useful for
// singleton and report resources.
func (c *Client) GetInto(ctx context.Context, path string, query url.Values, out any) error {
	return c.doJSON(ctx, http.MethodGet, path, query, nil, out)
}

// PostInto sends an arbitrary POST to a path under the API root.
func (c *Client) PostInto(ctx context.Context, path string, body, out any) error {
	return c.doJSON(ctx, http.MethodPost, path, nil, body, out)
}

// PutInto sends an arbitrary PUT to a path under the API root.
func (c *Client) PutInto(ctx context.Context, path string, body, out any) error {
	return c.doJSON(ctx, http.MethodPut, path, nil, body, out)
}

// RawMessage is re-exported for resource files that accept opaque JSON bodies.
type RawMessage = json.RawMessage
