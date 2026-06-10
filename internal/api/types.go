package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"regexp"
	"strconv"
)

// ID is a flexible identifier that decodes from either a JSON string ("12") or
// a JSON number (12), which Alegra uses inconsistently across endpoints. It
// always marshals back to a string and renders cleanly in tables.
type ID string

// UnmarshalJSON accepts string, number, or null.
func (i *ID) UnmarshalJSON(data []byte) error {
	data = bytes.TrimSpace(data)
	if len(data) == 0 || string(data) == "null" {
		*i = ""
		return nil
	}
	if data[0] == '"' {
		var s string
		if err := json.Unmarshal(data, &s); err != nil {
			return err
		}
		*i = ID(s)
		return nil
	}
	// number (possibly float); normalize to a clean integer string when whole.
	var n json.Number
	if err := json.Unmarshal(data, &n); err != nil {
		return err
	}
	// Int64 first: round-tripping through float64 would silently lose
	// precision for ids above 2^53.
	if v, err := n.Int64(); err == nil {
		*i = ID(strconv.FormatInt(v, 10))
		return nil
	}
	if f, err := n.Float64(); err == nil && f == float64(int64(f)) {
		*i = ID(strconv.FormatInt(int64(f), 10))
		return nil
	}
	*i = ID(n.String())
	return nil
}

// MarshalJSON emits the id as a string (Alegra accepts string ids on input).
func (i ID) MarshalJSON() ([]byte, error) {
	if i == "" {
		return []byte("null"), nil
	}
	return json.Marshal(string(i))
}

// String returns the raw identifier.
func (i ID) String() string { return string(i) }

// Int is a flexible integer that decodes from a JSON number or numeric string
// ("30"), which Alegra uses interchangeably for counts like days, decimal
// precision, and sequence numbers.
type Int int64

// UnmarshalJSON accepts number, numeric string, or null.
func (i *Int) UnmarshalJSON(data []byte) error {
	data = bytes.TrimSpace(data)
	if len(data) == 0 || string(data) == "null" {
		*i = 0
		return nil
	}
	if data[0] == '"' {
		var s string
		if err := json.Unmarshal(data, &s); err != nil {
			return err
		}
		if s == "" {
			*i = 0
			return nil
		}
		// ParseInt first: the float64 fallback (for "30.0"-style strings)
		// would lose precision above 2^53.
		if v, err := strconv.ParseInt(s, 10, 64); err == nil {
			*i = Int(v)
			return nil
		}
		f, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return err
		}
		// int64(NaN/Inf) is implementation-defined; reject instead.
		if math.IsNaN(f) || math.IsInf(f, 0) {
			return fmt.Errorf("invalid integer value %q", s)
		}
		*i = Int(int64(f))
		return nil
	}
	var n json.Number
	if err := json.Unmarshal(data, &n); err != nil {
		return err
	}
	if v, err := n.Int64(); err == nil {
		*i = Int(v)
		return nil
	}
	f, err := n.Float64()
	if err != nil {
		return err
	}
	*i = Int(int64(f))
	return nil
}

// MarshalJSON emits a plain integer.
func (i Int) MarshalJSON() ([]byte, error) {
	return json.Marshal(int64(i))
}

