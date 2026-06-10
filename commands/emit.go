package commands

import (
	"encoding/json"
	"errors"
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

			// 2. Idempotency guard. A cache that exists but cannot be read is a
			// hard stop: proceeding with an empty guard could re-emit invoices
			// that were already stamped (a fiscal duplicate, not undoable).
			profile := currentProfileName()
			cache, cerr := loadEmitCache()
			if cerr != nil {
				return fmt.Errorf("cannot read the emission idempotency cache: %w\nInspect or remove %s, then retry", cerr, emitCachePath())
			}
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
				// Persist after every batch, and stop on failure: stamping more
				// invoices the guard cannot record would risk re-emission on the
				// next run.
				if werr := saveEmitCache(cache); werr != nil {
					return fmt.Errorf("stamped %v but could not record them in the idempotency cache: %w\nRecord these ids in %s before re-running, or a re-run may emit them twice", batch, werr, emitCachePath())
				}
			}
			if !flagDryRun {
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

// loadEmitCache reads the emitted-ids cache. A missing file is a fresh start;
// an unreadable or corrupt file is an error — silently treating it as empty
// would drop the idempotency guard and allow double emission.
func loadEmitCache() (map[string]bool, error) {
	cache := map[string]bool{}
	data, err := os.ReadFile(emitCachePath()) //nolint:gosec // path under config dir
	if errors.Is(err, os.ErrNotExist) {
		return cache, nil
	}
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(data, &cache); err != nil {
		return nil, fmt.Errorf("corrupt cache %s: %w", emitCachePath(), err)
	}
	return cache, nil
}

// saveEmitCache persists atomically (write temp + rename): a crash mid-write
// must never leave a torn emitted.json, which would read as corrupt and block
// (or, worse, lose) the guard.
func saveEmitCache(cache map[string]bool) error {
	path := emitCachePath()
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return err
	}
	data, err := json.Marshal(cache)
	if err != nil {
		return err
	}
	tmp, err := os.CreateTemp(dir, "emitted-*.json")
	if err != nil {
		return err
	}
	defer func() { _ = os.Remove(tmp.Name()) }() // no-op once renamed
	if err := tmp.Chmod(0o600); err != nil {
		_ = tmp.Close()
		return err
	}
	if _, err := tmp.Write(data); err != nil {
		_ = tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	return os.Rename(tmp.Name(), path)
}

// currentProfileName returns the active profile name for cache scoping.
func currentProfileName() string {
	if cfg, err := config.Load(); err == nil {
		return cfg.ActiveProfileName(flagProfile)
	}
	return "default"
}
