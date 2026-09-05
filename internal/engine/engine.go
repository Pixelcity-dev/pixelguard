package engine

import (
	"os"
	"path/filepath"
	"time"

	"github.com/PixelCity-dev/pixelguard/internal/analyzer"
	"github.com/PixelCity-dev/pixelguard/internal/parser"
	"github.com/PixelCity-dev/pixelguard/internal/resolver"
	"github.com/PixelCity-dev/pixelguard/internal/types"
)

type Engine struct {
	config   types.ScanConfig
	callback types.ProgressCallback
}

func New(config types.ScanConfig, cb types.ProgressCallback) *Engine {
	if cb == nil {
		cb = func(event types.ProgressEvent) {}
	}
	return &Engine{
		config:   config,
		callback: cb,
	}
}

func (e *Engine) Scan() (*types.Report, error) {
	start := time.Now()
	report := &types.Report{
		ProjectPath: e.config.ProjectPath,
		ScanTime:    start,
	}

	// Phase 1: Detect project type
	e.emit(types.ProgressEvent{Phase: "detect", Message: "Detecting project type...", Percent: 0})

	projectType := parser.DetectProjectType(e.config.ProjectPath)
	if projectType == types.ProjectTypeUnknown {
		e.emitError("detect", "No pom.xml or build.gradle found in "+e.config.ProjectPath)
		return nil, os.ErrInvalid
	}

	e.emit(types.ProgressEvent{Phase: "detect", Message: "Detected " + projectType.String() + " project", Percent: 10, Done: true})

	// Phase 2: Parse project
	e.emit(types.ProgressEvent{Phase: "parse", Message: "Parsing build configuration...", Percent: 15})

	var info *types.ProjectInfo
	var deps []types.Dependency

	switch projectType {
	case types.ProjectTypeMaven:
		info2, deps2, err := parser.ParseMaven(e.config.ProjectPath)
		if err != nil {
			e.emitError("parse", err.Error())
			return nil, err
		}
		info = info2
		deps = deps2
	case types.ProjectTypeGradle:
		info2, deps2, err := parser.ParseGradle(e.config.ProjectPath)
		if err != nil {
			e.emitError("parse", err.Error())
			return nil, err
		}
		info = info2
		deps = deps2
	}

	// Phase 3: Resolve dependencies
	e.emit(types.ProgressEvent{Phase: "parse", Message: "Found " + itoa(info.DirectDeps) + " direct dependencies", Percent: 25, Done: true})
	report.ProjectInfo = *info

	e.emit(types.ProgressEvent{Phase: "resolve", Message: "Resolving dependency graph...", Percent: 30})

	graph := resolver.NewGraph()
	resolved := graph.Resolve(deps)
	report.ProjectInfo.ResolvedDeps = len(resolved)

	e.emit(types.ProgressEvent{Phase: "resolve", Message: "Resolved " + itoa(len(resolved)) + " dependencies", Percent: 45, Done: true})

	// Phase 4: Find conflicts
	e.emit(types.ProgressEvent{Phase: "conflict", Message: "Checking for dependency conflicts...", Percent: 50})

	conflicts := graph.FindConflicts(resolved)
	report.Conflicts = conflicts

	e.emit(types.ProgressEvent{Phase: "conflict", Message: "Found " + itoa(len(conflicts)) + " conflicts", Percent: 60, Done: true})

	// Phase 5: Build analysis
	if !e.config.SkipBuild {
		e.emit(types.ProgressEvent{Phase: "build", Message: "Analyzing build logs...", Percent: 65})

		buildAnalyzer := analyzer.NewBuildLogAnalyzer(e.config.ProjectPath)
		buildErrors, warnings := buildAnalyzer.Analyze()
		report.BuildErrors = buildErrors
		report.Summary.TotalWarnings = warnings

		e.emit(types.ProgressEvent{Phase: "build", Message: "Found " + itoa(len(buildErrors)) + " build errors", Percent: 75, Done: true})
	}

	// Phase 6: Security analysis
	if !e.config.SkipSecurity {
		e.emit(types.ProgressEvent{Phase: "security", Message: "Running security analysis...", Percent: 80})

		sca := analyzer.NewSCAnalyzer()
		ecosystem := "Maven"
		securityFindings := sca.Scan(resolved, ecosystem)
		report.Security = securityFindings

		e.emit(types.ProgressEvent{Phase: "security", Message: "Found " + itoa(len(securityFindings)) + " security findings", Percent: 95, Done: true})
	}

	// Phase 7: Compile summary
	report.Summary.TotalBuildErrors = len(report.BuildErrors)
	report.Summary.TotalConflicts = len(report.Conflicts)
	report.Summary.TotalSecurity = len(report.Security)

	for _, f := range report.Security {
		switch f.Severity {
		case types.SeverityCritical:
			report.Summary.CriticalCount++
		case types.SeverityHigh:
			report.Summary.HighCount++
		case types.SeverityMedium:
			report.Summary.MediumCount++
		case types.SeverityLow:
			report.Summary.LowCount++
		}
	}

	report.Duration = time.Since(start)

	e.emit(types.ProgressEvent{Phase: "done", Message: "Scan complete", Percent: 100, Done: true})

	return report, nil
}

func (e *Engine) emit(event types.ProgressEvent) {
	e.callback(event)
}

func (e *Engine) emitError(phase, msg string) {
	e.callback(types.ProgressEvent{
		Phase:   phase,
		Message: msg,
		Error:   msg,
		Done:    true,
	})
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	if n < 0 {
		return "-" + itoa(-n)
	}
	digits := make([]byte, 0, 10)
	for n > 0 {
		digits = append([]byte{byte('0' + n%10)}, digits...)
		n /= 10
	}
	return string(digits)
}

func AbsPath(path string) string {
	if abs, err := filepath.Abs(path); err == nil {
		return abs
	}
	return path
}
