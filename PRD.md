Project: pea — Fast Prompt Storage & Retrieval CLI

Summary

pea is a small, fast CLI to store short Markdown-like text under simple names and retrieve it instantly. It focuses on minimal UX, shell completion, and plain files that work well with Git.

Goals

- Minimal interface for speedy storage & retrieval
- Portable plain-file storage that works with Git
- Fast retrieval for copy/paste usage

User stories

- Store by name (editor): `pea add notes` → opens editor → saved as notes.md
- Store by importing a file: `pea add notes existing.md` → copies file contents
- Retrieve for paste: `pea notes` → prints and copies to clipboard if stdout is a TTY
- Scripted usage: `pea notes > out.txt` → prints only; clipboard untouched
- Discover entries: `pea ls`
- Fast completion: `pea <TAB>`, `pea add <TAB>`

Commands

- `pea add <name>` opens $EDITOR <store>/<name>.md; creates if missing (legacy `.txt` still readable)
- `pea add <name> <file>` imports file; writes to <store>/<name>.md by default
- `echo "text" | pea add <name>` reads stdin
- `pea <name>` reads <store>/<name>.md (select latest version by default: HEAD in Git history; legacy `.txt` also supported)
- `pea ls` lists entry names (without extension), one per line, sorted lexicographically
- `pea completion zsh|bash` outputs completion script; `pea completion install` installs into common locations
- `pea rm <name>` performs Git-backed delete (remove and commit); default policy: Git delete
- `pea mv <old> <new>` uses Git-backed rename; normal Git semantics apply

Install

- `go install github.com/<user>/pea@latest`

Notes

- pea is a Unix-style primitive for text reuse: store text under a name, retrieve it instantly, with speed, simplicity, and predictable semantics.

Planned (from roadmap)
- Cross-platform clipboard: use platform-specific clipboard backends (Linux/Windows) with graceful fallback and TTY-aware behavior; keep macOS parity.
- Metadata: commands to add/edit/view description and tags (YAML front matter); listing can optionally include metadata columns.
- Search/filter: name substring and tag-based queries over stored entries; results respect current store selection.
- Version awareness: allow retrieval at specific git refs via `--rev`; history is already available; produce clear errors for missing ref or file.
- Packaging/portability: publish binaries or install scripts for macOS/Linux; keep dependency footprint small; document requirements (git, clipboard tools).
- Multi-store profiles: central config maps store names to absolute paths; selection precedence `--store` > `PEA_STORE` > default; each store remains an isolated git repo.
