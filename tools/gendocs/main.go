// Command gendocs generates the CLI command reference as Markdown from the
// cobra command tree into docs/commands, for the MkDocs site.
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra/doc"

	"github.com/jjuanrivvera/alegra-cli/commands"
)

const outDir = "docs/commands"

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "gendocs:", err)
		os.Exit(1)
	}
}

func run() error {
	root := commands.RootCmd()
	root.DisableAutoGenTag = true

	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return err
	}

	// Front-matter prepender so MkDocs renders a clean title per page.
	filePrepender := func(filename string) string {
		base := filepath.Base(filename)
		name := strings.TrimSuffix(base, filepath.Ext(base))
		title := strings.ReplaceAll(name, "_", " ")
		return fmt.Sprintf("---\ntitle: %s\n---\n\n", title)
	}
	linkHandler := func(name string) string { return name }

	if err := doc.GenMarkdownTreeCustom(root, outDir, filePrepender, linkHandler); err != nil {
		return err
	}

	// Write an index page listing the top-level commands.
	var b strings.Builder
	b.WriteString("---\ntitle: Command Reference\n---\n\n# Command Reference\n\n")
	b.WriteString("Auto-generated from the CLI. See [alegra](alegra.md) for the root command.\n\n")
	for _, c := range root.Commands() {
		if c.Hidden || c.Name() == "help" || c.Name() == "completion" {
			continue
		}
		fmt.Fprintf(&b, "- [%s](alegra_%s.md) — %s\n", c.Name(), c.Name(), c.Short)
	}
	return os.WriteFile(filepath.Join(outDir, "index.md"), []byte(b.String()), 0o644)
}
