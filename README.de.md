<p align="center">
  <img src="https://github.com/PixelCity-dev/pixelguard/blob/master/.github/assets/hero.png?raw=true" alt="PixelGuard - Interaktiver Projekt-Auditor" width="100%">
</p>

<h1 align="center">PixelGuard</h1>

<p align="center">
  <strong>Interaktiver Projekt-Auditor — Build-Fehler, Abhängigkeitskonflikte und Sicherheitslücken.</strong>
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
  <img src="https://img.shields.io/badge/Sprache-Go-00ADD8?style=flat-square&logo=go&logoColor=white" alt="Go">
  <img src="https://img.shields.io/badge/Plattform-Linux%20%7C%20macOS%20%7C%20Windows-2EA44F?style=flat-square&logo=linux&logoColor=white" alt="Plattform">
  <img src="https://img.shields.io/badge/Build-System-Maven%20%7C%20Gradle-FF5722?style=flat-square&logo=apachemaven&logoColor=white" alt="Build-System">
  <img src="https://img.shields.io/badge/Sicherheit-OSV.dev-orange?style=flat-square&logo=security&logoColor=white" alt="Sicherheit">
</p>

<br>

<p align="center">
  <a href="https://github.com/PixelCity-dev/pixelguard">
    <img src="https://img.shields.io/badge/⭐_Star_Historie-Ansehen-FFD700?style=for-the-badge&logo=github" alt="Star History">
  </a>
</p>

<p align="center">
  <a href="https://star-history.com/#PixelCity-dev/pixelguard&Date">
    <picture>
      <source media="(prefers-color-scheme: dark)" srcset="https://api.star-history.com/svg?repos=PixelCity-dev/pixelguard&type=Date&theme=dark" />
      <img alt="Star History Diagramm" src="https://api.star-history.com/svg?repos=PixelCity-dev/pixelguard&type=Date&theme=light" width="600" />
    </picture>
  </a>
</p>

<br>

<p align="center">
  Sprache: 
  <a href="https://github.com/PixelCity-dev/pixelguard/blob/master/README.md">English</a> | 
  <a href="https://github.com/PixelCity-dev/pixelguard/blob/master/README.de.md">Deutsch</a>
</p>

<br>

<p align="center">
  <a href="https://github.com/PixelCity-dev/pixelguard">
    <img src="https://img.shields.io/badge/GitHub-PixelCity--dev/pixelguard-181717?style=for-the-badge&logo=github" alt="GitHub">
  </a>
  <a href="https://github.com/PixelCity-dev/pixelguard/blob/master/LICENSE">
    <img src="https://img.shields.io/badge/Lizenz-MIT-blue?style=for-the-badge" alt="MIT Lizenz">
  </a>
</p>

<br>

<p align="center">
  <img src="https://img.shields.io/badge/📦_Go_Sprache-00ADD8?style=flat-square&logo=go" alt="Go">
  <img src="https://img.shields.io/badge/🔧_Cobra_CLI-FF5722?style=flat-square&logo=cobra" alt="Cobra">
  <img src="https://img.shields.io/badge/💻_Bubbletea_TUI-9B59B6?style=flat-square&logo=bubbletea" alt="Bubbletea">
  <img src="https://img.shields.io/badge/🛡️_OSV.dev-2EA44F?style=flat-square&logo=shield" alt="OSV">
  <img src="https://img.shields.io/badge/📊_SARIF-3498DB?style=flat-square&logo=sarif" alt="SARIF">
</p>

---

> **100% Go** | **Maven + Gradle** | **OSV.dev SCA** | **Interaktives TUI** | **JSON / Markdown / SARIF Export**

---

**Der interaktive Projekt-Auditor für Java-Entwickler. Gebaut aus realen Security-Engineering-Workflows.**

Nicht nur ein Scanner. Ein komplettes System: Projekt-Parsing, Abhängigkeitsgraph-Auflösung, Build-Log-Analyse, Security Composition Analysis und eine interaktive Drill-Down-Oberfläche. Einmal scannen, ewig erkunden.

Richte ihn auf ein **Maven**- oder **Gradle**-Projekt und erhalte einen einzigen, vertrauenswürdigen Bericht über:

1. **Build-Fehler** — was gerade kaputt ist
2. **Abhängigkeitskonflikte** — Diamond-Dependencies, Versionskonflikte
3. **Sicherheitslücken** — CVEs nach echtem Schweregrad sortiert, aus OSV.dev

Einmal scannen, dann die Ergebnisse interaktiv erkunden — ohne Neuaufruf.

---

## Funktionen

