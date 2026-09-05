package report

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/PixelCity-dev/pixelguard/internal/types"
)

func WriteJSON(report *types.Report, path string) error {
	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal report: %w", err)
	}
	return os.WriteFile(path, data, 0644)
}

func WriteMarkdown(report *types.Report, path string) error {
	var sb strings.Builder

	sb.WriteString("# PixelGuard Scan Report\n\n")
	sb.WriteString(fmt.Sprintf("**Project:** `%s`\n", report.ProjectPath))
	sb.WriteString(fmt.Sprintf("**Scan Time:** %s\n", report.ScanTime.Format(time.RFC3339)))
	sb.WriteString(fmt.Sprintf("**Duration:** %s\n", report.Duration.Round(time.Millisecond)))
	sb.WriteString("\n")

	sb.WriteString("## Executive Summary\n\n")
	sb.WriteString(fmt.Sprintf("Scanned **%s** project with %d source files and %d resolved dependencies.\n",
		report.ProjectInfo.Type, report.ProjectInfo.SourceFiles, report.ProjectInfo.ResolvedDeps))
	sb.WriteString(fmt.Sprintf("Found **%d** build errors, **%d** dependency conflicts, and **%d** security findings.\n",
		report.Summary.TotalBuildErrors, report.Summary.TotalConflicts, report.Summary.TotalSecurity))
	if report.Summary.CriticalCount > 0 || report.Summary.HighCount > 0 {
		sb.WriteString(fmt.Sprintf("\n⚠️  **%d Critical** and **%d High** severity findings require immediate attention.\n",
			report.Summary.CriticalCount, report.Summary.HighCount))
	}
	sb.WriteString("\n---\n\n")

	if len(report.BuildErrors) > 0 {
		sb.WriteString("## Build Errors\n\n")
		for i, err := range report.BuildErrors {
			sb.WriteString(fmt.Sprintf("%d. **%s** (line %d)\n", i+1, err.File, err.Line))
			sb.WriteString(fmt.Sprintf("   %s\n\n", err.Message))
		}
		sb.WriteString("---\n\n")
	}

	if len(report.Conflicts) > 0 {
		sb.WriteString("## Dependency Conflicts\n\n")
		sb.WriteString("| Package | Conflicting Versions | Recommended |\n")
		sb.WriteString("|---------|---------------------|-------------|\n")
		for _, c := range report.Conflicts {
			sb.WriteString(fmt.Sprintf("| %s | %s | %s |\n",
				c.Package, strings.Join(c.Versions, ", "), c.RecommendedVersion))
		}
		sb.WriteString("\n---\n\n")
	}

	if len(report.Security) > 0 {
		sb.WriteString("## Security Findings\n\n")

		bySeverity := map[types.Severity][]types.SecurityFinding{
			types.SeverityCritical: {},
			types.SeverityHigh:     {},
			types.SeverityMedium:   {},
			types.SeverityLow:      {},
		}
		for _, f := range report.Security {
			bySeverity[f.Severity] = append(bySeverity[f.Severity], f)
		}

		for _, sev := range []types.Severity{types.SeverityCritical, types.SeverityHigh, types.SeverityMedium, types.SeverityLow} {
			findings := bySeverity[sev]
			if len(findings) == 0 {
				continue
			}
			sb.WriteString(fmt.Sprintf("### %s Severity (%d)\n\n", sev, len(findings)))
			for _, f := range findings {
				sb.WriteString(fmt.Sprintf("- **%s** — `%s:%s` v%s\n", f.Title, f.Package, f.Package, f.Version))
				if f.CVE != "" {
					sb.WriteString(fmt.Sprintf("  - CVE: %s\n", f.CVE))
				}
				if f.URL != "" {
					sb.WriteString(fmt.Sprintf("  - Reference: %s\n", f.URL))
				}
				sb.WriteString(fmt.Sprintf("  - Impact: %s\n", f.Impact))
				sb.WriteString(fmt.Sprintf("  - Fix: %s\n\n", f.Remediation))
			}
		}
		sb.WriteString("---\n\n")
	}

	sb.WriteString("## Recommendations\n\n")
	sb.WriteString(generateRecommendations(report))
	sb.WriteString("\n")

	return os.WriteFile(path, []byte(sb.String()), 0644)
}

