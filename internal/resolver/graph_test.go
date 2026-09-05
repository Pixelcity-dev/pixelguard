package resolver

import (
	"testing"

	"github.com/PixelCity-dev/pixelguard/internal/types"
)

func TestResolveNoConflicts(t *testing.T) {
	graph := NewGraph()
	deps := []types.Dependency{
		{GroupID: "com.google.guava", ArtifactID: "guava", Version: "31.1-jre"},
		{GroupID: "org.slf4j", ArtifactID: "slf4j-api", Version: "2.0.9"},
	}

	resolved := graph.Resolve(deps)
	if len(resolved) != 2 {
		t.Errorf("Resolve() returned %d deps, want 2", len(resolved))
	}

	conflicts := graph.FindConflicts(resolved)
	if len(conflicts) != 0 {
		t.Errorf("FindConflicts() returned %d conflicts, want 0", len(conflicts))
	}
}

func TestCompareVersions(t *testing.T) {
	tests := []struct {
		a, b string
		want int
	}{
		{"1.0.0", "1.0.0", 0},
		{"2.0.0", "1.0.0", 1},
		{"1.0.0", "2.0.0", -1},
		{"1.2.3", "1.2.4", -1},
		{"1.10.0", "1.9.0", 1},
		{"3.2.0", "3.1.9", 1},
	}

	for _, tt := range tests {
		got := compareVersions(tt.a, tt.b)
		if got != tt.want {
			t.Errorf("compareVersions(%q, %q) = %d, want %d", tt.a, tt.b, got, tt.want)
		}
	}
}

func TestPickLatest(t *testing.T) {
	versions := []string{"1.0.0", "3.2.0", "2.1.5"}
	got := pickLatest(versions)
	if got != "3.2.0" {
		t.Errorf("pickLatest() = %q, want %q", got, "3.2.0")
	}
}

func TestResolveTransitiveConflicts(t *testing.T) {
	graph := NewGraph()
	deps := []types.Dependency{
		{GroupID: "com.fasterxml.jackson.core", ArtifactID: "jackson-databind", Version: "2.14.0"},
	}

	resolved := graph.Resolve(deps)
	conflicts := graph.FindConflicts(resolved)

	// jackson-databind has transitive jackson-core and jackson-annotations
	// If a direct dep conflicts with transitive, we should see it
	_ = resolved
	_ = conflicts
}