| Funktion | Beschreibung |
|----------|--------------|
| **Vollständige Projektaufnahme** | Liest `pom.xml`, `build.gradle`, Quellbaum und Build-Logs |
| **Abhängigkeitsgraph-Auflösung** | Zeigt Diamond-Dependency-Probleme und transitive Konflikte |
| **Security Composition Analysis (SCA)** | Kreuzreferenziert Abhängigkeiten mit der OSV.dev Schwachstelldatenbank |
| **Build-Log-Analyse** | Parst Maven/Gradle Compiler-Ausgaben auf echte Fehler |
| **Interaktives TUI** | Live-Fortschritt, dann abfragbare Ergebnisse mit Drill-Down |
| **Exportierbare Berichte** | JSON, Markdown, SARIF — CI/CD Pipeline bereit |
| **KI-Analyse-Agent** | System-Prompt für LLM-gestützte Synthese (siehe `AGENT_PROMPT.md`) |

---

## Installation

```bash
go install github.com/PixelCity-dev/pixelguard/cmd/pixelguard@latest
```

Oder aus Quelle bauen:

```bash
git clone https://github.com/PixelCity-dev/pixelguard.git
cd pixelguard
go build -o pixelguard ./cmd/pixelguard
```

---

## Verwendung

### Interaktiver Scan (startet TUI)

```bash
pixelguard ./mein-java-projekt
```

### Berichte exportieren (nicht-interaktiv)

```bash
pixelguard export ./mein-java-projekt
```

Ausgaben:
- `pixelguard-report.json` — maschinenlesbar
- `pixelguard-report.md` — menschenlesbar
- `pixelguard-report.sarif` — CI/CD kompatibel

### Optionen

```
--skip-security    OSV.dev Schwachstellenscan überspringen
--skip-build       Build-Log-Analyse überspringen
-v, --verbose      Ausführliche Ausgabe
```

---

## TUI-Steuerung

| Taste | Aktion |
|-------|--------|
| `↑` / `k` | Auswahl nach oben |
| `↓` / `j` | Auswahl nach unten |
| `Enter` / `Space` | Fundstelle im Detail anzeigen |
| `1` | Zur Zusammenfassungsansicht wechseln |
| `2` | Zur Fundstellenansicht wechseln |
| `h` / `?` | Hilfe ein-/ausschalten |
| `Esc` | Zurück / Detail schließen |
| `q` / `Ctrl+C` | Beenden |

---

## Architektur

```
pixelguard/
├── cmd/pixelguard/main.go         # CLI-Einstiegspunkt (cobra)
├── internal/
│   ├── types/                     # Gemeinsame Typen (Severity, Dependency, Report)
│   ├── parser/
│   │   ├── parser.go              # Maven pom.xml + Gradle build.gradle Parser
│   │   └── parser_test.go         # Parser-Tests
│   ├── resolver/
│   │   ├── graph.go               # Abhängigkeitsgraph-Auflöser + transitive Simulation
│   │   └── graph_test.go          # Resolver-Tests
│   ├── analyzer/
│   │   ├── build.go               # Build-Log-Analyser (Java/Gradle Fehlermuster)
│   │   └── security.go            # OSV.dev SCA-Connector
│   ├── engine/
│   │   └── engine.go              # Orchestrierer (Fortschritts-Callbacks, Phasen-Scan)
│   ├── report/
│   │   └── report.go              # Export: JSON, Markdown, SARIF
│   └── ui/
│       ├── banner.go              # Farbige CLI-Ausgabehilfen
│       └── tui.go                 # Interaktives Bubbletea TUI
└── AGENT_PROMPT.md                # KI-Analyse-Agent System-Prompt
```

---

## KI-Analyse-Agent

PixelGuard enthält einen System-Prompt für eine KI-Analyse-Engine in [`AGENT_PROMPT.md`](AGENT_PROMPT.md). Er ist darauf ausgelegt, Scanner-Ausgaben zu kohärenten, priorisierten Berichten zusammenzufassen — ohne CVEs oder Versionsnummern aus dem Gedächtnis zu erfinden

Verbinde ihn mit jedem tool-calling-fähigen LLM (Claude, GPT, etc.) indem du die strukturierte Scan-Ausgabe als Kontext einspeist.

---

## Beispiel-Terminal-Sitzung

```
$ pixelguard scan ./mein-java-projekt

 Projektstruktur wird gescannt...
✔ Erkannt: Maven-Projekt (Java 17, pom.xml gefunden)
✔ 187 Quelldateien, 42 deklarierte Abhängigkeiten, 118 aufgelöst (inkl. transitive)

Build:        3 Fehler, 12 Warnungen
Konflikte:    5 Abhängigkeitsversionskonflikte
Sicherheit:   7 Fundstellen — 2 Kritisch, 3 Hoch, 2 Mittel
```

---

## Beiträge

Beiträge willkommen! Öffne ein Issue oder sende einen PR.

## Lizenz

MIT Lizenz — siehe [LICENSE](LICENSE) für Details.

---

<p align="center">
  Mit ❤️ gemacht von <a href="https://github.com/jaylukdevelopment">Jay</a> & <a href="https://github.com/PixelCity-dev">PixelCity</a>
</p>

<p align="center">
  <a href="https://github.com/PixelCity-dev/pixelguard/stargazers">
    <img src="https://img.shields.io/github/stars/PixelCity-dev/pixelguard?style=social" alt="PixelGuard sternen">
  </a>
</p>
