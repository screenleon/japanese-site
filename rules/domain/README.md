# Domain Rules — japanese-site

Domain rules apply per the `Workspace boundaries` table in
`project/project-manifest.md`. Each file follows the standard rule entry
format from `agent-playbook-template/rules/domain/backend-api.md`.

| File | Domain | Applies to |
|---|---|---|
| `content-source.md` | content-source | every persisted learning-content row |
| `jlpt-content-accuracy.md` | jlpt-content-accuracy | rows with a JLPT level + question generation |
| `grading-feedback.md` | grading-feedback | grading endpoint responses |
| `connector-credential.md` | connector-credential | server/connector credential handling (slated for extraction at M4) |
| `corpus-storage.md` | corpus-storage | three-tier storage (L1 curated / L2 LLM cache / L3 external) and cache promotion |
| `kokugo-content-authoring.md` | kokugo-content-authoring | Kokugo L1 unit provenance, stage fit, native-perspective review, and pre-merge checklist |
| `backend-api.md` | backend-api | public HTTP handlers under `/api/*` (envelope, evolution, body limits) |
| `frontend-components.md` | frontend-components | React components under `web/src/` (state location, async UI states, API coupling) |
