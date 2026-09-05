package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/PixelCity-dev/pixelguard/internal/types"
)

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7D56F4")).
			Padding(0, 1).
			MarginBottom(1)

	sectionStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#BD93F9")).
			MarginTop(1)

	selectedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#44475A")).
			Padding(0, 1)

	normalStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#F8F8F2")).
			Padding(0, 1)

	detailStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#8BE9FD")).
			PaddingLeft(2)

	statusBarStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#F8F8F2")).
			Background(lipgloss.Color("#6272A4")).
			Padding(0, 1)

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#6272A4"))
)

type phase int

const (
	phaseScanning phase = iota
	phaseSummary
	phaseFindings
	phaseDetail
	phaseHelp
)

type finding struct {
	severity types.Severity
	title    string
	detail   string
	source   string
	fix      string
	package_ string
	version  string
	url      string
}

type Model struct {
	phase           phase
	report          *types.Report
	findings        []finding
	selected        int
	width           int
	height          int
	scanning        bool
	spinner         int
	spinnerChars    []string
	progressPercent float64
	progressMsg     string
	err             error
	inputMode       bool
	inputBuffer     string
	inputPrompt     string
}

func NewModel() Model {
	return Model{
		phase:        phaseScanning,
		spinnerChars: []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"},
		spinner:      0,
		scanning:     true,
	}
}

func (m *Model) SetReport(report *types.Report) {
	m.report = report
	m.phase = phaseSummary
}

func (m *Model) SetFindings(findings []finding) {
	m.findings = findings
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		return m.handleKey(msg)
	}
	return m, nil
}

func (m Model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch m.phase {
	case phaseScanning:
		return m, nil

	case phaseSummary, phaseFindings:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "up", "k":
			if m.selected > 0 {
				m.selected--
			}
		case "down", "j":
			if m.selected < len(m.findings)-1 {
				m.selected++
			}
		case "enter", " ":
			if len(m.findings) > 0 {
				m.phase = phaseDetail
			}
		case "1":
			m.phase = phaseSummary
			m.selected = 0
		case "2":
			m.phase = phaseFindings
			m.selected = 0
		case "h", "?":
			m.phase = phaseHelp
		case "r":
			return m, tea.Quit
		}

	case phaseDetail:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "esc", "backspace":
			m.phase = phaseFindings
		case "up", "k":
			if m.selected > 0 {
				m.selected--
				m.phase = phaseFindings
			}
		case "down", "j":
			if m.selected < len(m.findings)-1 {
				m.selected++
				m.phase = phaseFindings
			}
		}

	case phaseHelp:
		switch msg.String() {
		case "q", "ctrl+c", "esc", "h", "?":
			m.phase = phaseFindings
		}
	}

	return m, nil
}

func (m Model) View() string {
	var sb strings.Builder

	sb.WriteString(titleStyle.Render(" PixelGuard Interactive Report "))
	sb.WriteString("\n\n")

	switch m.phase {
	case phaseScanning:
		sb.WriteString(m.viewScanning())
	case phaseSummary:
		sb.WriteString(m.viewSummary())
	case phaseFindings:
		sb.WriteString(m.viewFindings())
	case phaseDetail:
		sb.WriteString(m.viewDetail())
	case phaseHelp:
		sb.WriteString(m.viewHelp())
	}

	sb.WriteString("\n")
	sb.WriteString(m.viewStatusBar())

	return sb.String()
}

func (m Model) viewScanning() string {
	var sb strings.Builder

	spinner := m.spinnerChars[m.spinner%len(m.spinnerChars)]
	sb.WriteString(fmt.Sprintf("  %s %s\n", infoStyle.Render(spinner), m.progressMsg))

	bar := renderProgressBar(m.progressPercent, 40)
	sb.WriteString(fmt.Sprintf("  %s\n", bar))

	return sb.String()
}

