package commands

import (
	"encoding/json"
	"fmt"
	"strings"
)

// validateForCreate runs conservative, client-side pre-flight checks on a
// request body before it's sent, returning human-readable problems. It only
// flags issues that would almost certainly fail server-side, so false positives
// stay rare. country is lowercased (e.g. "colombia"); "" means country-agnostic.
func validateForCreate(resource, country string, body map[string]any) []string {
	switch resource {
	case "contacts":
		return validateContact(body)
	case "invoices", "credit-notes", "estimates":
		return validateSalesDocument(body)
	default:
		return nil
	}
}

func validateContact(body map[string]any) []string {
	var problems []string
	if isBlank(body["name"]) {
		problems = append(problems, "name is required")
	}
	if id, ok := body["identification"]; ok {
		if _, isStr := id.(string); isStr {
			problems = append(problems,
				`identification must be an object {type, number}, e.g. --set 'identification={"type":"NIT","number":"901123456"}'`)
		}
	}
	return problems
}

func validateSalesDocument(body map[string]any) []string {
	var problems []string
	if isBlank(body["client"]) {
		problems = append(problems, "client is required (e.g. --set 'client={\"id\":12}')")
	}
	if isBlank(body["date"]) {
		problems = append(problems, "date is required (YYYY-MM-DD)")
	}
	if items, ok := body["items"].([]any); !ok || len(items) == 0 {
		problems = append(problems, "at least one line item is required (items[])")
	}
	// Electronic emission always needs a numbering resolution.
	if truthy(getPath(body, "stamp.generateStamp")) && isBlank(getPath(body, "numberTemplate.id")) {
		problems = append(problems,
			"electronic emission (stamp) requires numberTemplate.id — see `alegra number-templates list`")
	}
	return problems
}

// --- helpers ---

func isBlank(v any) bool {
	switch t := v.(type) {
	case nil:
		return true
	case string:
		return strings.TrimSpace(t) == ""
	case map[string]any:
		return len(t) == 0
	case []any:
		return len(t) == 0
	default:
		return false
	}
}

func truthy(v any) bool {
	switch t := v.(type) {
	case bool:
		return t
	case string:
		return t == "true" || t == "1"
	case float64:
		return t != 0
	default:
		return false
	}
}

func getPath(body map[string]any, path string) any {
	parts := strings.Split(path, ".")
	var cur any = body
	for _, p := range parts {
		m, ok := cur.(map[string]any)
		if !ok {
			return nil
		}
		cur = m[p]
	}
	return cur
}

// bodyToMap unmarshals a JSON body into a map for validation. Non-object bodies
// (e.g. a JSON array) return ok=false and are skipped by validation.
func bodyToMap(raw []byte) (map[string]any, bool) {
	var m map[string]any
	if err := json.Unmarshal(raw, &m); err != nil {
		return nil, false
	}
	return m, true
}

// formatValidationError renders problems as a single actionable error.
func formatValidationError(resource, country string, problems []string) error {
	ctx := resource
	if country != "" {
		ctx = fmt.Sprintf("%s, country: %s", resource, country)
	}
	return fmt.Errorf("pre-flight validation failed (%s):\n  - %s\n(use --no-validate to send anyway)",
		ctx, strings.Join(problems, "\n  - "))
}
