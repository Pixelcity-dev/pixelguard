# PixelGuard

**Interactive project auditor — build errors, dependency conflicts, and security findings.**

Created by [Jay](https://github.com/jaylukdevelopment) & [PixelCity](https://github.com/PixelCity-dev)

---

## What it does

Point PixelGuard at a Java project (Maven or Gradle) and get a single, trustworthy report covering:

1. **Build errors** — what's broken right now
2. **Dependency conflicts** — diamond dependencies, version mismatches
3. **Security exposure** — CVEs ranked by real severity, sourced from OSV.dev

Scan once, then explore findings interactively without re-running.

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

## TUI Controls

| Key | Action |
|-----|--------|
| `↑`/`k` | Move up |
| `↓`/`j` | Move down |
| `Enter` | View finding detail |
| `1` | Summary view |
| `2` | Findings view |
| `h` / `?` | Help |
| `Esc` | Back / close detail |
| `q` | Quit |

## Architecture

```
pixelguard/
├── cmd/pixelguard/main.go         # CLI (cobra)
└── internal/
    ├── types/                     # Shared types
    ├── parser/                    # Maven + Gradle parsers
    ├── resolver/                  # Dependency graph resolver
    ├── analyzer/
    │   ├── build.go               # Build log analyzer
    │   └── security.go            # OSV.dev SCA connector
    ├── engine/                    # Orchestrator
    ├── report/                    # Export: JSON, Markdown, SARIF
    └── ui/                        # Bubbletea TUI + CLI output
```

## AI Analysis Agent

PixelGuard includes a system prompt for an AI analysis engine in [`AGENT_PROMPT.md`](AGENT_PROMPT.md). It's designed to synthesize scanner outputs into coherent, prioritized reports — without inventing CVEs or version numbers from memory.

Wire it into any tool-calling LLM (Claude, GPT, etc.) by feeding the scan's structured output as context.

## License

MIT
