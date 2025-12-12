# Product Requirements Document (PRD)

Project: p (pea) — Fast Prompt Storage & Retrieval CLI
Last updated: 2025-12-12T21:25:53.894Z

1. Purpose

- p is a fast, local CLI for storing short Markdown-like text under simple names and retrieving it instantly to stdout
  or the system clipboard. Optimized for daily repetitive use: speed, minimal keystrokes, completion, predictability.

2. Problem

- Developers frequently reuse structured text (prompts, notes, instructions, snippets). Existing tools are slow,
  UI-heavy, SaaS-bound, or lack fast CLI retrieval and completion.

3. Goals

- Primary: instant reuse, minimal typing, first-class shell completion, fast enough to feel invisible.
- Secondary: scriptable/pipe-friendly, plain files (Git-friendly), no daemons, trivial install via go install.

4. Non-Goals

- Not an LLM client or prompt executor; no cloud sync; no search in v0; no fuzzy picker/UI; no templates/variables in
  v0.

5. Target Users

- Software engineers and terminal power users who frequently use chat-based AI tools.

6. Core Use Cases

- UC1 Store by name (editor): `p add notes` → opens editor → saved as notes.txt
- UC2 Store by importing a file: `p add notes existing.txt` → copies file contents
- UC3 Retrieve for paste: `p notes` → prints and copies to clipboard if stdout is a TTY
- UC4 Scripted usage: `p notes > out.txt` → prints only; clipboard untouched
- UC5 Discover entries: `p ls`
- UC6 Fast completion: `p <TAB>`, `p add <TAB>`

7. Functional Requirements

- FR1 Storage
    - Entries are plain-text files: snake_case .txt (lowercase; strip invalid chars)
    - One entry per file; metadata stored inline as YAML front matter
    - Default directory: ~/.pea/prompts
    - Config file: ~/.pea/config.toml (may override defaults)
    - On first run, auto-create ~/.pea/prompts and config file if missing
- FR2 Add Entry
    - `p add <name>` opens $EDITOR <store>/<name>.txt; creates if missing
    - `p add <name> <file>` imports file; writes to <store>/<name>.txt
    - `echo "text" | p add <name>` reads stdin
    - Versioning is managed via Git commits; changes are recorded and latest (HEAD) is the default selection
- FR3 Retrieve Entry
    - `p <name>` reads <store>/<name>.txt (select latest version by default: HEAD in Git history)
    - Writes content to stdout
    - If stdout is a TTY, copy full output to clipboard; if redirected/piped, do not copy
- FR4 List Entries
    - `p ls` lists entry names (without .txt), one per line, sorted lexicographically
- FR5 Auto-Completion
    - `p completion zsh|bash` outputs completion script; `p completion install` installs into common locations
    - Completion based on <store>/*.txt (names without extension)
- FR6 Delete Entry
    - `p rm <name>` performs Git-backed delete (remove and commit); default policy: Git delete
- FR7 Rename Entry
    - `p mv <old> <new>` uses Git-backed rename; normal Git semantics apply

8. Non-Functional Requirements

- Performance: Go single binary; instant startup; lightweight config parsing via ~/.pea/config.toml; no background
  processes.
- UX: minimal flags; stable, predictable commands; completion-first discovery; use $EDITOR only (error if unset).
- Portability: macOS (v0); clipboard via library abstraction (golang-design/clipboard).

9. Error Handling

- Errors to stderr; non-zero exit on failure; clear, actionable messages; no partial output on error.

10. Installation

- `go install github.com/<user>/p@latest`
- Binary installs to $(go env GOPATH)/bin

11. Success Criteria

- Store and retrieve in <2 seconds; used multiple times per day; users rely on completion; repository contains simple,
  readable files; `--help` is sufficient for usage.

12. Summary

- p is a Unix-style primitive for text reuse: store text under a name, retrieve it instantly, with speed, simplicity,
  predictability, and excellent shell integration.
