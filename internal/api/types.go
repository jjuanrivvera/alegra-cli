package api

import (
	"bytes"
	"encoding/json"
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
		f, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return err
		}
		*i = Int(int64(f))
		return nil
	}
	var f float64
	if err := json.Unmarshal(data, &f); err != nil {
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

// Money is an amount Alegra may serialize as number or numeric string.
type Money float64

// UnmarshalJSON accepts number, numeric string, or null.
func (m *Money) UnmarshalJSON(data []byte) error {
	data = bytes.TrimSpace(data)
	if len(data) == 0 || string(data) == "null" {
		*m = 0
		return nil
	}
	if data[0] == '"' {
		var s string
		if err := json.Unmarshal(data, &s); err != nil {
			return err
		}
		if s == "" {
			*m = 0
			return nil
		}
		f, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return err
		}
		*m = Money(f)
		return nil
	}
	var f float64
	if err := json.Unmarshal(data, &f); err != nil {
		return err
	}
	*m = Money(f)
	return nil
}

// MarshalJSON emits a plain number.
func (m Money) MarshalJSON() ([]byte, error) {
	return json.Marshal(float64(m))
}
