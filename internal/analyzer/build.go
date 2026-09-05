package analyzer

import (
	"bufio"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/PixelCity-dev/pixelguard/internal/types"
)

var (
	javaErrorPatterns = []*regexp.Regexp{
		regexp.MustCompile(`^.*\.java:(\d+):\s*error:\s*(.+)$`),
		regexp.MustCompile(`^\[ERROR\]\s+(.+\.java):(\d+)(?::(\d+))?\s*(.*)$`),
		regexp.MustCompile(`^error: (.+)$`),
		regexp.MustCompile(`cannot find symbol`),
		regexp.MustCompile(`method does not override or implement a method from a supertype`),
		regexp.MustCompile(`incompatible types: (.+) cannot be converted to (.+)`),
		regexp.MustCompile(`(.+) is not abstract and does not override abstract method (.+)`),
		regexp.MustCompile(`non-static method (.+) cannot be referenced from a static context`),
		regexp.MustCompile(`cannot access (.+) in (.+)`),
		regexp.MustCompile(`package (.+) does not exist`),
		regexp.MustCompile(`class (.+) is public, should be declared in a file named (.+)`),
		regexp.MustCompile(`unreported exception (.+); must be caught or declared to be thrown`),
		regexp.MustCompile(`illegal start of expression`),
		regexp.MustCompile(` ';' expected`),
		regexp.MustCompile(` reached end of file while parsing`),
		regexp.MustCompile(`class (.+) extends ([^\s]+) is abstract`),
	}

	javaWarningPatterns = []*regexp.Regexp{
		regexp.MustCompile(`^\[WARNING\]\s+(.+\.java):(\d+)(?::(\d+))?\s*(.*)$`),
		regexp.MustCompile(`warning: (.+)`),
		regexp.MustCompile(`(.+) is deprecated`),
		regexp.MustCompile(`(.+) unused`),
	}

	gradleErrorPatterns = []*regexp.Regexp{
		regexp.MustCompile(`> Compilation failed with (\d+) errors?`),
		regexp.MustCompile(` FAILURE: Build failed with an exception\.`),
		regexp.MustCompile(`\* What went wrong:\s*$`),
		regexp.MustCompile(`(.+\.java):(\d+):\s*(.+)$`),
		regexp.MustCompile(`error: (.+)`),
		regexp.MustCompile(`cannot find symbol`),
		regexp.MustCompile(`package (.+) does not exist`),
	}

	failedPattern = regexp.MustCompile(`BUILD FAILURE|BUILD FAILED`)
)

type BuildLogAnalyzer struct {
	dir string
}

func NewBuildLogAnalyzer(dir string) *BuildLogAnalyzer {
	return &BuildLogAnalyzer{dir: dir}
}

func (a *BuildLogAnalyzer) Analyze() ([]types.BuildError, int) {
	var errors []types.BuildError
	warnings := 0

	// Try Maven build log
	mavenLog := filepath.Join(a.dir, "target", "build.log")
	if _, err := os.Stat(mavenLog); err == nil {
		e, w := a.parseBuildLog(mavenLog, "maven")
		errors = append(errors, e...)
		warnings += w
	}

	// Try Gradle build log
	gradleLog := filepath.Join(a.dir, "build", "reports")
	if _, err := os.Stat(gradleLog); err == nil {
		filepath.Walk(gradleLog, func(path string, info os.FileInfo, err error) error {
			if err != nil || info.IsDir() {
				return nil
			}
			if strings.HasSuffix(path, ".log") || strings.HasSuffix(path, ".txt") {
				e, w := a.parseBuildLog(path, "gradle")
				errors = append(errors, e...)
				warnings += w
			}
			return nil
		})
	}

	for _, logFile := range []string{"build.log", "mvn.log", "gradle.log"} {
		path := filepath.Join(a.dir, logFile)
		if _, err := os.Stat(path); err == nil {
			e, w := a.parseBuildLog(path, "generic")
			errors = append(errors, e...)
			warnings += w
		}
	}

	return errors, warnings
}

func (a *BuildLogAnalyzer) parseBuildLog(path, buildType string) ([]types.BuildError, int) {
	file, err := os.Open(path)
	if err != nil {
		return nil, 0
	}
	defer file.Close()

	var errors []types.BuildError
	warnings := 0
	scanner := bufio.NewScanner(file)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)

		for _, pattern := range javaErrorPatterns {
			matches := pattern.FindStringSubmatch(line)
			if matches != nil {
				errLine := 0
				if len(matches) >= 2 {
					if n, e := strconv.Atoi(matches[1]); e == nil {
						errLine = n
					}
				}
				msg := trimmed
				if len(matches) >= 3 {
					msg = matches[len(matches)-1]
				}

				errors = append(errors, types.BuildError{
					File:     path,
					Line:     errLine,
					Message:  msg,
					Severity: types.SeverityHigh,
					Raw:      line,
				})
				break
			}
		}

		for _, pattern := range gradleErrorPatterns {
			matches := pattern.FindStringSubmatch(line)
			if matches != nil {
				errLine := 0
				if len(matches) >= 3 {
					if n, e := strconv.Atoi(matches[2]); e == nil {
						errLine = n
					}
				}
				msg := trimmed
				if len(matches) >= 4 {
					msg = matches[3]
				}

				errors = append(errors, types.BuildError{
					File:     path,
					Line:     errLine,
					Message:  msg,
					Severity: types.SeverityHigh,
					Raw:      line,
				})
				break
			}
		}

		for _, pattern := range javaWarningPatterns {
			if pattern.MatchString(line) {
				warnings++
				break
			}
		}

		if failedPattern.MatchString(line) {
			errors = append(errors, types.BuildError{
				File:     path,
				Message:  "Build failed",
				Severity: types.SeverityCritical,
				Raw:      line,
			})
		}
	}

	return errors, warnings
}
