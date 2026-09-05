package types

import "time"

type Severity int

const (
	SeverityLow Severity = iota
	SeverityMedium
	SeverityHigh
	SeverityCritical
)

func (s Severity) String() string {
	switch s {
	case SeverityLow:
		return "Low"
	case SeverityMedium:
		return "Medium"
	case SeverityHigh:
		return "High"
	case SeverityCritical:
		return "Critical"
	default:
		return "Unknown"
	}
}

func (s Severity) Symbol() string {
	switch s {
	case SeverityCritical:
		return "🔴"
	case SeverityHigh:
		return "🟠"
	case SeverityMedium:
		return "🟡"
	case SeverityLow:
		return "🟢"
	default:
		return "⚪"
	}
}

type ProjectType int

const (
	ProjectTypeMaven ProjectType = iota
	ProjectTypeGradle
	ProjectTypeUnknown
)

func (p ProjectType) String() string {
	switch p {
	case ProjectTypeMaven:
		return "Maven"
	case ProjectTypeGradle:
		return "Gradle"
	default:
		return "Unknown"
	}
}

type Dependency struct {
	GroupID    string `xml:"groupId" json:"group_id"`
	ArtifactID string `xml:"artifactId" json:"artifact_id"`
	Version    string `xml:"version" json:"version"`
	Scope      string `xml:"scope" json:"scope,omitempty"`
}

type ResolvedDependency struct {
	Dependency
	ResolvedVersion string   `json:"resolved_version"`
	ConflictsWith   []string `json:"conflicts_with,omitempty"`
	IsTransitive    bool     `json:"is_transitive"`
	Parent          string   `json:"parent,omitempty"`
}

type BuildError struct {
	File     string   `json:"file"`
	Line     int      `json:"line,omitempty"`
	Column   int      `json:"column,omitempty"`
	Message  string   `json:"message"`
	Severity Severity `json:"severity"`
	Raw      string   `json:"raw"`
}

type SecurityFinding struct {
	ID          string   `json:"id"`
	Title       string   `json:"title"`
	Severity    Severity `json:"severity"`
	CVE         string   `json:"cve,omitempty"`
	Advisory    string   `json:"advisory,omitempty"`
	Package     string   `json:"package"`
	Version     string   `json:"version"`
	FixedIn     string   `json:"fixed_in,omitempty"`
	Impact      string   `json:"impact"`
	Remediation string   `json:"remediation"`
	Source      string   `json:"source"`
	URL         string   `json:"url,omitempty"`
}

type DependencyConflict struct {
	Package            string   `json:"package"`
	Versions           []string `json:"versions"`
	PulledBy           []string `json:"pulled_by"`
	RecommendedVersion string   `json:"recommended_version"`
	Severity           Severity `json:"severity"`
}

type ProjectInfo struct {
	Type             ProjectType `json:"type"`
	Language         string      `json:"language"`
	JavaVersion      string      `json:"java_version,omitempty"`
	SourceFiles      int         `json:"source_files"`
	DirectDeps       int         `json:"direct_deps"`
	ResolvedDeps     int         `json:"resolved_deps"`
	Modules          []string    `json:"modules,omitempty"`
	BuildToolVersion string      `json:"build_tool_version,omitempty"`
}

type Report struct {
	ProjectPath string              `json:"project_path"`
	ProjectInfo ProjectInfo         `json:"project_info"`
	ScanTime    time.Time           `json:"scan_time"`
	Duration    time.Duration       `json:"duration"`
	BuildErrors []BuildError        `json:"build_errors"`
	Conflicts   []DependencyConflict `json:"conflicts"`
	Security    []SecurityFinding   `json:"security_findings"`
	Summary     Summary             `json:"summary"`
}

type Summary struct {
	TotalBuildErrors int `json:"total_build_errors"`
	TotalWarnings    int `json:"total_warnings"`
	TotalConflicts   int `json:"total_conflicts"`
	TotalSecurity    int `json:"total_security"`
	CriticalCount    int `json:"critical_count"`
	HighCount        int `json:"high_count"`
	MediumCount      int `json:"medium_count"`
	LowCount         int `json:"low_count"`
}

type ScanConfig struct {
	ProjectPath  string
	MaxConcurrency int
	SkipSecurity bool
	SkipBuild    bool
	OutputFormat string
	OutputPath   string
	Verbose      bool
}

type ProgressEvent struct {
	Phase   string  `json:"phase"`
	Message string  `json:"message"`
	Percent float64 `json:"percent"`
	Done    bool    `json:"done"`
	Error   string  `json:"error,omitempty"`
}

type ProgressCallback func(event ProgressEvent)
