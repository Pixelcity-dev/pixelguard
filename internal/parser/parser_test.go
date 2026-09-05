package parser

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/PixelCity-dev/pixelguard/internal/types"
)

func TestDetectProjectTypeMaven(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "pom.xml"), []byte("<project/>"), 0644)

	got := DetectProjectType(dir)
	if got != types.ProjectTypeMaven {
		t.Errorf("DetectProjectType() = %v, want Maven", got)
	}
}

func TestDetectProjectTypeGradle(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "build.gradle"), []byte(""), 0644)

	got := DetectProjectType(dir)
	if got != types.ProjectTypeGradle {
		t.Errorf("DetectProjectType() = %v, want Gradle", got)
	}
}

func TestDetectProjectTypeUnknown(t *testing.T) {
	dir := t.TempDir()

	got := DetectProjectType(dir)
	if got != types.ProjectTypeUnknown {
		t.Errorf("DetectProjectType() = %v, want Unknown", got)
	}
}

func TestParseMaven(t *testing.T) {
	dir := t.TempDir()
	pom := `<?xml version="1.0" encoding="UTF-8"?>
<project xmlns="http://maven.apache.org/POM/4.0.0"
         xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
         xsi:schemaLocation="http://maven.apache.org/POM/4.0.0 http://maven.apache.org/xsd/maven-4.0.0.xsd">
    <modelVersion>4.0.0</modelVersion>
    <groupId>com.example</groupId>
    <artifactId>demo</artifactId>
    <version>1.0.0</version>
    <properties>
        <java.version>17</java.version>
    </properties>
    <dependencies>
        <dependency>
            <groupId>org.springframework.boot</groupId>
            <artifactId>spring-boot-starter-web</artifactId>
        </dependency>
        <dependency>
            <groupId>org.projectlombok</groupId>
            <artifactId>lombok</artifactId>
            <scope>provided</scope>
        </dependency>
    </dependencies>
</project>`
	os.WriteFile(filepath.Join(dir, "pom.xml"), []byte(pom), 0644)
	os.MkdirAll(filepath.Join(dir, "src", "main", "java"), 0755)
	os.WriteFile(filepath.Join(dir, "src", "main", "java", "App.java"), []byte("public class App {}"), 0644)

	info, deps, err := ParseMaven(dir)
	if err != nil {
		t.Fatalf("ParseMaven() error = %v", err)
	}

	if info.JavaVersion != "17" {
		t.Errorf("JavaVersion = %q, want %q", info.JavaVersion, "17")
	}
	if info.SourceFiles != 1 {
		t.Errorf("SourceFiles = %d, want 1", info.SourceFiles)
	}
	if len(deps) != 2 {
		t.Errorf("deps count = %d, want 2", len(deps))
	}
	if deps[0].ArtifactID != "spring-boot-starter-web" {
		t.Errorf("deps[0].ArtifactID = %q, want %q", deps[0].ArtifactID, "spring-boot-starter-web")
	}
}

func TestParseGradle(t *testing.T) {
	dir := t.TempDir()
	gradle := `plugins {
    id 'java'
}

sourceCompatibility = '17'

dependencies {
    implementation 'org.springframework.boot:spring-boot-starter-web:3.2.0'
    compileOnly 'org.projectlombok:lombok:1.18.30'
    testImplementation 'org.junit.jupiter:junit-jupiter:5.10.1'
}`
	os.WriteFile(filepath.Join(dir, "build.gradle"), []byte(gradle), 0644)

	info, deps, err := ParseGradle(dir)
	if err != nil {
		t.Fatalf("ParseGradle() error = %v", err)
	}

	if info.JavaVersion != "17" {
		t.Errorf("JavaVersion = %q, want %q", info.JavaVersion, "17")
	}
	if len(deps) != 3 {
		t.Errorf("deps count = %d, want 3", len(deps))
	}
}

func TestCountSourceFiles(t *testing.T) {
	dir := t.TempDir()
	os.MkdirAll(filepath.Join(dir, "src"), 0755)
	os.WriteFile(filepath.Join(dir, "src", "A.java"), []byte(""), 0644)
	os.WriteFile(filepath.Join(dir, "src", "B.java"), []byte(""), 0644)
	os.WriteFile(filepath.Join(dir, "src", "C.txt"), []byte(""), 0644)

	count := countSourceFiles(dir, ".java")
	if count != 2 {
		t.Errorf("countSourceFiles() = %d, want 2", count)
	}
}