// Ref is the common {id, name} nested reference Alegra embeds for related
// records (seller, term, priceList, category, ...).
type Ref struct {
	ID   ID     `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

// Refs decodes a related-records field that Alegra serializes as either a
// single {id,name} object or an array of them (e.g. costCenter on income
// debit notes, documented in both shapes across operations).
type Refs []Ref

// UnmarshalJSON accepts an object, an array of objects, or null.
func (r *Refs) UnmarshalJSON(data []byte) error {
	data = bytes.TrimSpace(data)
	if len(data) == 0 || string(data) == "null" {
		*r = nil
		return nil
	}
	if data[0] == '[' {
		var arr []Ref
		if err := json.Unmarshal(data, &arr); err != nil {
			return err
		}
		*r = arr
		return nil
	}
	var one Ref
	if err := json.Unmarshal(data, &one); err != nil {
		return err
	}
	*r = []Ref{one}
	return nil
}

// StringOrSlice decodes a JSON value that Alegra serializes as either a single
// string or an array of strings (e.g. contact "type", which is "client" in
// simple mode but ["client","provider"] in advanced mode).
type StringOrSlice []string

// UnmarshalJSON accepts a string, an array of strings, or null.
func (s *StringOrSlice) UnmarshalJSON(data []byte) error {
	data = bytes.TrimSpace(data)
	if len(data) == 0 || string(data) == "null" {
		*s = nil
		return nil
	}
	if data[0] == '[' {
		var arr []string
		if err := json.Unmarshal(data, &arr); err != nil {
			return err
		}
		*s = arr
		return nil
	}
	var one string
	if err := json.Unmarshal(data, &one); err != nil {
		return err
	}
	*s = []string{one}
	return nil
}

// Money is an amount Alegra may serialize as number or numeric string. It
// stores the exact decimal text it was decoded from rather than a float64:
// accounting amounts must survive decode→encode untouched, and float64
// silently rewrites digits once an amount exceeds ~15 significant figures
// (realistic for COP/CLP ledger totals). The CLI never does money arithmetic,
// so lexical fidelity is the property that matters.
type Money string

// UnmarshalJSON accepts number, numeric string, or null.
func (m *Money) UnmarshalJSON(data []byte) error {
	data = bytes.TrimSpace(data)
	if len(data) == 0 || string(data) == "null" {
		*m = ""
		return nil
	}
	if data[0] == '"' {
		var s string
		if err := json.Unmarshal(data, &s); err != nil {
			return err
		}
		if s == "" {
			*m = ""
			return nil
		}
		norm, err := normalizeMoneyText(s)
		if err != nil {
			return err
		}
		*m = Money(norm)
		return nil
	}
	// The decoder already validated the token as a JSON number; keep its
	// exact bytes so re-marshaling cannot alter a single digit.
	var n json.Number
	if err := json.Unmarshal(data, &n); err != nil {
		return err
	}
	*m = Money(n.String())
	return nil
}

// normalizeMoneyText validates an amount string and returns text that is safe
// to emit as a bare JSON number. Text already in JSON number form passes
// through untouched (full fidelity). Anything else ParseFloat accepts but JSON
// can't represent — "NaN"/"Inf" (would poison any later re-marshal of the
// document), hex floats, leading "+" — is rejected or reformatted.
func normalizeMoneyText(s string) (string, error) {
	if jsonNumberRe.MatchString(s) {
		return s, nil
	}
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return "", fmt.Errorf("invalid money value %q", s)
	}
	if math.IsNaN(f) || math.IsInf(f, 0) {
		return "", fmt.Errorf("invalid money value %q", s)
	}
	return strconv.FormatFloat(f, 'g', -1, 64), nil
}

// jsonNumberRe is the JSON number grammar (RFC 8259 §6).
var jsonNumberRe = regexp.MustCompile(`^-?(0|[1-9][0-9]*)(\.[0-9]+)?([eE][-+]?[0-9]+)?$`)

// MarshalJSON emits a plain number; the zero value emits 0 (struct fields tag
// Money with omitempty, so an absent amount is normally dropped instead).
func (m Money) MarshalJSON() ([]byte, error) {
	if m == "" {
		return []byte("0"), nil
	}
	if !jsonNumberRe.MatchString(string(m)) {
		return nil, fmt.Errorf("invalid money value %q", string(m))
	}
	return []byte(m), nil
}

// Float64 returns the amount as a float64 for display or comparison. Never use
// it to compute ledger values that flow back to the API.
func (m Money) Float64() (float64, error) {
	if m == "" {
		return 0, nil
	}
	return strconv.ParseFloat(string(m), 64)
}

// String returns the decimal text; the zero value reads as "0".
func (m Money) String() string {
	if m == "" {
		return "0"
	}
	return string(m)
}
