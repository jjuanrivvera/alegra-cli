// Package ui provides small terminal-styling helpers that respect NO_COLOR and
// automatically disable themselves when output is not an interactive terminal.
package ui

import (
	"io"
	"os"

	"golang.org/x/term"
)

// forceNoColor disables color globally regardless of the writer (wired to the
// --no-color flag).
var forceNoColor bool

// SetNoColor turns coloring off (or back on) process-wide.
func SetNoColor(v bool) { forceNoColor = v }

// Palette renders styled strings when color is enabled for its writer.
type Palette struct{ enabled bool }

// For returns a Palette that styles output only when w is a color-capable
// terminal and color hasn't been disabled.
func For(w io.Writer) Palette { return Palette{enabled: colorEnabled(w)} }

// Enabled reports whether this palette will emit color.
func (p Palette) Enabled() bool { return p.enabled }

func colorEnabled(w io.Writer) bool {
	// https://no-color.org/ — any non-empty NO_COLOR disables color.
	if forceNoColor || os.Getenv("NO_COLOR") != "" {
		return false
	}
	f, ok := w.(*os.File)
	if !ok {
		return false
	}
	return term.IsTerminal(int(f.Fd()))
}

const (
	reset  = "\x1b[0m"
	bold   = "\x1b[1m"
	dim    = "\x1b[2m"
	red    = "\x1b[31m"
	green  = "\x1b[32m"
	yellow = "\x1b[33m"
	cyan   = "\x1b[36m"
)

func (p Palette) wrap(code, s string) string {
	if !p.enabled || s == "" {
		return s
	}
	return code + s + reset
}

// Green styles success text.
func (p Palette) Green(s string) string { return p.wrap(green, s) }

// Red styles error text.
func (p Palette) Red(s string) string { return p.wrap(red, s) }

// Yellow styles warning text.
func (p Palette) Yellow(s string) string { return p.wrap(yellow, s) }

// Cyan styles informational accents.
func (p Palette) Cyan(s string) string { return p.wrap(cyan, s) }

// Bold styles emphasized text.
func (p Palette) Bold(s string) string { return p.wrap(bold, s) }

// Dim styles secondary text.
func (p Palette) Dim(s string) string { return p.wrap(dim, s) }
