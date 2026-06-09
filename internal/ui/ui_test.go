package ui

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestColorDisabledForNonTerminal(t *testing.T) {
	// A bytes.Buffer is not a terminal → no color.
	p := For(&bytes.Buffer{})
	assert.False(t, p.Enabled())
	assert.Equal(t, "ok", p.Green("ok"))
	assert.Equal(t, "bad", p.Red("bad"))
}

func TestNoColorEnv(t *testing.T) {
	t.Setenv("NO_COLOR", "1")
	assert.False(t, colorEnabled(&bytes.Buffer{}))
}

func TestSetNoColor(t *testing.T) {
	defer SetNoColor(false)
	SetNoColor(true)
	assert.False(t, colorEnabled(&bytes.Buffer{}))
}

func TestWrapWhenEnabled(t *testing.T) {
	// Construct an enabled palette directly to exercise the styling codes.
	p := Palette{enabled: true}
	for name, styled := range map[string]string{
		"green":  p.Green("x"),
		"red":    p.Red("x"),
		"yellow": p.Yellow("x"),
		"cyan":   p.Cyan("x"),
		"bold":   p.Bold("x"),
		"dim":    p.Dim("x"),
	} {
		assert.True(t, strings.HasPrefix(styled, "\x1b["), "%s should start with ANSI", name)
		assert.True(t, strings.HasSuffix(styled, reset), "%s should reset", name)
		assert.Contains(t, styled, "x")
	}
	// Empty string is never wrapped.
	assert.Equal(t, "", p.Green(""))
}
