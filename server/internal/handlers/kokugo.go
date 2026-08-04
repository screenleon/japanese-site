package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"unicode"

	"github.com/screenleon/japanese-site/server/internal/kokugo"
	"github.com/screenleon/japanese-site/server/internal/store"
)

const maxKokugoBodyBytes = 64 << 10 // 64 KB — artifact drafts + answers

// errKokugoBodyTooLarge is returned when a request body exceeds maxKokugoBodyBytes.
var errKokugoBodyTooLarge = errors.New("kokugo body too large")

type kokugoHandlers struct {
	db     *store.DB
	loader kokugo.Loader
	// enabled when progress store is sqlite-backed (local API full cycle).
	progressEnabled bool
}

func registerKokugo(mux *http.ServeMux, db *store.DB, loader kokugo.Loader, progressEnabled bool) {
	h := &kokugoHandlers{db: db, loader: loader, progressEnabled: progressEnabled}
	mux.HandleFunc("GET /api/kokugo/units", h.listUnits)
	mux.HandleFunc("GET /api/kokugo/units/{stage}/{id}", h.getUnit)
	mux.HandleFunc("GET /api/kokugo/skills", h.getSkills)
	mux.HandleFunc("GET /api/kokugo/progress", h.listProgress)
	mux.HandleFunc("GET /api/kokugo/progress/{stage}/{id}", h.getProgress)
	mux.HandleFunc("PUT /api/kokugo/progress/{stage}/{id}", h.putProgress)
	mux.HandleFunc("POST /api/kokugo/progress/{stage}/{id}/tasks/{taskId}", h.postTask)
	mux.HandleFunc("PUT /api/kokugo/progress/{stage}/{id}/artifact", h.putArtifact)
}

func (h *kokugoHandlers) listUnits(w http.ResponseWriter, r *http.Request) {
	if h.loader.Root == "" {
		writeJSON(w, http.StatusOK, map[string]any{"units": []any{}, "count": 0})
		return
	}
	units, err := h.loader.List()
	if err != nil {
		httpError(w, r, http.StatusInternalServerError, "internal", err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"units": units, "count": len(units)})
}

func (h *kokugoHandlers) getUnit(w http.ResponseWriter, r *http.Request) {
	stage, id := r.PathValue("stage"), r.PathValue("id")
	if !validKokugoSegment(stage) || !validKokugoSegment(id) {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid_id"})
		return
	}
	raw, err := h.loader.Get(stage, id)
	if errors.Is(err, kokugo.ErrUnitNotFound) {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "not_found"})
		return
	}
	if err != nil {
		httpError(w, r, http.StatusInternalServerError, "internal", err)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(raw)
}

// getSkills returns the Track B skill map + weak-skill review queue (JS-136).
// Available whenever the corpus loader is configured; attempt/artifact stats
// require progress store (otherwise skills are all unseen with review from units).
func (h *kokugoHandlers) getSkills(w http.ResponseWriter, r *http.Request) {
	if h.loader.Root == "" {
		writeJSON(w, http.StatusOK, kokugo.SkillMap{
			Skills:      []kokugo.SkillStat{},
			ReviewQueue: []kokugo.ReviewItem{},
		})
		return
	}
	units, err := h.loader.ListMaps()
	if err != nil {
		httpError(w, r, http.StatusInternalServerError, "internal", err)
		return
	}
	idx := kokugo.BuildUnitSkillIndex(units)

	var attempts []kokugo.AttemptRow
	var artifacts []kokugo.ArtifactRow
	var progress []kokugo.ProgressRow
	if h.progressEnabled && h.db != nil {
		rawAttempts, err := store.ListAllKokugoTaskAttempts(r.Context(), h.db)
		if err != nil {
			httpError(w, r, http.StatusInternalServerError, "internal", err)
			return
		}
		attempts = toSkillAttemptRows(rawAttempts)
		rawArts, err := store.ListAllKokugoArtifacts(r.Context(), h.db)
		if err != nil {
			httpError(w, r, http.StatusInternalServerError, "internal", err)
			return
		}
		artifacts = toSkillArtifactRows(rawArts)
		rawProg, err := store.ListKokugoProgress(r.Context(), h.db)
		if err != nil {
			httpError(w, r, http.StatusInternalServerError, "internal", err)
			return
		}
		progress = toSkillProgressRows(rawProg)
	}

	m := kokugo.AggregateSkillMap(idx, attempts, artifacts, progress)
	writeJSON(w, http.StatusOK, m)
}

