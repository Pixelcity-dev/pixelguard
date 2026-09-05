package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/fatih/color"
	"github.com/PixelCity-dev/pixelguard/internal/types"
)

var (
	headerStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7D56F4")).
			Padding(0, 1)

	successStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#04B575")).
			Bold(true)

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF5F5F")).
			Bold(true)

	warningStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFD93D")).
			Bold(true)

	infoStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#6C7ED4"))

	criticalStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF0000")).
			Bold(true)

	highStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF8C00")).
			Bold(true)

	mediumStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFD700"))

	lowStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#32CD32"))

	dimStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#6272A4"))
)

func PrintBanner() {
	banner := `
   ____  _           _     _     _
  |  _ \(_)_ __   __| |___| |__ (_)_   _
  | |_) | | '_ \ / _  / __| '_ \| \ \ / /
  |  __/| | | | | (_| \__ \ | | | |\ V /
  |_|   |_|_| |_|\__,_|___/_| |_|_| \_/

  Interactive Project Auditor v0.1.0
  By Jay & PixelCity
`
	fmt.Println(headerStyle.Render(banner))
}

func PrintScanHeader(projectType types.ProjectType, path string) {
	fmt.Printf("\n%s Scanning project: %s\n",
		infoStyle.Render("🔍"),
		path)
	fmt.Printf("%s Detected: %s project\n",
		successStyle.Render("✔"),
		projectType.String())
}

func PrintProjectInfo(info types.ProjectInfo) {
	if info.JavaVersion != "" {
		fmt.Printf("%s Java version: %s\n", successStyle.Render("✔"), info.JavaVersion)
	}
	fmt.Printf("%s %d source files, %d direct dependencies, %d resolved (incl. transitive)\n",
		successStyle.Render("✔"),
		info.SourceFiles, info.DirectDeps, info.ResolvedDeps)
}

func PrintScanSummary(report *types.Report) {
	fmt.Println()

	buildStr := fmt.Sprintf("Build:        %d errors, %d warnings",
		report.Summary.TotalBuildErrors, report.Summary.TotalWarnings)
	conflictStr := fmt.Sprintf("Conflicts:    %d dependency version conflicts",
		report.Summary.TotalConflicts)
	securityStr := fmt.Sprintf("Security:     %d findings — %d Critical, %d High, %d Medium, %d Low",
		report.Summary.TotalSecurity,
		report.Summary.CriticalCount,
		report.Summary.HighCount,
		report.Summary.MediumCount,
		report.Summary.LowCount)

	if report.Summary.TotalBuildErrors > 0 {
		buildStr = errorStyle.Render(buildStr)
	} else {
		buildStr = successStyle.Render(buildStr)
	}

	if report.Summary.TotalConflicts > 0 {
		conflictStr = warningStyle.Render(conflictStr)
	} else {
		conflictStr = successStyle.Render(conflictStr)
	}

	if report.Summary.TotalSecurity > 0 {
		securityStr = warningStyle.Render(securityStr)
	} else {
		securityStr = successStyle.Render(securityStr)
	}

	fmt.Println(buildStr)
	fmt.Println(conflictStr)
	fmt.Println(securityStr)
	fmt.Println()
}

func PrintBuildErrors(errors []types.BuildError) {
	if len(errors) == 0 {
		return
	}
	fmt.Println(headerStyle.Render(" Build Errors "))
	for i, e := range errors {
		loc := e.File
		if e.Line > 0 {
			loc = fmt.Sprintf("%s:%d", e.File, e.Line)
		}
		fmt.Printf("  %s %d. %s\n", errorStyle.Render("ERROR"), i+1, loc)
		fmt.Printf("     %s\n\n", e.Message)
	}
}

func PrintConflicts(conflicts []types.DependencyConflict) {
	if len(conflicts) == 0 {
		return
	}
	fmt.Println(headerStyle.Render(" Dependency Conflicts "))
	for i, c := range conflicts {
		fmt.Printf("  %s %d. %s\n", warningStyle.Render("⚡"), i+1, c.Package)
		fmt.Printf("     Versions: %s\n", strings.Join(c.Versions, ", "))
		fmt.Printf("     Recommended: %s\n\n", successStyle.Render(c.RecommendedVersion))
	}
}

func PrintSecurityFindings(findings []types.SecurityFinding) {
	if len(findings) == 0 {
		return
	}
	fmt.Println(headerStyle.Render(" Security Findings "))

	bySeverity := map[types.Severity][]types.SecurityFinding{
		types.SeverityCritical: {},
		types.SeverityHigh:     {},
		types.SeverityMedium:   {},
		types.SeverityLow:      {},
	}
	for _, f := range findings {
		bySeverity[f.Severity] = append(bySeverity[f.Severity], f)
	}

	idx := 1
	for _, sev := range []types.Severity{types.SeverityCritical, types.SeverityHigh, types.SeverityMedium, types.SeverityLow} {
		findings := bySeverity[sev]
		if len(findings) == 0 {
			continue
		}

		var style lipgloss.Style
		switch sev {
		case types.SeverityCritical:
			style = criticalStyle
		case types.SeverityHigh:
			style = highStyle
		case types.SeverityMedium:
			style = mediumStyle
		case types.SeverityLow:
			style = lowStyle
		}

		for _, f := range findings {
			fmt.Printf("  %s [%s] %d. %s\n",
				style.Render(sev.Symbol()),
				style.Render(sev.String()),
				idx, f.Title)
			fmt.Printf("     Package: %s:%s v%s\n", dimStyle.Render(f.Package), dimStyle.Render(f.Package), f.Version)
			if f.CVE != "" {
				fmt.Printf("     CVE: %s\n", f.CVE)
			}
			if f.URL != "" {
				fmt.Printf("     Reference: %s\n", dimStyle.Render(f.URL))
			}
			fmt.Printf("     Impact: %s\n", f.Impact)
			fmt.Printf("     Fix: %s\n\n", successStyle.Render(f.Remediation))
			idx++
		}
	}
}

func PrintProgress(event types.ProgressEvent) {
	if event.Error != "" {
		fmt.Printf("\r  %s %s\n", errorStyle.Render("✘"), event.Error)
		return
	}

	if event.Done {
		fmt.Printf("\r  %s %s\n", successStyle.Render("✔"), event.Message)
	} else {
		bar := renderProgressBar(event.Percent, 30)
		fmt.Printf("\r  %s %s %s", infoStyle.Render("⟳"), event.Message, bar)
	}
}

func renderProgressBar(percent float64, width int) string {
	filled := int(percent / 100 * float64(width))
	if filled > width {
		filled = width
	}
	bar := strings.Repeat("█", filled) + strings.Repeat("░", width-filled)
	return fmt.Sprintf("[%s] %.0f%%", dimStyle.Render(bar), percent)
}

func ColorError(msg string) string {
	return color.New(color.FgRed, color.Bold).Sprint(msg)
}

func ColorSuccess(msg string) string {
	return color.New(color.FgGreen, color.Bold).Sprint(msg)
}

func ColorWarning(msg string) string {
	return color.New(color.FgYellow, color.Bold).Sprint(msg)
}

func ColorInfo(msg string) string {
	return color.New(color.FgCyan).Sprint(msg)
}
