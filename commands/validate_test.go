package commands

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func bodyMap(t *testing.T, raw string) map[string]any {
	t.Helper()
	m, ok := bodyToMap([]byte(raw))
	assert.True(t, ok)
	return m
}

func TestValidateContact(t *testing.T) {
	// valid
	assert.Empty(t, validateForCreate("contacts", "colombia",
		bodyMap(t, `{"name":"Acme","identification":{"type":"NIT","number":"1"}}`)))
	// missing name
	assert.NotEmpty(t, validateForCreate("contacts", "", bodyMap(t, `{}`)))
	// identification as string
	probs := validateForCreate("contacts", "colombia", bodyMap(t, `{"name":"X","identification":"123"}`))
	assert.Len(t, probs, 1)
	assert.Contains(t, probs[0], "must be an object")
}

func TestValidateSalesDocument(t *testing.T) {
	valid := `{"client":{"id":1},"date":"2026-06-06","items":[{"id":1,"price":1,"quantity":1}]}`
	assert.Empty(t, validateForCreate("invoices", "colombia", bodyMap(t, valid)))

	// missing everything
	probs := validateForCreate("invoices", "colombia", bodyMap(t, `{}`))
	assert.Len(t, probs, 3) // client, date, items

	// stamping without numberTemplate
	stamping := `{"client":{"id":1},"date":"2026-06-06","items":[{"id":1}],"stamp":{"generateStamp":true}}`
	probs = validateForCreate("invoices", "colombia", bodyMap(t, stamping))
	assert.Len(t, probs, 1)
	assert.Contains(t, probs[0], "numberTemplate.id")

	// stamping with numberTemplate → ok
	good := `{"client":{"id":1},"date":"2026-06-06","items":[{"id":1}],"numberTemplate":{"id":7},"stamp":{"generateStamp":true}}`
	assert.Empty(t, validateForCreate("invoices", "colombia", bodyMap(t, good)))
}

func TestValidate_UnknownResource(t *testing.T) {
	assert.Nil(t, validateForCreate("taxes", "colombia", bodyMap(t, `{}`)))
}