func toSkillAttemptRows(in []store.KokugoTaskAttempt) []kokugo.AttemptRow {
	out := make([]kokugo.AttemptRow, 0, len(in))
	for _, a := range in {
		out = append(out, kokugo.AttemptRow{
			UnitKey:   a.UnitKey,
			TaskID:    a.TaskID,
			Correct:   a.Correct,
			CreatedAt: a.CreatedAt,
			ID:        a.ID,
		})
	}
	return out
}

func toSkillArtifactRows(in []store.KokugoArtifact) []kokugo.ArtifactRow {
	out := make([]kokugo.ArtifactRow, 0, len(in))
	for _, a := range in {
		out = append(out, kokugo.ArtifactRow{
			UnitKey:          a.UnitKey,
			Revision:         a.Revision,
			Body:             a.Body,
			ChecklistAllTrue: kokugo.ChecklistAllTrue(a.Checklist),
		})
	}
	return out
}

func toSkillProgressRows(in []store.KokugoUnitProgress) []kokugo.ProgressRow {
	out := make([]kokugo.ProgressRow, 0, len(in))
	for _, p := range in {
		out = append(out, kokugo.ProgressRow{
			UnitKey: p.UnitKey,
			Stage:   p.Stage,
			UnitID:  p.UnitID,
			Status:  p.Status,
		})
	}
	return out
}

func (h *kokugoHandlers) listProgress(w http.ResponseWriter, r *http.Request) {
	if !h.progressEnabled {
		writeJSON(w, http.StatusOK, map[string]any{"items": []any{}, "count": 0})
		return
	}
	items, err := store.ListKokugoProgress(r.Context(), h.db)
	if err != nil {
		httpError(w, r, http.StatusInternalServerError, "internal", err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": items, "count": len(items)})
}

func (h *kokugoHandlers) getProgress(w http.ResponseWriter, r *http.Request) {
	if !h.progressEnabled {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "progress_disabled"})
		return
	}
	stage, id := r.PathValue("stage"), r.PathValue("id")
	if !validKokugoSegment(stage) || !validKokugoSegment(id) {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid_id"})
		return
	}
	state, err := store.GetKokugoUnitState(r.Context(), h.db, stage, id)
	if err != nil {
		httpError(w, r, http.StatusInternalServerError, "internal", err)
		return
	}
	writeJSON(w, http.StatusOK, state)
}

func (h *kokugoHandlers) putProgress(w http.ResponseWriter, r *http.Request) {
	if !h.progressEnabled {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "progress_disabled"})
		return
	}
	stage, id := r.PathValue("stage"), r.PathValue("id")
	if !validKokugoSegment(stage) || !validKokugoSegment(id) {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid_id"})
		return
	}
	var body struct {
		Step   string `json:"step"`
		Status string `json:"status"`
	}
	if err := decodeKokugoBody(r, &body); err != nil {
		writeKokugoDecodeError(w, err)
		return
	}
	// Verify unit exists so we don't create progress for typos.
	unit, err := h.loader.GetMap(stage, id)
	if errors.Is(err, kokugo.ErrUnitNotFound) {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "not_found"})
		return
	} else if err != nil {
		httpError(w, r, http.StatusInternalServerError, "internal", err)
		return
	}

	// Completion is server-derived (ADR-0005 full cycle). Clients may request
	// completed only when attempts (+ artifacts when required) already exist.
	wantComplete := body.Status == "completed" || body.Step == "done"
	if wantComplete {
		ok, err := store.KokugoCycleComplete(r.Context(), h.db, stage, id, unit)
		if err != nil {
			httpError(w, r, http.StatusInternalServerError, "internal", err)
			return
		}
		if !ok {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "cycle_incomplete"})
			return
		}
		p, err := store.UpdateKokugoProgress(r.Context(), h.db, stage, id, "done", "completed")
		if err != nil {
			httpError(w, r, http.StatusInternalServerError, "internal", err)
			return
		}
		writeJSON(w, http.StatusOK, p)
		return
	}

	if body.Status == "" {
		body.Status = "in_progress"
	}
	if body.Status != "in_progress" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid_status"})
		return
	}
	if !validKokugoClientStep(body.Step) {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid_step"})
		return
	}
	p, err := store.UpdateKokugoProgress(r.Context(), h.db, stage, id, body.Step, body.Status)
	if err != nil {
		if errors.Is(err, store.ErrKokugoCompletedLocked) {
			writeJSON(w, http.StatusConflict, map[string]string{"error": "completed_locked"})
			return
		}
		if strings.Contains(err.Error(), "invalid status") {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid_status"})
			return
		}
		httpError(w, r, http.StatusInternalServerError, "internal", err)
		return
	}
	writeJSON(w, http.StatusOK, p)
}

