package commands

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/jjuanrivvera/alegra-cli/internal/api"
	"github.com/jjuanrivvera/alegra-cli/internal/config"
)

// maxStampBatch is Alegra's cap on invoices per POST /invoices/stamp call.
const maxStampBatch = 10

// newInvoicesEmitCmd builds `alegra invoices emit` — the safe electronic
// emission workflow: it gathers draft/open invoices, skips any already emitted
// (local idempotency guard), and stamps them in auto-chunked batches of 10.
func newInvoicesEmitCmd(sp resourceSpec[api.Invoice]) *cobra.Command {
	var (
		all   bool
		force bool
	)
	cmd := &cobra.Command{
		Use:   "emit [id...]",
		Short: "Emit (stamp) draft/open invoices electronically, in batches of 10",
		Long: `emit sends invoices to the tax authority for stamping.

Pass invoice ids explicitly, or use --all to emit every draft invoice. Already-
emitted ids are skipped (a local idempotency guard) unless --force is given —
emission is NOT idempotent on Alegra's side, so this prevents duplicates.`,
		Example: "  alegra invoices emit 1234 1235\n" +
			"  alegra invoices emit --all            # every draft\n" +
			"  alegra invoices emit --all --dry-run",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			res := sp.New(client)

			// 1. Gather target ids.
			ids := args
			if all {
				params := api.ListParams{Extra: map[string][]string{"status": {"draft"}}}
				items, lerr := res.ListAll(cmd.Context(), params, 0)
				if lerr != nil {
					return lerr
				}
				ids = nil
				for i := range items {
					if id := invoiceID(items[i]); id != "" {
						ids = append(ids, id)
					}
				}
			}
			if len(ids) == 0 {
				return fmt.Errorf("no invoices to emit (pass ids or --all)")
			}

			// 2. Idempotency guard.
			profile := currentProfileName()
			cache, _ := loadEmitCache()
			todo, skipped := filterEmitted(ids, cache, profile, force)
			out := cmd.OutOrStdout()
			if len(skipped) > 0 {
				fmt.Fprintf(out, "Skipping %d already-emitted invoice(s): %v (use --force to re-emit)\n", len(skipped), skipped)
			}
			if len(todo) == 0 {
				fmt.Fprintln(out, "Nothing to emit.")
				return nil
			}

			// 3. Stamp in chunks of <=10.
			var emitted, failedChunks int
			for _, batch := range chunk(todo, maxStampBatch) {
				if flagDryRun {
					fmt.Fprintf(out, "would stamp: %v\n", batch)
					continue
				}
				body := map[string]any{"ids": batch}
				var resp map[string]any
				if serr := res.CollectionAction(cmd.Context(), "stamp", body, &resp); serr != nil {
					failedChunks++
					fmt.Fprintf(cmd.ErrOrStderr(), "batch %v FAILED: %v\n", batch, serr)
					continue
				}
				for _, id := range batch {
					cache[emitKey(profile, id)] = true
					emitted++
				}
				fmt.Fprintf(out, "stamped: %v\n", batch)
			}
			if !flagDryRun {
				_ = saveEmitCache(cache)
				fmt.Fprintf(out, "Emitted %d invoice(s); %d batch(es) failed.\n", emitted, failedChunks)
			}
			if failedChunks > 0 {
				return fmt.Errorf("%d batch(es) failed", failedChunks)
			}
			return nil
		},
	}
	cmd.Flags().BoolVar(&all, "all", false, "Emit every draft invoice")
	cmd.Flags().BoolVar(&force, "force", false, "Re-emit even if locally recorded as emitted")
	return cmd
}

func invoiceID(inv api.Invoice) string { return inv.ID.String() }

// --- idempotency guard helpers (pure, unit-tested) ---

func emitKey(profile, id string) string { return profile + ":" + id }

// chunk splits ids into groups of at most size.
func chunk(ids []string, size int) [][]string {
	if size < 1 {
		size = 1
	}
	var out [][]string
	for i := 0; i < len(ids); i += size {
		end := min(i+size, len(ids))
		out = append(out, ids[i:end])
	}
	return out
}

// filterEmitted splits ids into those to emit and those already recorded as
// emitted for the profile. With force, nothing is skipped.
func filterEmitted(ids []string, cache map[string]bool, profile string, force bool) (todo, skipped []string) {
	for _, id := range ids {
		if !force && cache[emitKey(profile, id)] {
			skipped = append(skipped, id)
			continue
		}
		todo = append(todo, id)
	}
	return todo, skipped
}

// --- cache persistence ---

func emitCachePath() string {
	return filepath.Join(filepath.Dir(config.DefaultPath()), "emitted.json")
}

func loadEmitCache() (map[string]bool, error) {
	cache := map[string]bool{}
	data, err := os.ReadFile(emitCachePath()) //nolint:gosec // path under config dir
	if err != nil {
		return cache, err
	}
	_ = json.Unmarshal(data, &cache)
	return cache, nil
}

func saveEmitCache(cache map[string]bool) error {
	path := emitCachePath()
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return err
	}
	data, err := json.Marshal(cache)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o600)
}

// currentProfileName returns the active profile name for cache scoping.
func currentProfileName() string {
	if cfg, err := config.Load(); err == nil {
		return cfg.ActiveProfileName(flagProfile)
	}
	return "default"
}
