# pea Roadmap

Status: v0.1.0 (macOS-focused, local-only). Updated: 2025-12-14T20:19:39.520Z.

## Now (hardening v0.1.x)

- Cross-platform clipboard abstraction (Linux/Windows) with graceful TTY detection and fallback.

## Next (v0.2)

- Metadata commands: add/edit/view description/tags stored as YAML front matter; list with optional metadata columns.
- Search and filtering: name/substring and tag-based queries over entries.
- Version awareness: allow retrieving a specific revision.
- Packaging/portability: publish binaries or install scripts for macOS/Linux; lighter dependency footprint.

## Later (exploratory)

- Remote/sync backends (e.g., git remotes or cloud object stores) with opt-in encryption at rest.
- Templates/snippets library with variables; bulk operations (export/import).
- Multi-store management with profiles and per-store configuration.
