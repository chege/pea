# 🫛 pea

> **Fast, local prompt storage & retrieval.**

[![Go Version](https://img.shields.io/badge/go-1.25-blue.svg)](https://golang.org/dl/)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)
[![Build Status](https://img.shields.io/badge/build-passing-brightgreen.svg)]()

**pea** is a minimalist CLI designed to store short text snippets (prompts, notes, config fragments) under simple names and retrieve them instantly. It combines the speed of local files with the power of Git versioning and shell integration.

---

## ✨ Features

- **🚀 Instant Retrieval:** `pea get <name>` prints content to stdout.
- **📋 Smart Clipboard:** Auto-copies content to clipboard when running in a terminal (TTY).
- **🔎 Powerful Search:** Find entries by name, content, or tags.
- **📂 Plain Files:** Stores everything as simple `.md` files in `~/.pea/prompts`.
- **🐚 Shell Completion:** First-class autocomplete for commands *and* your snippet names.
- **🔄 Git-Backed:** Every change (add, edit, remove) is automatically committed to Git.
- **🏷️ Metadata:** Supports YAML front matter (stripped on retrieval).

## 🚀 Quick Start

```bash
# 1. Install (requires Go 1.25+)
go install github.com/user/pea@latest

# 2. Add your first snippet
echo "Hello, world!" | pea add hello

# 3. Retrieve it (copies to clipboard automatically!)
pea get hello

# 4. List all snippets
pea ls
```

## 🛠️ Usage

### Core Commands

| Command | Usage | Description |
| :--- | :--- | :--- |
| **Retrieve** | `pea get <name>` | Print content (and copy to clipboard). |
| **Copy** | `pea cp <name>` | Copy content to clipboard (no print). |
| **Add** | `pea add <name>` | Create/Edit via `$EDITOR` or stdin. |
| **List** | `pea ls` | List all entry names. |
| **Search** | `pea search <query>` | Search by name, content, or tags. |
| **Remove** | `pea rm <name>` | Delete an entry (versioned). |
| **Move** | `pea mv <old> <new>` | Rename an entry (versioned). |
| **History** | `pea history <name>` | View Git history of an entry. |
| **Remote** | `pea remote <url>` | Configure remote git sync. |
| **Create Repo** | `pea remote create <name>` | Create & sync with a new GitHub repo. |
| **Sync** | `pea sync` | Manual git pull (rebase) & push. |

### Adding Content

**From Stdin:**
```bash
echo "sk-12345-api-key" | pea add api_key
```

**From File:**
```bash
pea add project_docs ./README.md
```

**Using Editor:**
```bash
# Opens $EDITOR (vim, nano, code, etc.)
pea add my_notes
```

### Search & Tags

Entries can contain YAML front matter for organization:

```yaml
---
tags: [work, email]
---
Here is my email template...
```

**Search:**
```bash
pea search template          # Search by name/content
pea search --tag work        # Filter by tag
```

## ⚙️ Configuration

`pea` works out of the box with zero config. By default, it stores data in `~/.pea/prompts`.

**Environment Variables:**
*   `PEA_STORE`: Override the storage directory path.

**Config File:**
Located at `~/.pea/config.toml`:
```toml
store_dir = "/path/to/my/custom/store"
```

## 🔮 Shell Completion

Get super-fast autocomplete for both commands and your stored snippet names.

**Install (Recommended):**
```bash
# Detects your shell (zsh/bash), installs script, and updates your config (~/.zshrc or ~/.bashrc)
pea completion install
```

**Manual Generation:**
```bash
pea completion bash > ~/.bash_completion
pea completion zsh > ~/.zshrc
```

## 🧑‍💻 Development

Requirements: **Go 1.25+**, **Make** (or `just`).

```bash
# Build binary
just build

# Run all tests (Unit + E2E)
just check

# Run specific E2E tests
go test -v ./e2e
```

## 📄 License

MIT