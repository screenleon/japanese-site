# japanese-site

Japanese learning site covering grammar, vocabulary, kanji, quizzes, grading, corrective feedback, and a second track **国語教室** (school-style reading → evidence → expression → revision).

## Deployment modes

| Mode | Command / surface | What works |
|------|-------------------|------------|
| **Local API** (full product) | `make dev` or `make start` | Content browsing, quiz + grading, read tracking, progress, **国語 full cycle** |
| **Static cloud** (GitHub Pages / JS-018) | `make build-static` / site deploy | Content browsing + random grammar/vocab/kanji selection only. Quiz, sentence, and full 国語 progress are **off** (no backend / DB). |

Public URL: https://screenleon.github.io/japanese-site/ (enable **GitHub Pages → GitHub Actions** once if not already).

Contract details: `DECISIONS.md` → `2026-05-02: JS-018 GitHub Pages static deployment scope` and `docs/adr/0005-kokugo-track.md`.

---

## Local development

### Prerequisites

- Go (for `server/`)
- Node.js + npm (for `web/`)
- Optional: `jq` for static bake (`make build-static` / `make preview-static`)

### First-time setup

From the repository root:

```sh
make bootstrap
```

This installs Go modules and frontend deps, then runs the full seed pipeline (`make seed-all`, including external dataset download when needed). It can take a while the first time.

Equivalent pieces if you prefer to run them separately:

```sh
make bootstrap-server   # go mod download
make bootstrap-web      # npm install in web/
make seed-all           # build SQLite from L1 corpus + external data
```

### Daily development (recommended)

```sh
make dev
```

Starts both processes and stops them together with `Ctrl-C`:

| Service | URL | Notes |
|---------|-----|--------|
| **API** | http://localhost:8080 | Go server (`cmd/api`), SQLite progress |
| **Frontend** | http://127.0.0.1:**5173** | Vite hot-reload; `/api` is proxied to the backend |

Open the app at **http://127.0.0.1:5173** (not only `:8080`) while developing the React UI.

Override API listen address if needed:

```sh
LISTEN_ADDR=:9090 make dev
```

### One-port “production-style” local serve

Build the frontend bundle and serve SPA + API from the Go process on a single port:

```sh
make start
```

→ http://localhost:8080
(`STATIC_DIR=../web/dist`, same process as `./bin/api`)

### Backend-only / frontend-only

```sh
make run          # or: cd server && go run ./cmd/api
make web-dev      # or: cd web && npm run dev
```

### After content changes (corpus only)

Reload curated L1 corpus into the local DB without re-downloading JMdict/Tatoeba/etc.:

```sh
make db-update    # seed-corpus only
# or
make seed-corpus
```

Full re-seed (external downloads + corpus):

```sh
make seed-all
```

### Static mode preview (no Go / no DB)

Useful for checking the GitHub Pages–style build (content browse only):

```sh
make preview-static
```

→ http://localhost:4173 (override with `PREVIEW_PORT=5000 make preview-static`)

Requires `jq`. On Debian/Ubuntu: `apt install jq`; on macOS: `brew install jq`.

Bake + production static build (Pages base path `/japanese-site/`):

```sh
make build-static
```

Output: `web/dist/`. Serve with any static file server if you are not using `preview-static`.

### Useful make targets

```sh
make help              # short target list
make test              # Go + script tests
make lint              # rules + vocab + kokugo lint
make lint-kokugo       # School Kokugo unit JSON (ADR-0005)
make corpus-scale      # corpus floor report
make clean             # remove build artifacts and dev SQLite
```

---

## Project docs

| Doc | Role |
|-----|------|
| `ROADMAP.md` | Milestone narrative (incl. 国語 Phase 2) |
| `BACKLOG.md` / `project/backlog.yml` | Ticket queue |
| `DECISIONS.md` | Active architectural decisions |
| `docs/adr/0005-kokugo-track.md` | 国語教室 track design |
| `ARCHITECTURE.md` | System shape |
| `project/project-manifest.md` | Repo constraints and validation |
| `AGENTS.md` | Agent operating playbook |

**JS-134 shipped** (curated classmates after response + draft/改稿 side-by-side compare). Next 国語 work: **JS-135** unit pack 2, **JS-136** skill map. M4 LLM connector remains deferred until re-authorized.
