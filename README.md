# p (pea) — Fast Prompt Storage & Retrieval CLI

p is a fast, local CLI to store short Markdown-like text under simple names and retrieve it instantly to stdout (and to the clipboard when in a TTY). It favors speed, minimal keystrokes, shell completion, and plain files that work well with Git.

Status: v0.1.0 (macOS focus). Core commands implemented: add, retrieve, ls, rm, mv, completion; env/config; TTY-only clipboard; simple YAML front matter support; basic Git-backed versioning.

## Features
- Plain files: one entry per file, snake_case .txt (lowercase; invalid chars stripped)
- Instant retrieval: `p <name>` prints content; when stdout is a TTY, content is copied to clipboard
- Add content three ways: stdin, import a file, or open $EDITOR
- Discover: `p ls` lists names (sorted)
- Manage: `p rm <name>` and `p mv <old> <new>`
- Shell completion: `p completion bash|zsh` and `p completion install`
- Configurable storage dir via env or TOML config
- Git-backed changes: store is a Git repo; add/rm/mv create commits (best-effort)
- Front matter: YAML front matter delimited by `---` is stripped on retrieval

## Install
- Go 1.25+
- Build locally:
  - `go build -o bin/pea .` or `just build`
- Optional: `go install ./...` to place on PATH

## Quick start
```
# Add via stdin
echo "Hello world" | p add hello_world

# Retrieve (prints and copies to clipboard when TTY)
p hello_world

# List entries
p ls

# Rename and delete
p mv hello_world hello
p rm hello
```

## Commands
- `p [name]` — retrieve entry content (stdout; copies to clipboard in TTY)
- `p add <name> [file]` — add entry by name
  - stdin: `echo text | p add notes`
  - import file: `p add notes readme.txt`
  - editor: if no stdin/file, opens `$EDITOR` (set it, e.g., `export EDITOR=vim`); if unset, falls back to the OS default handler via `github.com/pkg/browser`
- `p ls` — list entry names (sorted)
- `p rm <name>` — delete entry (Git commit best-effort)
- `p mv <old> <new>` — rename entry (Git commit best-effort)
- `p completion [bash|zsh]` — output completion script (redirect to your shell’s completion directory)

Run `p --help` or `p <command> --help` for usage.

## Storage layout
- Default store: `~/.pea/prompts`
- Files: `<name>.txt` (snake_case)
- Content may include an optional YAML front matter block at the top:
  ```
  ---
  description: My snippet
  tags: [ai, daily]
  ---
  Actual content starts here
  ```
- Retrieval strips the front matter and prints only the body

## Configuration
- Environment:
  - `PEA_STORE` — absolute path to storage directory; overrides config file
- Config file: `~/.pea/config.toml`
  - Example:
    ```toml
    store_dir = "/absolute/path/to/store"
    ```
  - On first run, `~/.pea/config.toml` is created if missing

## Clipboard behavior (macOS v0)
- If stdout is a TTY, `p <name>` copies the printed content to the system clipboard (uses `pbcopy`)
- If output is redirected or piped, clipboard is not touched

## Shell completion
- Generate scripts:
  - Bash: `p completion bash > ~/.pea/p.bash`
  - Zsh: `p completion zsh > ~/.pea/_p`
- Automatic install to `~/.pea/`: `p completion install`
- Add to your shell profile to source these files as desired

## Git-backed versioning
- The store directory is initialized as a Git repo on first use
- `add`, `rm`, and `mv` attempt to create commits for changes (best-effort; does not fail the command on Git errors)

## Development
- Requirements: Go 1.25+, macOS (for v0 clipboard behavior)
- Common tasks (requires `just`):
  - `just build` — build local binary to `bin/pea`
  - `just test` — run tests
  - `just check` — fmt, vet, and test
  - `just tidy` — `go mod tidy`
- Without `just`:
  - `go build -o bin/pea .`
  - `go test ./...`

## Testing
- End-to-end tests live under `./e2e` and exercise CLI behavior (add/retrieve/list/delete/rename, completion, config, clipboard)
- Run: `go test ./e2e` or `just test`

## Roadmap
- Cross-platform clipboard abstraction (Linux/Windows)
- Richer metadata handling and commands
- Additional safety and UX polish

## Notes
- Latest stable library versions are used (e.g., cobra/pflag/toml); run `go mod tidy` as needed
- v0 targets macOS and a simple, predictable CLI UX
