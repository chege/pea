# 🫛 pea — Fast Prompt Storage & Retrieval CLI

pea is a fast, local CLI to store short Markdown-like text under simple names and retrieve it instantly to stdout (and to the clipboard when in a TTY). It favors speed, minimal keystrokes, shell completion, and plain files that work well with Git.

Status: v0.1.0 (macOS focus). Core commands implemented: add, retrieve, ls, rm, mv, completion; env/config; TTY-only clipboard; simple YAML front matter support; basic Git-backed versioning.

## ✨ Features
- Plain files: one entry per file, snake_case .md (lowercase; invalid chars stripped)
- Instant retrieval: `pea <name>` prints content; when stdout is a TTY, content is copied to clipboard
- Add content three ways: stdin, import a file, or open $EDITOR
- Discover: `pea ls` lists names (sorted)
- Manage: `pea rm <name>` and `pea mv <old> <new>`
- Search: `pea search <query>` finds entries by name, content, or tags
- Shell completion: `pea completion bash|zsh` and `pea completion install`
- Configurable storage dir via env or TOML config
- Git-backed changes: store is a Git repo; add/rm/mv create commits (best-effort)
- Front matter: YAML front matter delimited by `---` is stripped on retrieval

## ⚙️ Install
- Go 1.25+
- Build locally:
  - `go build -o bin/pea .` or `just build`
- Optional: `go install ./...` to place on PATH

## 🚀 Quick start
```
# Add via stdin
echo "Hello world" | pea add hello_world

# Retrieve (prints and copies to clipboard when TTY)
pea hello_world

# List entries
pea ls

# Rename and delete
pea mv hello_world hello
pea rm hello
```

## 🛠️ Commands
- `pea [name]` — retrieve entry content (stdout; copies to clipboard in TTY)
- `pea add <name> [file]` — add entry by name
  - stdin: `echo text | pea add notes`
  - import file: `pea add notes readme.md`
  - editor: if no stdin/file, opens `$EDITOR` (set it, e.g., `export EDITOR=vim`); if unset, falls back to the OS default handler via `github.com/pkg/browser`
- `pea ls` — list entry names (sorted)
- `pea rm <name>` — delete entry (Git commit best-effort)
- `pea mv <old> <new>` — rename entry (Git commit best-effort)
- `pea search <query>` — search entries by name substring, content, or tags
  - filter by tag: `pea search --tag coding`
- `pea completion [bash|zsh]` — output completion script (redirect to your shell’s completion directory)

Run `pea --help` or `pea <command> --help` for usage.

## 📁 Storage layout
- Default store: `~/.pea/prompts`
- Files: `<name>.md` (snake_case; legacy `.txt` still readable)
- Content may include an optional YAML front matter block at the top:
  ```
  ---
  description: My snippet
  tags: [ai, daily]
  ---
  Actual content starts here
  ```
- Retrieval strips the front matter and prints only the body

## ⚙️ Configuration
- Environment:
  - `PEA_STORE` — absolute path to storage directory; overrides config file
- Config file: `~/.pea/config.toml`
  - Example:
    ```toml
    store_dir = "/absolute/path/to/store"
    ```
  - On first run, `~/.pea/config.toml` is created if missing

## 📋 Clipboard behavior (macOS v0)
- If stdout is a TTY, `pea <name>` copies the printed content to the system clipboard (uses `pbcopy`)
- If output is redirected or piped, clipboard is not touched

## 🔁 Shell completion
- Generate scripts:
  - Bash: `pea completion bash > ~/.pea/pea.bash`
  - Zsh: `pea completion zsh > ~/.pea/_pea`
- Automatic install to `~/.pea/`: `pea completion install`
- Add to your shell profile to source these files as desired

## 🔀 Git-backed versioning
- The store directory is initialized as a Git repo on first use
- `add`, `rm`, and `mv` attempt to create commits for changes (best-effort; does not fail the command on Git errors)

## 🧑‍💻 Development
- Requirements: Go 1.25+, macOS (for v0 clipboard behavior)
- Common tasks (requires `just`):
  - `just build` — build local binary to `bin/pea`
  - `just test` — run tests
  - `just check` — fmt, vet, and test
  - `just tidy` — `go mod tidy`
- Without `just`:
  - `go build -o bin/pea .`
  - `go test ./...`

## ✅ Testing
- End-to-end tests live under `./e2e` and exercise CLI behavior (add/retrieve/list/delete/rename, completion, config, clipboard)
- Run: `go test ./e2e` or `just test`

## 🗺️ Roadmap
- Cross-platform clipboard abstraction (Linux/Windows)
- Richer metadata handling and commands
- Additional safety and UX polish

## 📝 Notes
- Latest stable library versions are used (e.g., cobra/pflag/toml); run `go mod tidy` as needed
- v0 targets macOS and a simple, predictable CLI UX
