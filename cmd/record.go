package cmd

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/goccy/go-json"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/lynicis/inzibat/config"
	_ "github.com/lynicis/inzibat/log"
	"github.com/lynicis/inzibat/recorder"
)

const (
	defaultRecordAddr   = "localhost:8080"
	defaultExportOutput = "recorded-session.json"
	defaultExportFormat = "json"
	exportFilePerm      = 0644
)

var (
	recordAddr         string
	recordExportOutput string
	recordExportFormat string
)

var recordCmd = &cobra.Command{
	Use:   "record",
	Short: "Manage request recordings from a running inzibat server",
	Long: `Manage request recordings from a running inzibat server.

The server must be running with the --record flag enabled.
Use subcommands to list, export, or clear recorded requests.`,
}

var recordListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls", "l"},
	Short:   "List recorded requests",
	Long: `List all recorded HTTP requests from the running inzibat server.

Displays a table with method, path, status code, duration, and timestamp
for each recorded request.`,
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		entries, err := fetchRecordedEntries(recordAddr)
		if err != nil {
			zap.L().Fatal("failed to fetch recorded entries", zap.Error(err))
		}

		if len(entries) == 0 {
			fmt.Println("No recorded requests.")
			return
		}

		rows := make([][]string, 0, len(entries))
		for _, entry := range entries {
			rows = append(rows, []string{
				entry.Request.Method,
				entry.Request.Path,
				fmt.Sprintf("%d", entry.Response.StatusCode),
				fmt.Sprintf("%dms", entry.DurationMs),
				entry.Timestamp.Format(time.RFC3339),
			})
		}

		t := table.New().
			Border(lipgloss.NormalBorder()).
			Headers("METHOD", "PATH", "STATUS", "DURATION", "TIMESTAMP").
			Rows(rows...)

		fmt.Println(t)
	},
}

var recordExportCmd = &cobra.Command{
	Use:     "export",
	Aliases: []string{"e"},
	Short:   "Export recorded session to a file",
	Long: `Export recorded HTTP requests from the running inzibat server.

Supports two formats:
  - json:    Raw recorded session data (default)
  - inzibat: Converted to inzibat mock configuration format`,
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		session, err := fetchRecordedSession(recordAddr)
		if err != nil {
			zap.L().Fatal("failed to fetch recorded session", zap.Error(err))
		}

		if len(session.Entries) == 0 {
			fmt.Println("No recorded requests to export.")
			return
		}

		switch recordExportFormat {
		case "inzibat":
			exportAsInzibatConfig(session)
		default:
			exportAsJSON(session)
		}
	},
}

var recordClearCmd = &cobra.Command{
	Use:     "clear",
	Aliases: []string{"clr"},
	Short:   "Clear all recorded requests",
	Long:    `Clear all recorded HTTP requests from the running inzibat server.`,
	Args:    cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		url := fmt.Sprintf("http://%s/_inzibat/recorder/clear", recordAddr)

		resp, err := http.Post(url, "application/json", nil) //nolint:noctx
		if err != nil {
			zap.L().Fatal(
				"failed to clear recorded entries",
				zap.String("addr", recordAddr),
				zap.Error(err),
			)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			zap.L().Fatal(
				"unexpected response from server",
				zap.Int("statusCode", resp.StatusCode),
			)
		}

		fmt.Println("All recorded requests cleared.")
	},
}

func fetchRecordedEntries(addr string) ([]recorder.RecordedEntry, error) {
	url := fmt.Sprintf("http://%s/_inzibat/recorder/entries", addr)

	resp, err := http.Get(url) //nolint:noctx
	if err != nil {
		return nil, fmt.Errorf(
			"could not connect to inzibat server at %s: %w",
			addr,
			err,
		)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var entries []recorder.RecordedEntry
	if err = json.Unmarshal(body, &entries); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return entries, nil
}

func fetchRecordedSession(addr string) (*recorder.RecordedSession, error) {
	url := fmt.Sprintf("http://%s/_inzibat/recorder/session", addr)

	resp, err := http.Get(url) //nolint:noctx
	if err != nil {
		return nil, fmt.Errorf(
			"could not connect to inzibat server at %s: %w",
			addr,
			err,
		)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var session recorder.RecordedSession
	if err = json.Unmarshal(body, &session); err != nil {
		return nil, fmt.Errorf("failed to parse session: %w", err)
	}

	return &session, nil
}

func exportAsJSON(session *recorder.RecordedSession) {
	data, err := json.MarshalIndent(session, "", "  ")
	if err != nil {
		zap.L().Fatal("failed to marshal session", zap.Error(err))
	}

	absPath, err := config.ResolveAbsolutePath(recordExportOutput)
	if err != nil {
		zap.L().Fatal("failed to resolve output path", zap.Error(err))
	}

	// #nosec G306 - Export files are user-visible data, not secrets
	if err = os.WriteFile(absPath, data, exportFilePerm); err != nil {
		zap.L().Fatal("failed to write export file", zap.Error(err))
	}

	zap.L().Info(
		"Session exported as JSON",
		zap.String("file", recordExportOutput),
		zap.Int("entries", len(session.Entries)),
	)
}

func exportAsInzibatConfig(session *recorder.RecordedSession) {
	cfg := recorder.ConvertToInzibatConfig(*session, 0)

	if err := config.WriteConfig(cfg, recordExportOutput); err != nil {
		zap.L().Fatal("failed to write inzibat config", zap.Error(err))
	}

	zap.L().Info(
		"Session exported as inzibat config",
		zap.String("file", recordExportOutput),
		zap.Int("routes", len(cfg.Routes)),
	)
}

func init() {
	recordCmd.PersistentFlags().StringVarP(
		&recordAddr,
		"addr",
		"a",
		defaultRecordAddr,
		"Address of the running inzibat server",
	)

	recordExportCmd.Flags().StringVarP(
		&recordExportOutput,
		"output",
		"o",
		defaultExportOutput,
		"Output file path",
	)
	recordExportCmd.Flags().StringVarP(
		&recordExportFormat,
		"format",
		"f",
		defaultExportFormat,
		"Export format: json or inzibat",
	)

	recordCmd.AddCommand(recordListCmd)
	recordCmd.AddCommand(recordExportCmd)
	recordCmd.AddCommand(recordClearCmd)
	rootCmd.AddCommand(recordCmd)
}
