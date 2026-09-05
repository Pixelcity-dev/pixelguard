package parser

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/PixelCity-dev/pixelguard/internal/types"
)

type MavenProject struct {
	XMLName      xml.Name          `xml:"project"`
	GroupID      string            `xml:"groupID"`
	ArtifactID   string            `xml:"artifactID"`
	Version      string            `xml:"version"`
	Packaging    string            `xml:"packaging"`
	Properties   MavenProperties   `xml:"properties"`
	Modules      []string          `xml:"modules>module"`
	Dependencies []Dependency      `xml:"dependencies>dependency"`
	Parent       *Parent           `xml:"parent"`
	Build        *Build            `xml:"build"`
}

type MavenProperties struct {
	Props map[string]string
}

func (mp *MavenProperties) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	mp.Props = make(map[string]string)
	for {
		token, err := d.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		if end, ok := token.(xml.EndElement); ok && end.Name == start.Name {
			break
		}
		if se, ok := token.(xml.StartElement); ok {
			var val string
			if err := d.DecodeElement(&val, &se); err != nil {
				return err
			}
			mp.Props[se.Name.Local] = val
		}
	}
	return nil
}

type Parent struct {
	GroupID    string `xml:"groupID"`
	ArtifactID string `xml:"artifactID"`
	Version    string `xml:"version"`
}

type Build struct {
	Plugins []Plugin `xml:"plugins>plugin"`
}

type Plugin struct {
	GroupID    string `xml:"groupID"`
	ArtifactID string `xml:"artifactID"`
	Version    string `xml:"version"`
}

type Dependency struct {
	GroupID    string `xml:"groupId"`
	ArtifactID string `xml:"artifactId"`
	Version    string `xml:"version"`
	Scope      string `xml:"scope"`
}

func DetectProjectType(dir string) types.ProjectType {
	pomPath := filepath.Join(dir, "pom.xml")
	if _, err := os.Stat(pomPath); err == nil {
		return types.ProjectTypeMaven
	}

	buildGradle := filepath.Join(dir, "build.gradle")
	if _, err := os.Stat(buildGradle); err == nil {
		return types.ProjectTypeGradle
	}

	settingsGradle := filepath.Join(dir, "settings.gradle")
	if _, err := os.Stat(settingsGradle); err == nil {
		return types.ProjectTypeGradle
	}

	return types.ProjectTypeUnknown
}

func ParseMaven(dir string) (*types.ProjectInfo, []types.Dependency, error) {
	pomPath := filepath.Join(dir, "pom.xml")
	data, err := os.ReadFile(pomPath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read pom.xml: %w", err)
	}

	var project MavenProject
	if err := xml.Unmarshal(data, &project); err != nil {
		return nil, nil, fmt.Errorf("failed to parse pom.xml: %w", err)
	}

	javaVersion := resolveJavaVersion(project.Properties.Props)
	sourceFiles := countSourceFiles(dir, ".java")

	directDeps := make([]types.Dependency, 0, len(project.Dependencies))
	for _, dep := range project.Dependencies {
		version := resolveProperty(dep.Version, project.Properties.Props)
		scope := dep.Scope
		if scope == "" {
			scope = "compile"
		}
		directDeps = append(directDeps, types.Dependency{
			GroupID:    dep.GroupID,
			ArtifactID: dep.ArtifactID,
			Version:    version,
			Scope:      scope,
		})
	}

	info := &types.ProjectInfo{
		Type:            types.ProjectTypeMaven,
		Language:        "Java",
		JavaVersion:     javaVersion,
		SourceFiles:     sourceFiles,
		DirectDeps:      len(directDeps),
		ResolvedDeps:    len(directDeps),
		Modules:         project.Modules,
	}

	return info, directDeps, nil
}

func ParseGradle(dir string) (*types.ProjectInfo, []types.Dependency, error) {
	gradlePath := filepath.Join(dir, "build.gradle")
	if _, err := os.Stat(gradlePath); err != nil {
		gradlePath = filepath.Join(dir, "build.gradle.kts")
		if _, err := os.Stat(gradlePath); err != nil {
			return nil, nil, fmt.Errorf("no build.gradle or build.gradle.kts found")
		}
	}

	data, err := os.ReadFile(gradlePath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read build.gradle: %w", err)
	}

	content := string(data)
	deps := parseGradleDeps(content)
	javaVersion := parseGradleJavaVersion(content)
	sourceFiles := countSourceFiles(dir, ".java")

	info := &types.ProjectInfo{
		Type:        types.ProjectTypeGradle,
		Language:    "Java",
		JavaVersion: javaVersion,
		SourceFiles: sourceFiles,
		DirectDeps:  len(deps),
	}

	return info, deps, nil
}

