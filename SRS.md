# pea CLI — Software Requirements Specification (SRS)

Last updated: 2025-12-13T07:35:00.694Z
Status: Draft — decisions incorporated

1. Introduction
- Purpose: define functional and non-functional requirements for the pea CLI.
- Scope: single-binary CLI with core commands (add, retrieve, list, delete, rename), platform behavior (clipboard, TTY), clear UX (errors, exit codes), and shell completion.
- Definitions: CLI (command-line interface), TTY (interactive terminal), E2E (end-to-end tests).

2. Overall Description
- Product perspective: standalone Go CLI (module: pea) that stores prompts as local files and interacts with env/config and clipboard.
- User classes: LLM users (CLI users).
- Operating environment: Go 1.25; macOS; shells: zsh and bash.
- Constraints: idiomatic Go; popular, well-maintained libraries/tools; keep dependencies updated to latest stable; Cobra for CLI; commits compile and tests pass; emphasize e2e tests; no telemetry.
- Assumptions/dependencies: local filesystem; standard shell; platform clipboard available.

3. Functional Requirements
- FR-1 Add Prompt
  - Inputs: name; content via stdin, file path, or $EDITOR (flags to select); optional metadata description (string), tags (list).
  - Behavior: create file in ~/.pea/prompts; filename snake_case with .md (lowercase; strip invalid chars); metadata stored inline as YAML front matter; versioning managed via Git commits; if $EDITOR is unset and no stdin/file is provided, open the target file with the system default handler via `github.com/pkg/browser` (error if launch fails); legacy `.txt` entries remain readable.
  - Outputs: confirmation with location; changes recorded in Git.
  - Errors: validation, conflicts, IO failures.
- FR-2 Retrieve Prompt
  - Behavior: print prompt content; when TTY, copy full prompt to clipboard; select latest by default (HEAD in Git).
  - Outputs: content (and optionally metadata).
  - Errors: not found, ambiguous match.
- FR-3 List Prompts
  - Behavior: list prompt names (without extension) in lexicographic order; output is names only in v0 (no metadata or counts).
  - Outputs: human-readable list of names, one per line.
  - Errors: IO failures.
- FR-4 Delete Prompt
  - Behavior: Git-backed delete (remove and commit); default policy uses Git delete semantics.
- FR-5 Rename Prompt
  - Behavior: Git-backed rename; normal Git semantics apply.
- FR-6 Shell Completion
  - Behavior: generate completion scripts for zsh and bash; provide `p completion install` to install into common locations.
- FR-7 Configuration Handling
  - Behavior: load config from TOML at ~/.pea/config.toml; env vars override config file; on first run, auto-create ~/.pea/prompts and config file if missing.
- FR-8 Clipboard Handling
  - Behavior: TTY-only clipboard copy on retrieve; use library abstraction (golang-design/clipboard).
- FR-9 TTY Detection
  - Behavior: adapt for TTY vs non-TTY.
- FR-10 Version and Help
  - Behavior: provide --version (SemVer) and --help; show usage and subcommands.
- FR-11 Output Format
  - Behavior: human-readable output only (no JSON now).

4. Non-Functional Requirements
- NFR-1 Usability: clear errors, consistent exit codes, readable defaults.
- NFR-2 Reliability: safe writes; avoid data corruption; handle partial failures.
- NFR-3 Performance: fast CLI UX; small files; no special constraints.
- NFR-4 Portability: works on macOS.
- NFR-5 Testability: idiomatic e2e tests for all commands; unit tests for core logic; go test ./... must pass; primary platform macOS.
- NFR-6 Maintainability: simple, readable code; minimal abstractions; standard tooling.
- NFR-7 Security/Privacy: no telemetry; respect file permissions; avoid leaking sensitive data (including to clipboard inadvertently).

5. External Interfaces
- Filesystem: read/write within ~/.pea/prompts.
- Config: TOML at ~/.pea/config.toml; env vars may override.
- Clipboard: platform-specific (TTY-only copy on retrieve).
- STDIN/STDOUT/STDERR: standard IO for scripting.

6. Data Model
- Prompt: plain-text .md file (legacy `.txt` supported); one per prompt.
- Naming: snake_case from prompt name.
- Metadata: description (string), tags (list); representation decided during design (inline header or sidecar index).
- Versioning: multiple versions per prompt; default operations act on latest; scheme defined in design.

7. Error Handling and Exit Codes
- Exit 0 on success; non-zero on error with helpful message; consistent mapping for common errors (not found, invalid input, IO).

8. Logging and Observability
- Minimal default logging; optional --verbose flag.

9. Documentation and Completion
- Built-in help for commands and flags; completion scripts for zsh and bash with install instructions.

10. Testing and Verification
- E2E tests for add/retrieve/list/delete/rename happy paths and edge cases; unit tests for parsing/validation/file IO; use idiomatic integration test tooling for the language (e.g., Go test/e2e frameworks); CI runs build and test across all packages.

11. Risks and Open Issues
- Versioning scheme specifics; metadata representation; UX for duplicate/version selection.

12. Acceptance Criteria
- Commands behave as specified; human-readable output; clipboard copy on retrieve in TTY; completion generated and documented; all tests pass; repository builds cleanly.

13. Future Enhancements (non-goals now)
- Tags-based search and filtering; remote sync/backends; templating/rendering; richer query/filter language.

14. Roadmap-aligned design notes (planned)
- Cross-platform clipboard: abstract ClipboardImpl for Linux/Windows; detect platform and use xclip/wl-copy equivalents with graceful fallback and TTY-aware guard rails.
- Metadata commands: add/edit/view description/tags persisted as YAML front matter; retrieval strips front matter; list supports optional metadata columns.
- Search and filtering: name substring and tag-based queries over stored entries; respect same store selection/precedence as other commands.
- Version awareness: retrieving specific revisions via `--rev <ref>` using `git show`; history already covered; clear errors when ref/file missing.
- Packaging/portability: publish install scripts/binaries for macOS/Linux; minimize dependencies; document required tools (git/clipboard).
- Multi-store profiles: central config maps store names to absolute paths; selection precedence `--store` > `PEA_STORE` > default; each store self-contained with its own .git.