func (m Model) viewSummary() string {
	var sb strings.Builder

	if m.report == nil {
		sb.WriteString("  No report available.\n")
		return sb.String()
	}

	sb.WriteString(sectionStyle.Render(" Executive Summary "))
	sb.WriteString("\n\n")

	sb.WriteString(fmt.Sprintf("  Project: %s\n", m.report.ProjectPath))
	sb.WriteString(fmt.Sprintf("  Type: %s | Files: %d | Dependencies: %d\n",
		m.report.ProjectInfo.Type,
		m.report.ProjectInfo.SourceFiles,
		m.report.ProjectInfo.ResolvedDeps))
	sb.WriteString(fmt.Sprintf("  Scan duration: %s\n\n", m.report.Duration.Round(1)))

	sb.WriteString(sectionStyle.Render(" Summary "))
	sb.WriteString("\n\n")

	rows := []struct {
		label string
		count int
		style lipgloss.Style
	}{
		{"Build Errors", m.report.Summary.TotalBuildErrors, errorStyle},
		{"Warnings", m.report.Summary.TotalWarnings, warningStyle},
		{"Dependency Conflicts", m.report.Summary.TotalConflicts, warningStyle},
		{"Critical Security", m.report.Summary.CriticalCount, criticalStyle},
		{"High Security", m.report.Summary.HighCount, highStyle},
		{"Medium Security", m.report.Summary.MediumCount, mediumStyle},
		{"Low Security", m.report.Summary.LowCount, lowStyle},
	}

	for _, row := range rows {
		sb.WriteString(fmt.Sprintf("  %s %s\n",
			row.style.Render(fmt.Sprintf("%-24s", row.label)),
			row.style.Render(fmt.Sprintf("%d", row.count))))
	}

	sb.WriteString("\n")
	sb.WriteString(helpStyle.Render("  [1] Summary  [2] Findings  [h] Help  [q] Quit"))
	sb.WriteString("\n")

	return sb.String()
}

func (m Model) viewFindings() string {
	var sb strings.Builder

	sb.WriteString(sectionStyle.Render(" Findings "))
	sb.WriteString("\n\n")

	if len(m.findings) == 0 {
		sb.WriteString("  No findings.\n")
		return sb.String()
	}

	start := 0
	end := len(m.findings)
	visibleRows := m.height - 8
	if visibleRows < 1 {
		visibleRows = 15
	}

	if end > visibleRows {
		if m.selected >= start+visibleRows {
			start = m.selected - visibleRows + 1
		}
		if start+visibleRows < end {
			end = start + visibleRows
		}
	}

	for i := start; i < end; i++ {
		f := m.findings[i]
		cursor := "  "
		if i == m.selected {
			cursor = selectedStyle.Render(" ▸ ")
		}

		var sevStyle lipgloss.Style
		switch f.severity {
		case types.SeverityCritical:
			sevStyle = criticalStyle
		case types.SeverityHigh:
			sevStyle = highStyle
		case types.SeverityMedium:
			sevStyle = mediumStyle
		default:
			sevStyle = lowStyle
		}

		line := fmt.Sprintf("%s %s [%s] %s",
			cursor,
			sevStyle.Render(fmt.Sprintf("%-8s", f.severity)),
			sevStyle.Render(f.severity.Symbol()),
			f.title)
		sb.WriteString(line + "\n")
	}

	sb.WriteString("\n")
	sb.WriteString(helpStyle.Render("  [↑↓] Navigate  [Enter] Detail  [1] Summary  [2] Findings  [h] Help  [q] Quit"))
	sb.WriteString("\n")

	return sb.String()
}