func parseGradleDeps(content string) []types.Dependency {
	var deps []types.Dependency
	lines := strings.Split(content, "\n")

	inDepsBlock := false
	depth := 0

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		if strings.Contains(trimmed, "dependencies") && strings.Contains(trimmed, "{") {
			inDepsBlock = true
			depth = 1
			continue
		}

		if inDepsBlock {
			depth += strings.Count(trimmed, "{") - strings.Count(trimmed, "}")
			if depth <= 0 {
				inDepsBlock = false
				continue
			}

			if strings.Contains(trimmed, "implementation") ||
				strings.Contains(trimmed, "compileOnly") ||
				strings.Contains(trimmed, "testImplementation") ||
				strings.Contains(trimmed, "compile") ||
				strings.Contains(trimmed, "api") ||
				strings.Contains(trimmed, "runtimeOnly") {
				dep := parseGradleDepLine(trimmed)
				if dep != nil {
					deps = append(deps, *dep)
				}
			}
		}
	}

	return deps
}

func parseGradleDepLine(line string) *types.Dependency {
	line = strings.TrimRight(line, "}")
	line = strings.TrimSpace(line)

	start := strings.Index(line, "'")
	if start == -1 {
		start = strings.Index(line, "\"")
	}
	if start == -1 {
		return nil
	}

	end := strings.LastIndex(line, "'")
	if end == -1 {
		end = strings.LastIndex(line, "\"")
	}
	if end == -1 || end <= start {
		return nil
	}

	coord := line[start+1 : end]
	parts := strings.Split(coord, ":")

	if len(parts) < 2 {
		return nil
	}

	dep := &types.Dependency{
		GroupID:    parts[0],
		ArtifactID: parts[1],
	}
	if len(parts) >= 3 {
		dep.Version = parts[2]
	}

	scope := "compile"
	if strings.Contains(line, "implementation") {
		scope = "compile"
	} else if strings.Contains(line, "compileOnly") {
		scope = "provided"
	} else if strings.Contains(line, "testImplementation") {
		scope = "test"
	} else if strings.Contains(line, "runtimeOnly") {
		scope = "runtime"
	}
	dep.Scope = scope

	return dep
}

func parseGradleJavaVersion(content string) string {
	markers := []string{"sourceCompatibility", "targetCompatibility", "JavaVersion.VERSION"}
	for _, marker := range markers {
		idx := strings.Index(content, marker)
		if idx != -1 {
			rest := content[idx:]
			eqIdx := strings.Index(rest, "=")
			if eqIdx != -1 {
				val := strings.TrimSpace(rest[eqIdx+1:])
				val = strings.TrimRight(val, "\n\r")
				parts := strings.Fields(val)
				if len(parts) > 0 {
					version := strings.Trim(parts[0], "'\"")
					if strings.Contains(version, "VERSION_") {
						version = strings.Replace(version, "VERSION_", "", 1)
						version = strings.Replace(version, "_", ".", -1)
					}
					return version
				}
			}
		}
	}
	return ""
}

func resolveJavaVersion(props map[string]string) string {
	for _, key := range []string{"java.version", "maven.compiler.source", "java.source.version"} {
		if v, ok := props[key]; ok {
			return v
		}
	}
	return ""
}

func resolveProperty(val string, props map[string]string) string {
	if !strings.HasPrefix(val, "${") {
		return val
	}
	key := strings.Trim(val, "${}")
	if v, ok := props[key]; ok {
		return v
	}
	return val
}

func countSourceFiles(dir, ext string) int {
	count := 0
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() {
			name := info.Name()
			if name == "target" || name == "build" || name == ".git" || name == "node_modules" {
				return filepath.SkipDir
			}
			return nil
		}
		if strings.HasSuffix(path, ext) {
			count++
		}
		return nil
	})
	return count
}