func (h *kokugoHandlers) postTask(w http.ResponseWriter, r *http.Request) {
	if !h.progressEnabled {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "progress_disabled"})
		return
	}
	stage, id, taskID := r.PathValue("stage"), r.PathValue("id"), r.PathValue("taskId")
	if !validKokugoSegment(stage) || !validKokugoSegment(id) || !validKokugoSegment(taskID) {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid_id"})
		return
	}
	var body struct {
		Answer json.RawMessage `json:"answer"`
	}
	if err := decodeKokugoBody(r, &body); err != nil {
		writeKokugoDecodeError(w, err)
		return
	}
	if len(body.Answer) == 0 {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid_body"})
		return
	}
	unit, err := h.loader.GetMap(stage, id)
	if errors.Is(err, kokugo.ErrUnitNotFound) {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "not_found"})
		return
	}
	if err != nil {
		httpError(w, r, http.StatusInternalServerError, "internal", err)
		return
	}
	grade, err := kokugo.GradeTask(unit, taskID, body.Answer)
	if err != nil {
		// Stable public code only — raw grader text stays server-side.
		slog.Warn("kokugo grade failed",
			"path", r.URL.Path,
			"task_id", taskID,
			"err", err,
		)
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "grade_failed"})
		return
	}
	fb, _ := json.Marshal(grade)
	attempt, err := store.SaveKokugoTaskAttempt(r.Context(), h.db, stage, id, taskID, body.Answer, grade.Correct, fb)
	if err != nil {
		httpError(w, r, http.StatusInternalServerError, "internal", err)
		return
	}
	// Units without artifact complete when every task has an attempt.
	if !store.UnitHasArtifact(unit) {
		if ok, cerr := store.KokugoCycleComplete(r.Context(), h.db, stage, id, unit); cerr != nil {
			httpError(w, r, http.StatusInternalServerError, "internal", cerr)
			return
		} else if ok {
			if _, err := store.UpdateKokugoProgress(r.Context(), h.db, stage, id, "done", "completed"); err != nil {
				httpError(w, r, http.StatusInternalServerError, "internal", err)
				return
			}
		}
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"attempt": attempt,
		"grade":   grade,
	})
}

