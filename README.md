# japanese-site

Japanese learning site covering grammar, vocabulary, kanji, quizzes, grading, and corrective feedback.

## Deployment

Public URL: https://screenleon.github.io/japanese-site/ (TBD until the first GitHub Pages deploy succeeds).

One-time setup: in the GitHub repository settings, enable GitHub Pages and set the source to **GitHub Actions**.

This project has two deployment modes. Local mode runs the Go API with the full feature set: content browsing, quizzes, grading, read tracking, and progress. Static cloud mode is the GitHub Pages build for JS-018: it bakes `server/data/corpus/**` into static files and serves content browsing plus random grammar/vocab/kanji selection only. Quiz and sentence tabs are hidden in static mode because there is no backend, database, grading state, or sentence bake. See `DECISIONS.md` -> `2026-05-02: JS-018 GitHub Pages static deployment scope` for the full contract.

Local development:

```sh
make seed-all
make dev
```

For backend-only local API work:

```sh
cd server && go run ./cmd/api
```

Static build:

```sh
make build-static
```

The static build writes `web/dist/`. Serve that directory with any static file server to preview the GitHub Pages build locally.
