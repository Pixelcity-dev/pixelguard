package resolver

import (
	"github.com/PixelCity-dev/pixelguard/internal/types"
)

type Graph struct {
	Nodes map[string]*Node
}

type Node struct {
	Dependency types.ResolvedDependency
	Children   []*Node
}

func NewGraph() *Graph {
	return &Graph{
		Nodes: make(map[string]*Node),
	}
}

func (g *Graph) AddDependency(dep types.ResolvedDependency) {
	key := depKey(dep)
	g.Nodes[key] = &Node{
		Dependency: dep,
	}
}

func (g *Graph) Resolve(direct []types.Dependency) []types.ResolvedDependency {
	resolved := make(map[string]*types.ResolvedDependency)

	for _, dep := range direct {
		rd := types.ResolvedDependency{
			Dependency:      dep,
			ResolvedVersion: dep.Version,
			IsTransitive:    false,
		}
		key := g.nodeKey(dep)
		resolved[key] = &rd
	}

	transitive := g.simulateTransitive(direct)
	for _, dep := range transitive {
		key := g.nodeKey(dep)
		if existing, ok := resolved[key]; ok {
			if existing.ResolvedVersion != dep.Version {
				existing.ConflictsWith = append(existing.ConflictsWith, dep.Version)
			}
		} else {
			rd := types.ResolvedDependency{
				Dependency:      dep,
				ResolvedVersion: dep.Version,
				IsTransitive:    true,
			}
			resolved[key] = &rd
		}
	}

	result := make([]types.ResolvedDependency, 0, len(resolved))
	for _, rd := range resolved {
		result = append(result, *rd)
	}
	return result
}

func (g *Graph) FindConflicts(resolved []types.ResolvedDependency) []types.DependencyConflict {
	conflictMap := make(map[string]*types.DependencyConflict)

	for _, dep := range resolved {
		if len(dep.ConflictsWith) > 0 {
			key := dep.ArtifactID
			if existing, ok := conflictMap[key]; ok {
				for _, v := range dep.ConflictsWith {
					found := false
					for _, ev := range existing.Versions {
						if ev == v {
							found = true
							break
						}
					}
					if !found {
						existing.Versions = append(existing.Versions, v)
					}
				}
			} else {
				versions := []string{dep.ResolvedVersion}
				versions = append(versions, dep.ConflictsWith...)
				conflictMap[key] = &types.DependencyConflict{
					Package:  dep.ArtifactID,
					Versions: versions,
					PulledBy: []string{dep.Parent},
					Severity: types.SeverityMedium,
				}
			}
		}
	}

	conflicts := make([]types.DependencyConflict, 0, len(conflictMap))
	for _, c := range conflictMap {
		c.RecommendedVersion = pickLatest(c.Versions)
		if len(c.Versions) > 2 {
			c.Severity = types.SeverityHigh
		}
		conflicts = append(conflicts, *c)
	}
	return conflicts
}

func (g *Graph) nodeKey(dep types.Dependency) string {
	return dep.GroupID + ":" + dep.ArtifactID
}

func depKey(dep types.ResolvedDependency) string {
	return dep.GroupID + ":" + dep.ArtifactID + ":" + dep.ResolvedVersion
}

func pickLatest(versions []string) string {
	if len(versions) == 0 {
		return ""
	}
	latest := versions[0]
	for _, v := range versions[1:] {
		if compareVersions(v, latest) > 0 {
			latest = v
		}
	}
	return latest
}

func compareVersions(a, b string) int {
	aParts := splitVersion(a)
	bParts := splitVersion(b)

	maxLen := len(aParts)
	if len(bParts) > maxLen {
		maxLen = len(bParts)
	}

	for i := 0; i < maxLen; i++ {
		ai, bi := 0, 0
		if i < len(aParts) {
			ai = aParts[i]
		}
		if i < len(bParts) {
			bi = bParts[i]
		}
		if ai > bi {
			return 1
		}
		if ai < bi {
			return -1
		}
	}
	return 0
}

func splitVersion(v string) []int {
	var parts []int
	current := 0
	for _, c := range v {
		if c >= '0' && c <= '9' {
			current = current*10 + int(c-'0')
		} else {
			if current > 0 || len(parts) > 0 {
				parts = append(parts, current)
				current = 0
			}
		}
	}
	parts = append(parts, current)
	return parts
}

func (g *Graph) simulateTransitive(direct []types.Dependency) []types.Dependency {
	var transitive []types.Dependency

	commonTransitive := map[string][]types.Dependency{
		"spring-boot-starter-web": {
			{GroupID: "org.springframework.boot", ArtifactID: "spring-boot-starter", Version: "3.2.0"},
			{GroupID: "org.springframework.boot", ArtifactID: "spring-boot-starter-tomcat", Version: "3.2.0"},
			{GroupID: "com.fasterxml.jackson.core", ArtifactID: "jackson-databind", Version: "2.15.3"},
			{GroupID: "org.springframework", ArtifactID: "spring-web", Version: "6.1.1"},
			{GroupID: "org.yaml", ArtifactID: "snakeyaml", Version: "2.2"},
		},
		"spring-boot-starter-data-jpa": {
			{GroupID: "org.springframework.boot", ArtifactID: "spring-boot-starter-aop", Version: "3.2.0"},
			{GroupID: "org.hibernate.orm", ArtifactID: "hibernate-core", Version: "6.4.0.Final"},
			{GroupID: "org.springframework.data", ArtifactID: "spring-data-jpa", Version: "3.2.0"},
		},
		"spring-boot-starter-test": {
			{GroupID: "org.junit.jupiter", ArtifactID: "junit-jupiter", Version: "5.10.1"},
			{GroupID: "org.mockito", ArtifactID: "mockito-core", Version: "5.7.0"},
			{GroupID: "org.assertj", ArtifactID: "assertj-core", Version: "3.24.2"},
		},
		"logback-classic": {
			{GroupID: "ch.qos.logback", ArtifactID: "logback-core", Version: "1.4.14"},
			{GroupID: "org.slf4j", ArtifactID: "slf4j-api", Version: "2.0.9"},
		},
		"jackson-databind": {
			{GroupID: "com.fasterxml.jackson.core", ArtifactID: "jackson-core", Version: "2.15.3"},
			{GroupID: "com.fasterxml.jackson.core", ArtifactID: "jackson-annotations", Version: "2.15.3"},
		},
		"mockito-core": {
			{GroupID: "net.bytebuddy", ArtifactID: "byte-buddy", Version: "1.14.10"},
			{GroupID: "net.bytebuddy", ArtifactID: "byte-buddy-agent", Version: "1.14.10"},
			{GroupID: "org.objenesis", ArtifactID: "objenesis", Version: "3.3"},
		},
		"commons-io": {
			{GroupID: "commons-io", ArtifactID: "commons-io", Version: "2.15.0"},
		},
		"mysql-connector-j": {
			{GroupID: "com.google.protobuf", ArtifactID: "protobuf-java", Version: "3.24.3"},
		},
	}

	visited := make(map[string]bool)

	var walk func(deps []types.Dependency)
	walk = func(deps []types.Dependency) {
		for _, dep := range deps {
			key := dep.ArtifactID
			if visited[key] {
				continue
			}
			visited[key] = true

			if transitiveDeps, ok := commonTransitive[key]; ok {
				transitive = append(transitive, transitiveDeps...)
				walk(transitiveDeps)
			}
		}
	}

	walk(direct)
	return transitive
}
