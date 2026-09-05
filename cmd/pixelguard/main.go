package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"

	"github.com/PixelCity-dev/pixelguard/internal/engine"
	"github.com/PixelCity-dev/pixelguard/internal/report"
	"github.com/PixelCity-dev/pixelguard/internal/types"
	"github.com/PixelCity-dev/pixelguard/internal/ui"
)

var (
	version   = "0.1.0"
	verbose   bool
	skipSec   bool
	skipBuild bool
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "pixelguard [project-path]",
		Short: "Interactive project auditor — build errors, dependency conflicts, and security findings",
		Long: `PixelGuard scans a Java project (Maven or Gradle) and produces
a single trustworthy report covering build errors, dependency
conflicts, and security vulnerabilities.

Created by Jay & PixelCity.`,
		Args: cobra.MaximumNArgs(1),
		RunE: runScan,
	}

	rootCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose output")
	rootCmd.Flags().BoolVar(&skipSec, "skip-security", false, "Skip security analysis")
	rootCmd.Flags().BoolVar(&skipBuild, "skip-build", false, "Skip build log analysis")

	rootCmd.AddCommand(&cobra.Command{
		Use:   "export [project-path]",
		Short: "Export report to file (json, markdown, sarif)",
		Long:  "Run a non-interactive scan and export results to a file.",
		Args:  cobra.MaximumNArgs(1),
		RunE:  runExport,
	})

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func runScan(cmd *cobra.Command, args []string) error {
	projectPath := "."
	if len(args) > 0 {
		projectPath = args[0]
	}

	ui.PrintBanner()

	var scanReport *types.Report
	scanDone := make(chan error, 1)

	go func() {
		cfg := types.ScanConfig{
			ProjectPath:  projectPath,
			SkipSecurity: skipSec,
			SkipBuild:    skipBuild,
			Verbose:      verbose,
		}

		eng := engine.New(cfg, func(event types.ProgressEvent) {
			ui.PrintProgress(event)
		})

		var err error
		scanReport, err = eng.Scan()
		scanDone <- err
	}()

	if err := <-scanDone; err != nil {
		fmt.Fprintf(os.Stderr, "\n%s Scan failed: %v\n", ui.ColorError("✘"), err)
		return err
	}

	ui.PrintScanHeader(scanReport.ProjectInfo.Type, projectPath)
	ui.PrintProjectInfo(scanReport.ProjectInfo)
	ui.PrintScanSummary(scanReport)

	if len(scanReport.BuildErrors) > 0 {
		ui.PrintBuildErrors(scanReport.BuildErrors)
	}

	if len(scanReport.Conflicts) > 0 {
		ui.PrintConflicts(scanReport.Conflicts)
	}

	if len(scanReport.Security) > 0 {
		ui.PrintSecurityFindings(scanReport.Security)
	}

	fmt.Println()
	fmt.Println("  Launching interactive TUI...")
	fmt.Println("  Press 'q' to quit, 'h' for help")
	fmt.Println()

	m := ui.NewModel()
	m.SetReport(scanReport)
	m.SetFindings(ui.BuildFindingsFromReport(scanReport))

	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "TUI error: %v\n", err)
		return err
	}

	return nil
}

func runExport(cmd *cobra.Command, args []string) error {
	projectPath := "."
	if len(args) > 0 {
		projectPath = args[0]
	}

	ui.PrintBanner()

	cfg := types.ScanConfig{
		ProjectPath:  projectPath,
		SkipSecurity: skipSec,
		SkipBuild:    skipBuild,
		Verbose:      verbose,
	}

	eng := engine.New(cfg, func(event types.ProgressEvent) {
		ui.PrintProgress(event)
	})

	scanReport, err := eng.Scan()
	if err != nil {
		fmt.Fprintf(os.Stderr, "\n%s Scan failed: %v\n", ui.ColorError("✘"), err)
		return err
	}

	ui.PrintScanHeader(scanReport.ProjectInfo.Type, projectPath)
	ui.PrintProjectInfo(scanReport.ProjectInfo)
	ui.PrintScanSummary(scanReport)

	if err := report.WriteJSON(scanReport, "pixelguard-report.json"); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to write JSON: %v\n", err)
	} else {
		fmt.Println(ui.ColorSuccess("  ✔ Exported: pixelguard-report.json"))
	}

	if err := report.WriteMarkdown(scanReport, "pixelguard-report.md"); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to write Markdown: %v\n", err)
	} else {
		fmt.Println(ui.ColorSuccess("  ✔ Exported: pixelguard-report.md"))
	}

	if err := report.WriteSARIF(scanReport, "pixelguard-report.sarif"); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to write SARIF: %v\n", err)
	} else {
		fmt.Println(ui.ColorSuccess("  ✔ Exported: pixelguard-report.sarif"))
	}

	fmt.Println()
	return nil
}
