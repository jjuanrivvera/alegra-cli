package commands

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/jjuanrivvera/alegra-cli/internal/ui"
	"github.com/jjuanrivvera/alegra-cli/internal/version"
)

const latestReleaseURL = "https://api.github.com/repos/jjuanrivvera/alegra-cli/releases/latest"

func init() {
	var asJSON, check bool
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Print version information",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			out := cmd.OutOrStdout()
			if asJSON {
				enc := json.NewEncoder(out)
				enc.SetIndent("", "  ")
				return enc.Encode(map[string]string{
					"version":   version.Short(),
					"commit":    version.Commit,
					"buildDate": version.BuildDate,
					"goVersion": runtime.Version(),
					"os":        runtime.GOOS,
					"arch":      runtime.GOARCH,
				})
			}
			fmt.Fprintln(out, version.Info())
			if check {
				printUpdateCheck(cmd.Context(), out)
			}
			return nil
		},
	}
	cmd.Flags().BoolVar(&asJSON, "json", false, "Output version information as JSON")
	cmd.Flags().BoolVar(&check, "check", false, "Check whether a newer release is available")
	rootCmd.AddCommand(cmd)
	rootCmd.Version = version.Short()
}

// printUpdateCheck queries the latest GitHub release and reports whether an
// update is available. Failures are non-fatal (network/offline tolerant).
func printUpdateCheck(ctx context.Context, out io.Writer) {
	c := ui.For(out)
	latest, err := fetchLatestRelease(ctx, &http.Client{Timeout: 5 * time.Second}, latestReleaseURL)
	if err != nil {
		fmt.Fprintf(out, "%s %v\n", c.Yellow("update check failed:"), err)
		return
	}
	if isNewer(latest, version.Short()) {
		fmt.Fprintf(out, "%s %s is available (you have %s) — run `brew upgrade alegra-cli` or `go install …@latest`.\n",
			c.Yellow("update available:"), latest, version.Short())
		return
	}
	fmt.Fprintln(out, c.Green("You're on the latest version."))
}

// fetchLatestRelease returns the tag name of the newest GitHub release.
func fetchLatestRelease(ctx context.Context, hc *http.Client, url string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	resp, err := hc.Do(req)
	if err != nil {
		return "", err
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("GitHub API returned %d", resp.StatusCode)
	}
	var body struct {
		TagName string `json:"tag_name"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return "", err
	}
	if body.TagName == "" {
		return "", fmt.Errorf("no release tag in response")
	}
	return body.TagName, nil
}

// isNewer reports whether latest is a higher semantic version than current.
// It tolerates a leading "v" and any pre-release/build/git-describe suffix.
func isNewer(latest, current string) bool {
	lv, cv := parseSemver(latest), parseSemver(current)
	if lv == nil || cv == nil {
		return false
	}
	for i := range 3 {
		if lv[i] != cv[i] {
			return lv[i] > cv[i]
		}
	}
	return false
}

func parseSemver(s string) []int {
	s = strings.TrimPrefix(strings.TrimSpace(s), "v")
	if i := strings.IndexAny(s, "-+"); i >= 0 {
		s = s[:i] // drop pre-release / build / git-describe suffix
	}
	parts := strings.Split(s, ".")
	if len(parts) < 3 {
		return nil
	}
	out := make([]int, 3)
	for i := range 3 {
		n, err := strconv.Atoi(parts[i])
		if err != nil {
			return nil
		}
		out[i] = n
	}
	return out
}
