# Copilot CLI – Working Instructions

Last updated: 2025-12-12T21:22:54.609Z

## Purpose

I will operate in this repository like a disciplined engineer, making small, correct, reviewable changes.

## Non‑Negotiable Rules

- I work in small, logical units; never batch unrelated changes.
- I keep the repo in a runnable state after every change.
- I prefer clarity and correctness over cleverness.
- I do not pause to ask questions when the PRD/SRS/WORK_PROMPT provide the answer; proceed decisively.
- I never invent features or behaviors not explicitly requested, but I will implement all features specified end-to-end without waiting for approval.
- **Correction Protocol:** If I discover an error immediately after committing (and before pushing/sharing), I will fix
  it and **amend** the previous commit rather than creating a new "fix" commit, to keep the history clean.

## Workflow Loop (Automated Continuous Mode)

Follow this loop for every change, automatically and continuously:

1. **Inspect** — Read PRD/SRS/WORK_PROMPT and current code.
2. **Plan** — Pick the next smallest complete unit (a command with tests).
3. **Implement** — Code only that unit.
4. **Verify** — Run `just build`/`go build` and `just test`/`go test`; fix failures.
5. **Commit** — Conventional Commit message; commit immediately after passing.
6. **Review** — Self-check the diff; ensure only intended files changed.
7. **Continue** — Immediately proceed to the next unit with no pauses or approvals until the app is complete.

What Makes a Good Unit:

- A complete command with its tests (unit and/or e2e) and docs when applicable
- Brings value: adds something useful or fixes something broken
- Tells a story: reader understands what changed and why
- Stands alone: makes sense without other commits
- Is reviewable in under 5 minutes

**The Story Test:**
Ask yourself: "If someone reads just this commit message and diff in 6 months, will they understand what was
accomplished and why?" If no, the unit isn't complete.

### Small, Focused Changes

- Each commit addresses **one logical unit of work**
- Never batch unrelated changes (e.g., don't mix refactoring with new features)
- Break large features into a sequence of small, valuable commits
- Each step should leave the project better than before

### Always Runnable

- The repository must compile and run after every single commit
- No "broken" intermediate states
- If tests exist, they must pass
- Think: "Could we deploy after this commit?" (even if we won't)

### Clarity Over Cleverness

- Prefer explicit, readable code over "smart" solutions
- Simple solutions beat complex ones
- Code should explain its intent clearly

### Ask, Don't Guess

- When requirements are ambiguous, stop and ask for clarification
- Never invent features or behaviors not explicitly requested
- Confirm architectural decisions before implementing

## Commit Standards

### Format: Conventional Commits (Angular Style)

```
<type>(<scope>): <subject>

<body>

<footer>
```

**Types:**

- `feat:` — New feature or capability
- `fix:` — Bug fix
- `docs:` — Documentation only
- `refactor:` — Code restructuring without behavior change
- `test:` — Adding or updating tests
- `chore:` — Maintenance (deps, tooling, config)
- `perf:` — Performance improvement
- `ci:` — CI/CD configuration

**Examples:**

- `feat: add command skeleton for 'pea add'`
- `feat(config): implement TOML config file parsing`
- `fix: handle missing file gracefully with error message`
- `refactor: extract validation logic into separate function`
- `test: add integration tests for list command`
- `docs: update README with installation instructions`

**Guidelines:**

- Subject line: present tense, imperative mood, lowercase, no period
- Keep subject under 72 characters
- Body (optional): explain why and what, not how
- One logical concept per commit
- Avoid vague messages like "wip", "update", "fix stuff"

## Project Context (pea)

- Language: Go (module: pea, go 1.25).
- Current files: go.mod, main.go.
- Status: main.go prints a welcome message and demonstrates a loop.

## Technology Stack and Preferences

- CLI framework: use Cobra; where scaffolding/generation is applicable, use cobra-cli.
- Testing: prefer testscript for end-to-end CLI behavior; use testing for unit logic.
- Dependencies: use the standard library first; pick popular, well-maintained libraries only when necessary.
- Always ensure the latest stable versions of libraries/tools are used; run `go mod tidy` and update modules as needed on each feature step.

## Recommended Implementation Order

1. **Bootstrap:** Ensure minimal runnable binary; initialize CLI framework.
2. **Core Helpers:** Establish error handling patterns and config reading.
3. **Primary Commands:** Implement command structure (e.g., `add`, `retrieve`, `list`) iteratively.
4. **UX Polish:** Refine error messages, exit codes, and output formatting.
5. **Completion:** Generate shell completions and installation docs.

## Coding Style

- Write idiomatic Go: simple, explicit code; short functions.
- Handle errors explicitly and fail fast with helpful messages.
- Avoid premature abstractions; copy-paste is better than a wrong abstraction.

## Verification

- Every commit must compile (go build ./...) and pass tests (go test ./...).
- Test strategy:
    - Unit tests: critical internal logic/algorithms.
    - Integration tests: use testscript (or similar) to verify end-to-end CLI behavior (flags, exit codes, output).
- Execution: go run ./... (or specific command).
- Automation: if a Taskfile exists, use task-defined tasks (e.g., task build, task test) for consistency.

## When Unsure

- Prefer action: if PRD/SRS specify behavior, implement without asking.
- Only ask if requirements are truly ambiguous or conflicting; otherwise proceed.

## Process Checklist

1. Define tiny work units; one goal per commit
2. Plan the step; write intent in commit body
3. Always build and run before/after changes
4. Each commit must compile and pass tests
5. Prioritize integration tests for CLI surface area
6. Use Conventional Commits (Angular prefixes)
7. Choose widely adopted, best-in-class tools/libs
8. Keep solutions simple; avoid early abstractions
9. Use repeatable tasks (task/Makefile) for build/test/run
10. Document non-obvious decisions briefly in code
11. Timebox exploration; commit learnings separately
12. Stop and ask when scope/requirements are unclear
13. Self-review diffs; ensure only intended files change
14. Maintain clear error handling and exit codes
15. Optimize after correctness and clarity are achieved
16. Amend commits if errors are found immediately after saving