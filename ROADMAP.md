# pea Roadmap

Status: v0.1.0 (macOS-focused, local-only). Updated: 2025-12-14T16:48:37Z.

## Now (hardening v0.1.x)

- Cross-platform clipboard abstraction (Linux/Windows) with graceful TTY detection and fallback.
- Tighten Git-backed versioning:
  - Deterministic auto-commit subjects for add/mv/rm (including chore vs refactor rename paths) and consistent add/rm semantics.
  - Surface Git errors to stderr without failing the primary command; clear guidance when Git is absent/misconfigured.
  - E2E coverage for add/mv/rm auto-commit flows and failure cases.
- UX polish: clearer errors for missing entries, invalid names, and editor launch failures; improve completion install
  paths.

## Next (v0.2)

- Metadata commands: add/edit/view description/tags stored as YAML front matter; list with optional metadata columns.
- Search and filtering: name/substring and tag-based queries over entries.
- Version awareness: show history for an entry and allow retrieving a specific revision.
- Packaging/portability: publish binaries or install scripts for macOS/Linux; lighter dependency footprint.
- Safety: optional confirmation/undo for destructive commands; dry-run flags where applicable.

## Later (exploratory)

- Remote/sync backends (e.g., git remotes or cloud object stores) with opt-in encryption at rest.
- Templates/snippets library with variables; bulk operations (export/import).
- Multi-store management with profiles and per-store configuration.
