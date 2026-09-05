<p align="center">
  <img src="https://github.com/PixelCity-dev/pixelguard/blob/master/.github/assets/hero.png?raw=true" alt="PixelGuard - Interactive Project Auditor" width="100%">
</p>

<h1 align="center">PixelGuard</h1>

<p align="center">
  <strong>Interactive project auditor — build errors, dependency conflicts, and security findings.</strong>
</p>

<p align="center">
  <a href="https://github.com/PixelCity-dev/pixelguard/stargazers"><img src="https://img.shields.io/github/stars/PixelCity-dev/pixelguard?style=flat-square&logo=github&color=yellow" alt="Stars"></a>
  <a href="https://github.com/PixelCity-dev/pixelguard/network/members"><img src="https://img.shields.io/github/forks/PixelCity-dev/pixelguard?style=flat-square&logo=github&color=blue" alt="Forks"></a>
  <a href="https://github.com/PixelCity-dev/pixelguard/issues"><img src="https://img.shields.io/github/issues/PixelCity-dev/pixelguard?style=flat-square&logo=github&color=red" alt="Issues"></a>
  <a href="https://github.com/PixelCity-dev/pixelguard/pulls"><img src="https://img.shields.io/github/issues-pr/PixelCity-dev/pixelguard?style=flat-square&logo=github&color=green" alt="Pull Requests"></a>
  <a href="https://github.com/PixelCity-dev/pixelguard/blob/master/LICENSE"><img src="https://img.shields.io/github/license/PixelCity-dev/pixelguard?style=flat-square&color=purple" alt="License"></a>
  <a href="https://github.com/PixelCity-dev/pixelguard/releases"><img src="https://img.shields.io/github/v/release/PixelCity-dev/pixelguard?style=flat-square&logo=github&color=orange" alt="Release"></a>
</p>

<p align="center">
  <img src="https://img.shields.io/badge/Language-Go-00ADD8?style=flat-square&logo=go&logoColor=white" alt="Go">
  <img src="https://img.shields.io/badge/Platform-Linux%20%7C%20macOS%20%7C%20Windows-2EA44F?style=flat-square&logo=linux&logoColor=white" alt="Platform">
  <img src="https://img.shields.io/badge/Build%20System-Maven%20%7C%20Gradle-FF5722?style=flat-square&logo=apachemaven&logoColor=white" alt="Build System">
  <img src="https://img.shields.io/badge/Security-OSV.dev-orange?style=flat-square&logo=security&logoColor=white" alt="Security">
</p>

<br>

<p align="center">
  <a href="https://github.com/PixelCity-dev/pixelguard">
    <img src="https://img.shields.io/badge/⭐_Star_History-Click_to_View-FFD700?style=for-the-badge&logo=github" alt="Star History">
  </a>
</p>

<p align="center">
  <a href="https://star-history.com/#PixelCity-dev/pixelguard&Date">
    <picture>
      <source media="(prefers-color-scheme: dark)" srcset="https://api.star-history.com/svg?repos=PixelCity-dev/pixelguard&type=Date&theme=dark" />
      <img alt="Star History Chart" src="https://api.star-history.com/svg?repos=PixelCity-dev/pixelguard&type=Date&theme=light" width="600" />
    </picture>
  </a>
</p>

<br>

<p align="center">
  Language: 
  <a href="#">English</a> | 
  <a href="#">Português (Brasil)</a> | 
  <a href="#">简体中文</a> | 
  <a href="#">繁體中文</a> | 
  <a href="#">日本語</a> | 
  <a href="#">한국어</a> | 
  <a href="#">Türkçe</a> | 
  <a href="#">Русский</a> | 
  <a href="#">Tiếng Việt</a> | 
  <a href="#">ไทย</a> | 
  <a href="#">Deutsch</a> | 
  <a href="#">Español</a> | 
  <a href="#">Українська</a>
</p>

<br>

<p align="center">
  <a href="https://github.com/PixelCity-dev/pixelguard">
    <img src="https://img.shields.io/badge/GitHub-PixelCity--dev/pixelguard-181717?style=for-the-badge&logo=github" alt="GitHub">
  </a>
  <a href="https://github.com/PixelCity-dev/pixelguard/blob/master/LICENSE">
    <img src="https://img.shields.io/badge/License-MIT-blue?style=for-the-badge" alt="MIT License">
  </a>
</p>

<br>

<p align="center">
  <img src="https://img.shields.io/badge/📦_Go_Language-00ADD8?style=flat-square&logo=go" alt="Go">
  <img src="https://img.shields.io/badge/🔧_Cobra_CLI-FF5722?style=flat-square&logo=cobra" alt="Cobra">
  <img src="https://img.shields.io/badge/💻_Bubbletea_TUI-9B59B6?style=flat-square&logo=bubbletea" alt="Bubbletea">
  <img src="https://img.shields.io/badge/🛡️_OSV.dev-2EA44F?style=flat-square&logo=shield" alt="OSV">
  <img src="https://img.shields.io/badge/📊_SARIF-3498DB?style=flat-square&logo=sarif" alt="SARIF">
</p>

