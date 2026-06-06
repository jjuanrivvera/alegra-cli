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
		"bare array":      `[{"id":"1","name":"a"},{"id":"2","name":"b"}]`,
		"data wrapper":    `{"metadata":{"total":2},"data":[{"id":"1","name":"a"},{"id":"2","name":"b"}]}`,
		"results wrapper": `{"total":"2","results":[{"id":"1","name":"a"},{"id":"2","name":"b"}]}`,
		"rows wrapper":    `{"rows":[{"id":"1","name":"a"},{"id":"2","name":"b"}]}`,
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