func WriteSARIF(report *types.Report, path string) error {
	sarif := map[string]interface{}{
		"$schema": "https://raw.githubusercontent.com/oasis-tcs/sarif-spec/master/Schemata/sarif-schema-2.1.0.json",
		"version": "2.1.0",
		"runs": []map[string]interface{}{
			{
				"tool": map[string]interface{}{
					"driver": map[string]interface{}{
						"name":            "PixelGuard",
						"semanticVersion": "0.1.0",
						"informationUri":  "https://github.com/PixelCity-dev/pixelguard",
						"rules":           generateSARIFRules(report),
					},
				},
				"results": generateSARIFResults(report),
			},
		},
	}

	data, err := json.MarshalIndent(sarif, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal SARIF: %w", err)
	}
	return os.WriteFile(path, data, 0644)
}

func generateRecommendations(report *types.Report) string {
	var sb strings.Builder
	recommendations := make([]string, 0)

	if report.Summary.CriticalCount > 0 {
		recommendations = append(recommendations, fmt.Sprintf(
			"🔴 **URGENT:** Address %d critical security findings immediately", report.Summary.CriticalCount))
	}
	if report.Summary.HighCount > 0 {
		recommendations = append(recommendations, fmt.Sprintf(
			"🟠 **HIGH:** Resolve %d high-severity security findings before next release", report.Summary.HighCount))
	}
	if report.Summary.TotalConflicts > 0 {
		recommendations = append(recommendations, fmt.Sprintf(
			"🟡 **MEDIUM:** Resolve %d dependency conflicts to prevent runtime errors", report.Summary.TotalConflicts))
	}
	if report.Summary.TotalBuildErrors > 0 {
		recommendations = append(recommendations, fmt.Sprintf(
			"🔧 **BUILD:** Fix %d compilation errors", report.Summary.TotalBuildErrors))
	}
	if report.Summary.MediumCount > 0 {
		recommendations = append(recommendations, fmt.Sprintf(
			"🟡 **PLAN:** Schedule remediation for %d medium-severity findings", report.Summary.MediumCount))
	}
	if len(recommendations) == 0 {
		recommendations = append(recommendations, "✅ No critical issues found. Project looks healthy!")
	}

	for i, r := range recommendations {
		sb.WriteString(fmt.Sprintf("%d. %s\n", i+1, r))
	}
	return sb.String()
}

func generateSARIFRules(report *types.Report) []map[string]interface{} {
	var rules []map[string]interface{}

	seen := make(map[string]bool)
	for _, f := range report.Security {
		if seen[f.ID] {
			continue
		}
		seen[f.ID] = true

		rules = append(rules, map[string]interface{}{
			"id":               f.ID,
			"shortDescription": map[string]string{"text": f.Title},
			"fullDescription":  map[string]string{"text": f.Impact},
			"defaultConfiguration": map[string]interface{}{
				"level": severityToSARIF(f.Severity),
			},
		})
	}
	return rules
}

func generateSARIFResults(report *types.Report) []map[string]interface{} {
	var results []map[string]interface{}

	for _, f := range report.Security {
		result := map[string]interface{}{
			"ruleId":  f.ID,
			"message": map[string]string{"text": f.Title + " — " + f.Impact},
			"level":   severityToSARIF(f.Severity),
			"locations": []map[string]interface{}{
				{
					"physicalLocation": map[string]interface{}{
						"artifactLocation": map[string]string{
							"uri": f.Package,
						},
					},
				},
			},
		}
		results = append(results, result)
	}

	for _, e := range report.BuildErrors {
		result := map[string]interface{}{
			"ruleId":  "build-error",
			"message": map[string]string{"text": e.Message},
			"level":   "error",
			"locations": []map[string]interface{}{
				{
					"physicalLocation": map[string]interface{}{
						"artifactLocation": map[string]string{
							"uri": e.File,
						},
						"region": map[string]int{
							"startLine": e.Line,
						},
					},
				},
			},
		}
		results = append(results, result)
	}

	return results
}

func severityToSARIF(sev types.Severity) string {
	switch sev {
	case types.SeverityCritical, types.SeverityHigh:
		return "error"
	case types.SeverityMedium:
		return "warning"
	case types.SeverityLow:
		return "note"
	default:
		return "none"
	}
}