---

> **100% Go** | **Maven + Gradle** | **OSV.dev SCA** | **Interactive TUI** | **JSON / Markdown / SARIF export**

---

**The interactive project auditor for Java developers. Built from real-world security engineering workflows.**

Not just a scanner. A complete system: project parsing, dependency graph resolution, build log analysis, security composition analysis, and an interactive drill-down interface. One scan, explore forever.

Point it at a **Maven** or **Gradle** project and get a single, trustworthy report covering:

1. **Build errors** — what's broken right now
2. **Dependency conflicts** — diamond dependencies, version mismatches
3. **Security exposure** — CVEs ranked by real severity, sourced from OSV.dev

Scan once, then explore findings interactively without re-running.

---

## Features

| Feature | Description |
|---------|-------------|
| **Full project ingestion** | Reads `pom.xml`, `build.gradle`, source tree, and build logs |
| **Dependency graph resolution** | Surfaces diamond dependency problems and transitive conflicts |
| **Security composition analysis (SCA)** | Cross-references dependencies against OSV.dev vulnerability feed |
| **Build log analysis** | Parses Maven/Gradle compiler output for real errors |
| **Interactive TUI** | Live scan progress, then queryable findings with drill-down |
| **Exportable reports** | JSON, Markdown, SARIF — CI/CD pipeline ready |
| **AI analysis agent** | System prompt for LLM-powered synthesis (see `AGENT_PROMPT.md`) |

---

## Install

```bash
go install github.com/PixelCity-dev/pixelguard/cmd/pixelguard@latest
```

Or build from source:

```bash
git clone https://github.com/PixelCity-dev/pixelguard.git
cd pixelguard
go build -o pixelguard ./cmd/pixelguard
```

---

## Usage

### Interactive scan (launches TUI)

```bash
pixelguard ./my-java-project
```

### Export reports (non-interactive)

```bash
pixelguard export ./my-java-project
```

Outputs:
- `pixelguard-report.json` — machine-readable
- `pixelguard-report.md` — human-readable
- `pixelguard-report.sarif` — CI/CD compatible

### Flags

```
--skip-security    Skip OSV.dev vulnerability scan
--skip-build       Skip build log analysis
-v, --verbose      Verbose output
```

---

## TUI Controls

| Key | Action |
|-----|--------|
| `↑` / `k` | Move selection up |
| `↓` / `j` | Move selection down |
| `Enter` / `Space` | View finding detail |
| `1` | Switch to Summary view |
| `2` | Switch to Findings view |
| `h` / `?` | Toggle help screen |
| `Esc` | Go back / close detail |
| `q` / `Ctrl+C` | Quit |

---

## Architecture

```
pixelguard/
├── cmd/pixelguard/main.go         # CLI entry point (cobra)
├── internal/
│   ├── types/                     # Shared types (Severity, Dependency, Report)
│   ├── parser/
│   │   ├── parser.go              # Maven pom.xml + Gradle build.gradle parsers
│   │   └── parser_test.go         # Parser tests
│   ├── resolver/
│   │   ├── graph.go               # Dependency graph resolver + transitive simulation
│   │   └── graph_test.go          # Resolver tests
│   ├── analyzer/
│   │   ├── build.go               # Build log analyzer (Java/Gradle error patterns)
│   │   └── security.go            # OSV.dev SCA connector
│   ├── engine/
│   │   └── engine.go              # Orchestrator (progress callbacks, phased scan)
│   ├── report/
│   │   └── report.go              # Export: JSON, Markdown, SARIF
│   └── ui/
│       ├── banner.go              # Colored CLI output helpers
│       └── tui.go                 # Interactive Bubbletea TUI
└── AGENT_PROMPT.md                # AI analysis agent system prompt
```

---

## AI Analysis Agent

PixelGuard includes a system prompt for an AI analysis engine in [`AGENT_PROMPT.md`](AGENT_PROMPT.md). It's designed to synthesize scanner outputs into coherent, prioritized reports — without inventing CVEs or version numbers from memory.

Wire it into any tool-calling LLM (Claude, GPT, etc.) by feeding the scan's structured output as context.

---

## Example Terminal Session

```
$ pixelguard scan ./my-java-project

 Scanning project structure...
✔ Detected: Maven project (Java 17, pom.xml found)
✔ 187 source files, 42 declared dependencies, 118 resolved (incl. transitive)

Build:        3 errors, 12 warnings
Conflicts:    5 dependency version conflicts
Security:     7 findings — 2 Critical, 3 High, 2 Medium
```

---

## Contributing

Contributions welcome! Open an issue or submit a PR.

## License

MIT License — see [LICENSE](LICENSE) for details.

---

<p align="center">
  Made with ❤️ by <a href="https://github.com/jaylukdevelopment">Jay</a> & <a href="https://github.com/PixelCity-dev">PixelCity</a>
</p>

<p align="center">
  <a href="https://github.com/PixelCity-dev/pixelguard/stargazers">
    <img src="https://img.shields.io/github/stars/PixelCity-dev/pixelguard?style=social" alt="Star PixelGuard">
  </a>
</p>
