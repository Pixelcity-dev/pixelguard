package analyzer

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/PixelCity-dev/pixelguard/internal/types"
)

const osvAPIURL = "https://api.osv.dev/v1/query"

type osvQuery struct {
	Package struct {
		Name      string `json:"name"`
		Ecosystem string `json:"ecosystem"`
	} `json:"package"`
	Version string `json:"version"`
}

type osvResponse struct {
	Vulns []osvVuln `json:"vulns"`
}

type osvVuln struct {
	ID       string `json:"id"`
	Summary  string `json:"summary"`
	Details  string `json:"details"`
	Severity []struct {
		Type  string `json:"type"`
		Score string `json:"score"`
	} `json:"severity"`
	Affected []struct {
		Package struct {
			Name      string `json:"name"`
			Ecosystem string `json:"ecosystem"`
		} `json:"package"`
		Ranges []struct {
			Type   string `json:"type"`
			Events []struct {
				Introduced string `json:"introduced"`
				Fixed      string `json:"fixed"`
			} `json:"events"`
		} `json:"ranges"`
	} `json:"affected"`
	References []struct {
		Type string `json:"type"`
		URL  string `json:"url"`
	} `json:"references"`
}

type SCAnalyzer struct {
	client *http.Client
}

func NewSCAnalyzer() *SCAnalyzer {
	return &SCAnalyzer{
		client: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

func (a *SCAnalyzer) Scan(deps []types.ResolvedDependency, ecosystem string) []types.SecurityFinding {
	var findings []types.SecurityFinding

	for _, dep := range deps {
		if dep.Scope == "test" || dep.Scope == "provided" {
			continue
		}
		if dep.ResolvedVersion == "" {
			continue
		}

		vulns := a.queryOSV(dep.ArtifactID, dep.ResolvedVersion, ecosystem)
		for _, v := range vulns {
			severity := classifySeverity(v)
			fixedIn := findFixedVersion(v, dep.ResolvedVersion)

			finding := types.SecurityFinding{
				ID:          v.ID,
				Title:       v.Summary,
				Severity:    severity,
				Package:     dep.ArtifactID,
				Version:     dep.ResolvedVersion,
				FixedIn:     fixedIn,
				Impact:      v.Details,
				Remediation: generateRemediation(dep, fixedIn),
				Source:      "OSV.dev",
			}

			for _, ref := range v.References {
				if strings.Contains(ref.URL, "github.com/advisories") || strings.Contains(ref.URL, "nvd.nist.gov") {
					finding.URL = ref.URL
					break
				}
			}

			if finding.URL == "" && len(v.References) > 0 {
				finding.URL = v.References[0].URL
			}

			findings = append(findings, finding)
		}
	}

	return findings
}

func (a *SCAnalyzer) queryOSV(name, version, ecosystem string) []osvVuln {
	query := osvQuery{
		Version: version,
	}
	query.Package.Name = name
	query.Package.Ecosystem = ecosystem

	body, err := json.Marshal(query)
	if err != nil {
		return nil
	}

	resp, err := a.client.Post(osvAPIURL, "application/json", strings.NewReader(string(body)))
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil
	}

	var osvResp osvResponse
	if err := json.Unmarshal(respBody, &osvResp); err != nil {
		return nil
	}

	return osvResp.Vulns
}

func classifySeverity(vuln osvVuln) types.Severity {
	summary := strings.ToLower(vuln.Summary)
	if strings.Contains(summary, "remote code execution") ||
		strings.Contains(summary, "rce") ||
		strings.Contains(summary, "arbitrary code") ||
		strings.Contains(summary, "sql injection") ||
		strings.Contains(summary, "authentication bypass") {
		return types.SeverityCritical
	}

	if strings.Contains(summary, "denial of service") ||
		strings.Contains(summary, "information disclosure") ||
		strings.Contains(summary, "cross-site scripting") ||
		strings.Contains(summary, "xss") {
		return types.SeverityMedium
	}

	if strings.Contains(summary, "overflow") ||
		strings.Contains(summary, "injection") ||
		strings.Contains(summary, "deserialization") {
		return types.SeverityHigh
	}

	return types.SeverityMedium
}

func findFixedVersion(vuln osvVuln, currentVersion string) string {
	for _, affected := range vuln.Affected {
		for _, r := range affected.Ranges {
			for _, event := range r.Events {
				if event.Fixed != "" {
					return event.Fixed
				}
			}
		}
	}
	return ""
}

func generateRemediation(dep types.ResolvedDependency, fixedIn string) string {
	if fixedIn != "" {
		return fmt.Sprintf("Upgrade %s:%s from %s to %s or later",
			dep.GroupID, dep.ArtifactID, dep.ResolvedVersion, fixedIn)
	}
	return fmt.Sprintf("Upgrade %s:%s from %s to the latest available version",
		dep.GroupID, dep.ArtifactID, dep.ResolvedVersion)
}
