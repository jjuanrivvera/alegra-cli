package api

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestID_Unmarshal(t *testing.T) {
	cases := map[string]string{
		`"12"`:   "12",
		`12`:     "12",
		`12.0`:   "12",
		`12.5`:   "12.5",
		`null`:   "",
		`""`:     "",
		`"abc1"`: "abc1",
		// above 2^53: must not round-trip through float64
		`9007199254740993`:    "9007199254740993",
		`9223372036854775807`: "9223372036854775807",
	}
	for in, want := range cases {
		var id ID
		require.NoError(t, json.Unmarshal([]byte(in), &id), "input %s", in)
		assert.Equal(t, want, id.String(), "input %s", in)
	}
}

func TestID_Marshal(t *testing.T) {
	b, err := json.Marshal(ID("12"))
	require.NoError(t, err)
	assert.JSONEq(t, `"12"`, string(b))

	b, err = json.Marshal(ID(""))
	require.NoError(t, err)
	assert.Equal(t, "null", string(b))
}

func TestInt_RoundTrip(t *testing.T) {
	cases := map[string]int64{
		`30`:   30,
		`"30"`: 30,
		`""`:   0,
		`null`: 0,
		`12.9`: 12,
		// above 2^53: must not round-trip through float64
		`9007199254740993`:   9007199254740993,
		`"9007199254740993"`: 9007199254740993,
	}
	for in, want := range cases {
		var n Int
		require.NoError(t, json.Unmarshal([]byte(in), &n), "input %s", in)
		assert.Equal(t, want, int64(n), "input %s", in)
	}
	b, err := json.Marshal(Int(7))
	require.NoError(t, err)
	assert.Equal(t, "7", string(b))

	// ParseFloat accepts these; the decoder must not (int64(NaN) is undefined).
	for _, in := range []string{`"NaN"`, `"Inf"`, `"-Inf"`} {
		var n Int
		assert.Error(t, json.Unmarshal([]byte(in), &n), "input %s", in)
	}
}

func TestMoney_RoundTrip(t *testing.T) {
	cases := map[string]Money{
		`1000`:      "1000",
		`"1000.50"`: "1000.50",
		`""`:        "",
		`null`:      "",
		`59500`:     "59500",
		// Amounts beyond float64's ~15 significant digits must survive intact;
		// a float64-backed Money rewrote the trailing cents here.
		`1234567890123456.78`:   "1234567890123456.78",
		`"9999999999999999.99"`: "9999999999999999.99",
	}
	for in, want := range cases {
		var m Money
		require.NoError(t, json.Unmarshal([]byte(in), &m), "input %s", in)
		assert.Equal(t, want, m, "input %s", in)
	}
	b, err := json.Marshal(Money("1234.5"))
	require.NoError(t, err)
	assert.Equal(t, "1234.5", string(b))

	// The zero value (absent amount) still emits a valid number.
	b, err = json.Marshal(Money(""))
	require.NoError(t, err)
	assert.Equal(t, "0", string(b))

	// A Money cast from non-numeric text must fail loudly at marshal time
	// rather than emit invalid JSON.
	_, err = json.Marshal(Money("abc"))
	assert.Error(t, err)

	f, err := Money("1250.50").Float64()
	require.NoError(t, err)
	assert.InDelta(t, 1250.50, f, 0.0001)

	// ParseFloat accepts these; a NaN/Inf Money would poison any re-marshal
	// of the containing document (JSON cannot represent them).
	for _, in := range []string{`"NaN"`, `"nAn"`, `"Inf"`, `"-Infinity"`} {
		var m Money
		assert.Error(t, json.Unmarshal([]byte(in), &m), "input %s", in)
	}
}

func TestRefs_Unmarshal(t *testing.T) {
	var r Refs
	require.NoError(t, json.Unmarshal([]byte(`[{"id":1,"name":"A"},{"id":"2","name":"B"}]`), &r))
	assert.Equal(t, Refs{{ID: "1", Name: "A"}, {ID: "2", Name: "B"}}, r)

	// Alegra also serializes the field as a single object (e.g. costCenter on
	// income debit notes); it must decode as a one-element slice.
	require.NoError(t, json.Unmarshal([]byte(`{"id":7,"name":"Central"}`), &r))
	assert.Equal(t, Refs{{ID: "7", Name: "Central"}}, r)

	require.NoError(t, json.Unmarshal([]byte(`null`), &r))
	assert.Nil(t, r)

	assert.Error(t, json.Unmarshal([]byte(`"nope"`), &r))
}

func TestStringOrSlice_Unmarshal(t *testing.T) {
	var s StringOrSlice
	require.NoError(t, json.Unmarshal([]byte(`"client"`), &s))
	assert.Equal(t, StringOrSlice{"client"}, s)

	require.NoError(t, json.Unmarshal([]byte(`["client","provider"]`), &s))
	assert.Equal(t, StringOrSlice{"client", "provider"}, s)

	require.NoError(t, json.Unmarshal([]byte(`null`), &s))
	assert.Nil(t, s)
}

func TestRef_Decode(t *testing.T) {
	var r Ref
	require.NoError(t, json.Unmarshal([]byte(`{"id":7,"name":"IVA"}`), &r))
	assert.Equal(t, "7", r.ID.String())
	assert.Equal(t, "IVA", r.Name)
}

func TestID_UnmarshalInvalid(t *testing.T) {
	var id ID
	assert.Error(t, json.Unmarshal([]byte(`{`), &id))
}
