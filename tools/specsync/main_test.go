package main

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseManifest_Fixture(t *testing.T) {
	body, err := os.ReadFile("testdata/llms.txt")
	require.NoError(t, err)

	man, err := parseManifest(body)
	require.NoError(t, err)

	// Modern REST slugs: method uppercased, "-id" expanded, duplicates dropped,
	// output sorted by path then method. get_bills stays REST — a single
	// underscore is not enough to call a method-prefixed slug an MCP tool.
	require.Len(t, man.RESTEndpoints, 6)
	assert.Equal(t, RESTEndpoint{Method: "GET", Path: "/bills", Slug: "get_bills"}, man.RESTEndpoints[0])
	assert.Equal(t, RESTEndpoint{Method: "GET", Path: "/taxes", Slug: "get_taxes"}, man.RESTEndpoints[1])
	assert.Equal(t, RESTEndpoint{Method: "DELETE", Path: "/taxes/{id}", Slug: "delete_taxes-id"}, man.RESTEndpoints[2])
	assert.Equal(t, RESTEndpoint{Method: "GET", Path: "/taxes/{id}", Slug: "get_taxes-id"}, man.RESTEndpoints[3])
	assert.Equal(t, RESTEndpoint{Method: "PUT", Path: "/taxes/{id}", Slug: "put_taxes-id"}, man.RESTEndpoints[4])
	assert.Equal(t, RESTEndpoint{Method: "POST", Path: "/webhooks-subscriptions", Slug: "post_webhooks-subscriptions"}, man.RESTEndpoints[5])

	// Both separators: legacy group__tool and the current group_tool form,
	// including a tool name that itself contains underscores.
	assert.Equal(t, []string{
		"accounting__createcostcenter",
		"accounting__listjournals",
		"expenses_list-bills",
		"items_getitems",
		"reports_get_cash_flow",
	}, man.MCPTools)
	assert.Equal(t, []string{"createAdditionalCharge", "listreconciliations"}, man.LegacySlugs)
	assert.Equal(t, []string{"accounting", "bills", "expenses", "items", "reports", "taxes", "webhooks-subscriptions"}, man.Resources)

	// Guide pages (docs/*.md) must never leak into the manifest.
	for _, e := range man.RESTEndpoints {
		assert.NotContains(t, e.Slug, "docs/")
	}
}

func TestParseManifest_LiveIndexShape(t *testing.T) {
	// The committed manifest came from a real llms.txt; if a cached copy is
	// present locally (specsync runs from the repo root), re-parse it as a
	// broader regression check.
	body, err := os.ReadFile("../../.alegra-spec/llms.txt")
	if err != nil {
		t.Skip("no cached .alegra-spec/llms.txt; run specsync to fetch one")
	}
	man, err := parseManifest(body)
	require.NoError(t, err)
	assert.Greater(t, len(man.RESTEndpoints), 50, "real index documents far more than 50 REST endpoints")
	assert.NotEmpty(t, man.MCPTools)
	assert.NotEmpty(t, man.Resources)
}

func TestParseManifest_EmptyInputFailsLoudly(t *testing.T) {
	// specsync feeds spec-check's drift detection: an upstream format change
	// must abort, not write an empty manifest that disables detection.
	for name, body := range map[string]string{
		"empty file": "",
		"no links":   "# Alegra API\n\nplain text only\n",
		"only mcp/legacy": "- [a](https://developer.alegra.com/reference/accounting__listjournals.md)\n" +
			"- [b](https://developer.alegra.com/reference/listreconciliations.md)\n",
	} {
		_, err := parseManifest([]byte(body))
		require.Error(t, err, name)
		assert.True(t, strings.Contains(err.Error(), "0 REST endpoints"), name)
	}
}

func TestParseManifest_ZeroMCPToolsFailsLoudly(t *testing.T) {
	// Alegra dropped the group__tool separator for group_tool. A parser that
	// only recognizes "__" records every tool as legacy and writes mcpTools
	// null, which silently drops MCP-only resources from drift detection.
	body := "- [a](https://developer.alegra.com/reference/get_bills.md)\n" +
		"- [b](https://developer.alegra.com/reference/listreconciliations.md)\n"
	_, err := parseManifest([]byte(body))
	require.Error(t, err)
	assert.Contains(t, err.Error(), "0 MCP tools")
}
