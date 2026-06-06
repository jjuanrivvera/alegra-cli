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
		Short: "Read-only Alegra sales reports",
		Long: `reports fetches read-only Alegra sales reports.

Each subcommand GETs a report subpath under /reports and aggregates sales
documents over a date range. Use --from / --to to bound the range and
--start / --limit to paginate.`,
	}

	reportsCmd.AddCommand(
		newReportSalesDocumentsCmd(),
		newReportSalesTotalsCmd(),
		newReportSalesByClientCmd(),
		newReportSalesBySellerCmd(),
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

// addRangeFlags registers the date-range and pagination flags common to every
// report subcommand and returns a pointer to the bound values.
func addRangeFlags(cmd *cobra.Command) *reportFlags {
	f := &reportFlags{}
	fs := cmd.Flags()
	fs.StringVar(&f.from, "from", "", "Range start date (YYYY-MM-DD)")
	fs.StringVar(&f.to, "to", "", "Range end date (YYYY-MM-DD)")
	fs.IntVar(&f.start, "start", 0, "Pagination offset")
	fs.IntVar(&f.limit, "limit", 10, "Number of rows per page")
	return f
}

// values builds the url.Values for a report request from the shared flags,
// applying only the flags the caller actually set.
func (f *reportFlags) values(cmd *cobra.Command, paginated bool) url.Values {
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

// fetchReport GETs a report subpath and renders the resulting rows.
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

func newReportSalesDocumentsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sales-documents",
		Short: "List individual sales documents (invoices, credit/debit notes) in a range",
		Args:  cobra.NoArgs,
	}
	f := addRangeFlags(cmd)
	var docType, docStatus, docNumber string
	cmd.Flags().StringVar(&docType, "document-types", "", "Comma-separated types: invoice, creditNote, incomeDebitNote")
	cmd.Flags().StringVar(&docStatus, "document-status", "", "Document status: open, closed, applied")
	cmd.Flags().StringVar(&docNumber, "document-number", "", "Filter by document number")
	cmd.RunE = func(cmd *cobra.Command, _ []string) error {
		q := f.values(cmd, true)
		if docType != "" {
			q.Set("documentTypes", docType)
		}
		if docStatus != "" {
			q.Set("documentStatus", docStatus)
		}
		if docNumber != "" {
			q.Set("documentNumber", docNumber)
		}
		return fetchReport(cmd, "reports/sales-documents", q,
			[]string{"id", "documentNumber", "documentType", "date", "status", "total"})
	}
	return cmd
}

func newReportSalesTotalsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sales-totals",
		Short: "Sales totals grouped by period (day, month, or year)",
		Args:  cobra.NoArgs,
	}
	f := addRangeFlags(cmd)
	var groupBy, docStatus string
	cmd.Flags().StringVar(&groupBy, "group-by", "month", "Temporal grouping: day, month, year")
	cmd.Flags().StringVar(&docStatus, "document-status", "", "Document status: open, closed, applied")
	cmd.RunE = func(cmd *cobra.Command, _ []string) error {
		q := f.values(cmd, false)
		if groupBy != "" {
			q.Set("groupBy", groupBy)
		}
		if docStatus != "" {
			q.Set("documentStatus", docStatus)
		}
		return fetchReport(cmd, "reports/sales-totals", q,
			[]string{"period", "beforeTaxes", "tax", "discount", "creditNote", "total"})
	}
	return cmd
}

func newReportSalesByClientCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sales-by-client",
		Short: "Sales grouped by client over a date range",
		Args:  cobra.NoArgs,
	}
	f := addRangeFlags(cmd)
	var clientName, orderField, orderDirection string
	cmd.Flags().StringVar(&clientName, "client-name", "", "Filter by client name")
	cmd.Flags().StringVar(&orderField, "order-field", "", "Order by: totalDocuments, subTotal, total")
	cmd.Flags().StringVar(&orderDirection, "order-direction", "", "Order direction: ASC or DESC")
	cmd.RunE = func(cmd *cobra.Command, _ []string) error {
		q := f.values(cmd, true)
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

func newReportSalesBySellerCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sales-by-seller",
		Short: "Sales grouped by seller over a date range",
		Args:  cobra.NoArgs,
	}
	f := addRangeFlags(cmd)
	var sellerName, orderField, orderDirection string
	cmd.Flags().StringVar(&sellerName, "seller-name", "", "Filter by seller name")
	cmd.Flags().StringVar(&orderField, "order-field", "", "Order by: totalDocuments, totalPayed, beforeTaxes, total")
	cmd.Flags().StringVar(&orderDirection, "order-direction", "", "Order direction: ASC or DESC")
	cmd.RunE = func(cmd *cobra.Command, _ []string) error {
		q := f.values(cmd, true)
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