func (m Model) viewDetail() string {
	var sb strings.Builder

	if len(m.findings) == 0 || m.selected >= len(m.findings) {
		sb.WriteString("  No finding selected.\n")
		return sb.String()
	}

	f := m.findings[m.selected]

	sb.WriteString(sectionStyle.Render(" Finding Detail "))
	sb.WriteString("\n\n")

	sb.WriteString(fmt.Sprintf("  %s %s\n",
		criticalStyle.Render("Title:"),
		f.title))
	sb.WriteString(fmt.Sprintf("  %s %s:%s v%s\n",
		infoStyle.Render("Package:"),
		f.package_, f.package_, f.version))
	sb.WriteString(fmt.Sprintf("  %s %s\n",
		warningStyle.Render("Severity:"),
		f.severity))
	if f.source != "" {
		sb.WriteString(fmt.Sprintf("  %s %s\n", infoStyle.Render("Source:"), f.source))
	}
	if f.url != "" {
		sb.WriteString(fmt.Sprintf("  %s %s\n", infoStyle.Render("Reference:"), f.url))
	}
	sb.WriteString("\n")

	sb.WriteString(detailStyle.Render("Impact:"))
	sb.WriteString("\n")
	sb.WriteString(fmt.Sprintf("    %s\n\n", f.detail))

	sb.WriteString(detailStyle.Render("Remediation:"))
	sb.WriteString("\n")
	sb.WriteString(fmt.Sprintf("    %s\n", successStyle.Render(f.fix)))

	sb.WriteString("\n")
	sb.WriteString(helpStyle.Render("  [Esc] Back  [↑↓] Prev/Next  [q] Quit"))
	sb.WriteString("\n")

	return sb.String()
}

func (m Model) viewHelp() string {
	var sb strings.Builder

	sb.WriteString(sectionStyle.Render(" Keyboard Shortcuts "))
	sb.WriteString("\n\n")

	shortcuts := []struct {
		key         string
		description string
	}{
		{"↑/k", "Move selection up"},
		{"↓/j", "Move selection down"},
		{"Enter/Space", "View finding detail"},
		{"1", "Switch to Summary view"},
		{"2", "Switch to Findings view"},
		{"h/?", "Toggle this help screen"},
		{"Esc", "Go back / Exit detail"},
		{"q/Ctrl+C", "Quit"},
	}

	for _, s := range shortcuts {
		sb.WriteString(fmt.Sprintf("  %s %s\n",
			selectedStyle.Render(fmt.Sprintf(" %-12s", s.key)),
			s.description))
	}

	sb.WriteString("\n")
	sb.WriteString(helpStyle.Render("  Press any key to return"))
	sb.WriteString("\n")

	return sb.String()
}

func (m Model) viewStatusBar() string {
	var parts []string

	if m.report != nil {
		parts = append(parts, fmt.Sprintf("Project: %s", m.report.ProjectInfo.Type))
		parts = append(parts, fmt.Sprintf("Findings: %d", len(m.findings)))
	}
	parts = append(parts, fmt.Sprintf("Phase: %s", phaseName(m.phase)))

	return statusBarStyle.Render(strings.Join(parts, " │ "))
}

func phaseName(p phase) string {
	switch p {
	case phaseScanning:
		return "Scanning"
	case phaseSummary:
		return "Summary"
	case phaseFindings:
		return "Findings"
	case phaseDetail:
		return "Detail"
	case phaseHelp:
		return "Help"
	default:
		return "Unknown"
	}
}

func BuildFindingsFromReport(report *types.Report) []finding {
	var findings []finding

	for _, e := range report.BuildErrors {
		findings = append(findings, finding{
			severity: e.Severity,
			title:    fmt.Sprintf("Build Error: %s", e.Message),
			detail:   e.Raw,
			source:   "Build System",
			fix:      "Fix the compilation error in " + e.File,
		})
	}

	for _, c := range report.Conflicts {
		findings = append(findings, finding{
			severity: c.Severity,
			title:    fmt.Sprintf("Dependency Conflict: %s", c.Package),
			detail:   fmt.Sprintf("Multiple versions found: %s\nRecommended: %s", strings.Join(c.Versions, ", "), c.RecommendedVersion),
			source:   "Dependency Resolver",
			fix:      fmt.Sprintf("Upgrade to %s and exclude conflicting transitive dependencies", c.RecommendedVersion),
		})
	}

	for _, s := range report.Security {
		f := finding{
			severity: s.Severity,
			title:    s.Title,
			detail:   s.Impact,
			source:   s.Source,
			fix:      s.Remediation,
			package_: s.Package,
			version:  s.Version,
			url:      s.URL,
		}
		if s.CVE != "" {
			f.title = fmt.Sprintf("[%s] %s", s.CVE, s.Title)
		}
		findings = append(findings, f)
	}

	return findings
}
