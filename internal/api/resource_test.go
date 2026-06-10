package api

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDecodeList_Shapes(t *testing.T) {
	type row struct {
		ID   ID     `json:"id"`
		Name string `json:"name"`
	}

	cases := map[string]string{
		"bare array":            `[{"id":"1","name":"a"},{"id":"2","name":"b"}]`,
		"data wrapper":          `{"metadata":{"total":2},"data":[{"id":"1","name":"a"},{"id":"2","name":"b"}]}`,
		"results wrapper":       `{"total":"2","results":[{"id":"1","name":"a"},{"id":"2","name":"b"}]}`,
		"rows wrapper":          `{"rows":[{"id":"1","name":"a"},{"id":"2","name":"b"}]}`,
		"subscriptions wrapper": `{"subscriptions":[{"id":"1","name":"a"},{"id":"2","name":"b"}]}`,
	}
	for name, body := range cases {
		t.Run(name, func(t *testing.T) {
			out, err := decodeList[row]([]byte(body))
			require.NoError(t, err)
			require.Len(t, out, 2)
			assert.Equal(t, ID("1"), out[0].ID)
			assert.Equal(t, "b", out[1].Name)
		})
	}
}

func TestDecodeList_Empty(t *testing.T) {
	out, err := decodeList[Contact](nil)
	require.NoError(t, err)
	assert.Empty(t, out)

	// Empty wrapper arrays must yield an empty slice, not an error or a
	// fall-through to the bare-array parser.
	for _, body := range []string{`{"data":[]}`, `{"results":[]}`, `{"rows":[]}`, `{"data":null}`} {
		out, err := decodeList[Contact]([]byte(body))
		require.NoError(t, err, "body %s", body)
		assert.Empty(t, out, "body %s", body)
	}
}

func TestCount_TotalShapes(t *testing.T) {
	cases := map[string]string{
		"metadata.total":  `{"metadata":{"total":42},"data":[{"id":"1"}]}`,
		"top-level total": `{"total":"7","results":[{"id":"1"}]}`,
	}
	want := map[string]int64{"metadata.total": 42, "top-level total": 7}
	for name, body := range cases {
		t.Run(name, func(t *testing.T) {
			c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "true", r.URL.Query().Get("metadata"))
				_, _ = w.Write([]byte(body))
			})
			n, err := NewResource[Contact](c, "x").Count(context.Background(), ListParams{})
			require.NoError(t, err)
			assert.Equal(t, want[name], n)
		})
	}
}

// TestList_ResultsWrapper exercises List end-to-end against a {total,results}
// response (the shape /taxes returns).
func TestList_ResultsWrapper(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`{"total":"1","results":[{"id":"1","name":"IVA"}]}`))
	})
	res := NewResource[Contact](c, "taxes")
	items, err := res.List(context.Background(), ListParams{})
	require.NoError(t, err)
	require.Len(t, items, 1)
	assert.Equal(t, "IVA", items[0].Name)
}