func (h *kokugoHandlers) putArtifact(w http.ResponseWriter, r *http.Request) {
	if !h.progressEnabled {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "progress_disabled"})
		return
	}
	stage, id := r.PathValue("stage"), r.PathValue("id")
	if !validKokugoSegment(stage) || !validKokugoSegment(id) {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid_id"})
		return
	}
	var body struct {
		Revision         int    `json:"revision"`
		Body             string `json:"body"`
		ChecklistChecked []bool `json:"checklist_checked"`
		// ExpectedVersion is required for updates; omit on first insert.
		ExpectedVersion *int `json:"expected_version"`
	}
	if err := decodeKokugoBody(r, &body); err != nil {
		writeKokugoDecodeError(w, err)
		return
	}
	if body.Revision != 0 && body.Revision != 1 {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid_revision"})
		return
	}
	unit, err := h.loader.GetMap(stage, id)
	if errors.Is(err, kokugo.ErrUnitNotFound) {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "not_found"})
		return
	}
	if err != nil {
		httpError(w, r, http.StatusInternalServerError, "internal", err)
		return
	}
	grade := kokugo.GradeArtifact(unit, body.Body, body.ChecklistChecked, body.Revision)
	checklistJSON, _ := json.Marshal(body.ChecklistChecked)
	opts := store.SaveKokugoArtifactOpts{}
	if body.ExpectedVersion != nil {
		opts.HasExpectedVersion = true
		opts.ExpectedVersion = *body.ExpectedVersion
	}
	// Progress steps track the pedagogical cycle.
	// Draft (rev 0): stay on "artifact" after save so learners can re-edit freely;
	// the UI advances to revise explicitly ("改稿へ").
	// Revision (rev 1): complete only when the full cycle is satisfied.
	// Artifact + progress share one SQLite transaction (risk-reviewer-F001).
	step, status := "artifact", "in_progress"
	passed := grade.Correct != nil && *grade.Correct
	if body.Revision == 1 {
		step = "revise"
		if passed {
			// Cycle check uses already-committed attempts/rev0; rev1 is about to be written
			// in the same transaction, so treat this save as providing rev1 when checking.
			ok, cerr := store.KokugoCycleCompleteAssumingArtifactRevision(
				r.Context(), h.db, stage, id, unit, body.Revision,
			)
			if cerr != nil {
				httpError(w, r, http.StatusInternalServerError, "internal", cerr)
				return
			}
			if ok {
				step, status = "done", "completed"
			}
		}
	}
	art, prog, err := store.SaveKokugoArtifactAndProgress(
		r.Context(), h.db, stage, id, body.Revision, body.Body, checklistJSON, opts, step, status,
	)
	if err != nil {
		if errors.Is(err, store.ErrKokugoDraftRequired) {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "draft_required"})
			return
		}
		if errors.Is(err, store.ErrKokugoStaleWrite) {
			writeJSON(w, http.StatusConflict, map[string]string{"error": "stale_write"})
			return
		}
		if errors.Is(err, store.ErrKokugoCompletedLocked) {
			writeJSON(w, http.StatusConflict, map[string]string{"error": "completed_locked"})
			return
		}
		httpError(w, r, http.StatusInternalServerError, "internal", err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"artifact": art,
		"grade":    grade,
		"progress": prog,
	})
}

func decodeKokugoBody(r *http.Request, dst any) error {
	defer r.Body.Close()
	// Read max+1 so an oversized body is distinguishable from a valid limit payload.
	data, err := io.ReadAll(io.LimitReader(r.Body, int64(maxKokugoBodyBytes)+1))
	if err != nil {
		return err
	}
	if len(data) > maxKokugoBodyBytes {
		return errKokugoBodyTooLarge
	}
	if len(data) == 0 {
		return io.ErrUnexpectedEOF
	}
	return json.Unmarshal(data, dst)
}

func writeKokugoDecodeError(w http.ResponseWriter, err error) {
	if errors.Is(err, errKokugoBodyTooLarge) {
		writeJSON(w, http.StatusRequestEntityTooLarge, map[string]string{"error": "body_too_large"})
		return
	}
	writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid_body"})
}

// validKokugoClientStep allows UI resume navigation steps only (not done/completed).
func validKokugoClientStep(step string) bool {
	if step == "" {
		return true // store defaults to predict
	}
	switch step {
	case "predict", "read", "task", "artifact", "revise":
		return true
	default:
		return strings.HasPrefix(step, "task:")
	}
}

func validKokugoSegment(s string) bool {
	if s == "" || len(s) > 64 || strings.Contains(s, "..") {
		return false
	}
	// Allow kebab-case + stage forms like e5-6
	for _, r := range s {
		if unicode.IsControl(r) {
			return false
		}
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '-' || r == '_' {
			continue
		}
		return false
	}
	return true
}
