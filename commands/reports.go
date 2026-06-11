package commands

import (
	"net/url"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/jjuanrivvera/alegra-cli/internal/api"
)

func init() {
	reportsCmd := &cobra.Command{
		Use:   "reports",
		Short: "Read-only Alegra reports",
		Long: `reports fetches read-only Alegra reports.

Each subcommand GETs a report subpath under /reports. Use --from / --to to bound
the date range and --start / --limit to paginate the paginated reports.

Note: some reports (income-statement, account-statement) are only available on
certain Alegra plans and may return HTTP 403 otherwise.`,
	}

	reportsCmd.AddCommand(
		readOnlyHints(newReportSalesByClientCmd()),
		readOnlyHints(newReportSalesByClientTotalsCmd()),
		readOnlyHints(newReportSalesBySellerCmd()),
		readOnlyHints(newReportIncomeStatementCmd()),
		readOnlyHints(newReportAccountStatementCmd()),
	)

	rootCmd.AddCommand(reportsCmd)
}

// reportFlags holds the flag values shared across report subcommands.
type reportFlags struct {
	from  string
	to    string
	start int
	limit int
}

// addRangeFlags registers the date-range and pagination flags common to report
// subcommands and returns a pointer to the bound values.
func addRangeFlags(cmd *cobra.Command) *reportFlags {
	f := &reportFlags{}
	fs := cmd.Flags()
	fs.StringVar(&f.from, "from", "", "Range start date (YYYY-MM-DD)")
	fs.StringVar(&f.to, "to", "", "Range end date (YYYY-MM-DD)")
	fs.IntVar(&f.start, "start", 0, "Pagination offset")
	fs.IntVar(&f.limit, "limit", 10, "Number of rows per page")
	return f
}

// values builds the url.Values for a report request from the shared flags.
func (f *reportFlags) values(paginated bool) url.Values {
	q := url.Values{}
	if f.from != "" {
		q.Set("from", f.from)
	}
	if f.to != "" {
		q.Set("to", f.to)
	}
	if paginated {
		q.Set("start", strconv.Itoa(f.start))
		q.Set("limit", strconv.Itoa(f.limit))
	}
	return q
}

// fetchReport GETs a report subpath and renders the rows from a {data:[...]}
// envelope.
func fetchReport(cmd *cobra.Command, path string, q url.Values, cols []string) error {
	client, err := getAPIClient()
	if err != nil {
		return err
	}
	var out api.ReportEnvelope
	if err := client.GetInto(cmd.Context(), path, q, &out); err != nil {
		return err
	}
	if flagDryRun {
		return nil
	}
	return render(cmd, out.Data, cols)
}

// fetchReportRaw GETs a report subpath and renders the whole response object
// (for reports that are not a {data:[...]} list).
func fetchReportRaw(cmd *cobra.Command, path string, q url.Values) error {
	client, err := getAPIClient()
	if err != nil {
		return err
	}
	var out any
	if err := client.GetInto(cmd.Context(), path, q, &out); err != nil {
		return err
	}
	if flagDryRun {
		return nil
	}
	return render(cmd, out, nil)
}

func newReportSalesByClientCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "sales-by-client",
		Short:   "Sales grouped by client over a date range",
		Args:    cobra.NoArgs,
		Example: "  alegra reports sales-by-client --from 2026-01-01 --to 2026-03-31",
	}
	f := addRangeFlags(cmd)
	var clientName, orderField, orderDirection string
	cmd.Flags().StringVar(&clientName, "client-name", "", "Filter by client name")
	cmd.Flags().StringVar(&orderField, "order-field", "", "Order by: totalDocuments, subTotal, total")
	cmd.Flags().StringVar(&orderDirection, "order-direction", "", "Order direction: ASC or DESC")
	cmd.RunE = func(cmd *cobra.Command, _ []string) error {
		q := f.values(true)
		if clientName != "" {
			q.Set("clientName", clientName)
		}
		if orderField != "" {
			q.Set("order_field", orderField)
		}
		if orderDirection != "" {
			q.Set("order_direction", orderDirection)
		}
		return fetchReport(cmd, "reports/sales-by-client", q,
			[]string{"clientName", "totalDocuments", "subTotal", "total"})
	}
	return cmd
}

func newReportSalesByClientTotalsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "sales-by-client-totals",
		Short:   "Grand totals of the sales-by-client report",
		Args:    cobra.NoArgs,
		Example: "  alegra reports sales-by-client-totals --from 2026-01-01 --to 2026-03-31",
	}
	f := addRangeFlags(cmd)
	cmd.RunE = func(cmd *cobra.Command, _ []string) error {
		return fetchReportRaw(cmd, "reports/sales-by-client-totals", f.values(false))
	}
	return cmd
}

func newReportSalesBySellerCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "sales-by-seller",
		Short:   "Sales grouped by seller over a date range",
		Args:    cobra.NoArgs,
		Example: "  alegra reports sales-by-seller --from 2026-01-01 --to 2026-03-31",
	}
	f := addRangeFlags(cmd)
	var sellerName, orderField, orderDirection string
	cmd.Flags().StringVar(&sellerName, "seller-name", "", "Filter by seller name")
	cmd.Flags().StringVar(&orderField, "order-field", "", "Order by: totalDocuments, totalPayed, beforeTaxes, total")
	cmd.Flags().StringVar(&orderDirection, "order-direction", "", "Order direction: ASC or DESC")
	cmd.RunE = func(cmd *cobra.Command, _ []string) error {
		q := f.values(true)
		if sellerName != "" {
			q.Set("sellerName", sellerName)
		}
		if orderField != "" {
			q.Set("order_field", orderField)
		}
		if orderDirection != "" {
			q.Set("order_direction", orderDirection)
		}
		return fetchReport(cmd, "reports/sales-by-seller", q,
			[]string{"sellerName", "totalDocuments", "totalPayed", "beforeTaxes", "total"})
	}
	return cmd
}

func newReportIncomeStatementCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "income-statement",
		Short:   "Income statement / profit & loss (estado de resultados) — plan-dependent",
		Args:    cobra.NoArgs,
		Example: "  alegra reports income-statement --from 2026-01-01 --to 2026-12-31",
	}
	f := addRangeFlags(cmd)
	cmd.RunE = func(cmd *cobra.Command, _ []string) error {
		return fetchReportRaw(cmd, "reports/income-statement", f.values(false))
	}
	return cmd
}

func newReportAccountStatementCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "account-statement",
		Short:   "Account statement for a contact (estado de cuenta) — plan-dependent",
		Args:    cobra.NoArgs,
		Example: "  alegra reports account-statement --client-id 12 --from 2026-01-01 --to 2026-12-31",
	}
	f := addRangeFlags(cmd)
	var clientID string
	cmd.Flags().StringVar(&clientID, "client-id", "", "Contact (client) ID")
	cmd.RunE = func(cmd *cobra.Command, _ []string) error {
		q := f.values(false)
		if clientID != "" {
			q.Set("client_id", clientID)
		}
		return fetchReportRaw(cmd, "reports/account-statement", q)
	}
	return cmd
}
