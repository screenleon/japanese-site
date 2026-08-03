package store

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
)

// ErrNotFound is returned when a kokugo progress/artifact row is missing.
var ErrNotFound = errors.New("not found")

// KokugoUnitProgress is the learner's progress through one unit cycle.
type KokugoUnitProgress struct {
	UnitKey     string  `json:"unit_key"`
	Stage       string  `json:"stage"`
	UnitID      string  `json:"unit_id"`
	Status      string  `json:"status"` // in_progress | completed
	Step        string  `json:"step"`
	StartedAt   string  `json:"started_at"`
	UpdatedAt   string  `json:"updated_at"`
	CompletedAt *string `json:"completed_at,omitempty"`
}

// KokugoTaskAttempt is one saved task answer (latest kept by query order).
type KokugoTaskAttempt struct {
	ID           int64           `json:"id"`
	UnitKey      string          `json:"unit_key"`
	TaskID       string          `json:"task_id"`
	Answer       json.RawMessage `json:"answer"`
	Correct      *bool           `json:"correct,omitempty"`
	Feedback     json.RawMessage `json:"feedback,omitempty"`
	CreatedAt    string          `json:"created_at"`
}

// KokugoArtifact is a draft (revision=0) or post-revision (revision=1) writing.
type KokugoArtifact struct {
	UnitKey   string          `json:"unit_key"`
	Revision  int             `json:"revision"`
	Body      string          `json:"body"`
	Checklist json.RawMessage `json:"checklist"`
	// Version is a monotonic optimistic-concurrency token (starts at 1).
	Version   int             `json:"version"`
	CreatedAt string          `json:"created_at"`
	UpdatedAt string          `json:"updated_at"`
}

// KokugoUnitState bundles progress + latest attempts + artifacts for one unit.
type KokugoUnitState struct {
	Progress  *KokugoUnitProgress `json:"progress,omitempty"`
	Attempts  []KokugoTaskAttempt `json:"attempts"`
	Artifacts []KokugoArtifact    `json:"artifacts"`
}

// EnsureKokugoProgress inserts a progress row if missing and returns current state.
func EnsureKokugoProgress(ctx context.Context, db *DB, stage, unitID, step string) (KokugoUnitProgress, error) {
	unitKey := KokugoUnitKey(stage, unitID)
	if step == "" {
		step = "predict"
	}
	_, err := db.ExecContext(ctx, `
		INSERT INTO kokugo_unit_progress (unit_key, stage, unit_id, status, step)
		VALUES (?, ?, ?, 'in_progress', ?)
		ON CONFLICT(unit_key) DO NOTHING`,
		unitKey, stage, unitID, step)
	if err != nil {
		return KokugoUnitProgress{}, fmt.Errorf("ensure kokugo progress: %w", err)
	}
	return GetKokugoProgress(ctx, db, stage, unitID)
}

// UpdateKokugoProgress sets step/status for a unit (creates row if needed).
func UpdateKokugoProgress(ctx context.Context, db *DB, stage, unitID, step, status string) (KokugoUnitProgress, error) {
	unitKey := KokugoUnitKey(stage, unitID)
	if step == "" {
		step = "predict"
	}
	if status == "" {
		status = "in_progress"
	}
	if status != "in_progress" && status != "completed" {
		return KokugoUnitProgress{}, fmt.Errorf("invalid status %q", status)
	}

	var completedAt any
	if status == "completed" {
		completedAt = time.Now().UTC().Format("2006-01-02 15:04:05")
	}

	_, err := db.ExecContext(ctx, `
		INSERT INTO kokugo_unit_progress (unit_key, stage, unit_id, status, step, completed_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP)
		ON CONFLICT(unit_key) DO UPDATE SET
			status = excluded.status,
			step = excluded.step,
			updated_at = CURRENT_TIMESTAMP,
			completed_at = CASE
				WHEN excluded.status = 'completed' THEN COALESCE(kokugo_unit_progress.completed_at, excluded.completed_at)
				ELSE NULL
			END`,
		unitKey, stage, unitID, status, step, completedAt)
	if err != nil {
		return KokugoUnitProgress{}, fmt.Errorf("update kokugo progress: %w", err)
	}
	return GetKokugoProgress(ctx, db, stage, unitID)
}

// GetKokugoProgress returns progress for one unit, or ErrNotFound.
func GetKokugoProgress(ctx context.Context, db *DB, stage, unitID string) (KokugoUnitProgress, error) {
	unitKey := KokugoUnitKey(stage, unitID)
	var p KokugoUnitProgress
	var completed sql.NullString
	err := db.QueryRowContext(ctx, `
		SELECT unit_key, stage, unit_id, status, step, started_at, updated_at, completed_at
		FROM kokugo_unit_progress WHERE unit_key = ?`, unitKey).
		Scan(&p.UnitKey, &p.Stage, &p.UnitID, &p.Status, &p.Step, &p.StartedAt, &p.UpdatedAt, &completed)
	if errors.Is(err, sql.ErrNoRows) {
		return KokugoUnitProgress{}, ErrNotFound
	}
	if err != nil {
		return KokugoUnitProgress{}, err
	}
	if completed.Valid {
		p.CompletedAt = &completed.String
	}
	return p, nil
}

// ListKokugoProgress returns all unit progress rows.
func ListKokugoProgress(ctx context.Context, db *DB) ([]KokugoUnitProgress, error) {
	rows, err := db.QueryContext(ctx, `
		SELECT unit_key, stage, unit_id, status, step, started_at, updated_at, completed_at
		FROM kokugo_unit_progress
		ORDER BY updated_at DESC, unit_key`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []KokugoUnitProgress
	for rows.Next() {
		var p KokugoUnitProgress
		var completed sql.NullString
		if err := rows.Scan(&p.UnitKey, &p.Stage, &p.UnitID, &p.Status, &p.Step, &p.StartedAt, &p.UpdatedAt, &completed); err != nil {
			return nil, err
		}
		if completed.Valid {
			p.CompletedAt = &completed.String
		}
		out = append(out, p)
	}
	if out == nil {
		out = []KokugoUnitProgress{}
	}
	return out, rows.Err()
}

// SaveKokugoTaskAttempt appends a task attempt.
func SaveKokugoTaskAttempt(ctx context.Context, db *DB, stage, unitID, taskID string, answer json.RawMessage, correct *bool, feedback json.RawMessage) (KokugoTaskAttempt, error) {
	unitKey := KokugoUnitKey(stage, unitID)
	if _, err := EnsureKokugoProgress(ctx, db, stage, unitID, "task:"+taskID); err != nil {
		return KokugoTaskAttempt{}, err
	}
	var correctVal any
	if correct != nil {
		if *correct {
			correctVal = 1
		} else {
			correctVal = 0
		}
	}
	var feedbackVal any
	if len(feedback) > 0 {
		feedbackVal = string(feedback)
	}
	res, err := db.ExecContext(ctx, `
		INSERT INTO kokugo_task_attempt (unit_key, task_id, answer_json, correct, feedback_json)
		VALUES (?, ?, ?, ?, ?)`,
		unitKey, taskID, string(answer), correctVal, feedbackVal)
	if err != nil {
		return KokugoTaskAttempt{}, fmt.Errorf("save kokugo attempt: %w", err)
	}
	id, _ := res.LastInsertId()
	return KokugoTaskAttempt{
		ID:        id,
		UnitKey:   unitKey,
		TaskID:    taskID,
		Answer:    answer,
		Correct:   correct,
		Feedback:  feedback,
		CreatedAt: time.Now().UTC().Format("2006-01-02 15:04:05"),
	}, nil
}

// ListKokugoTaskAttempts returns attempts for a unit, newest first.
func ListKokugoTaskAttempts(ctx context.Context, db *DB, stage, unitID string) ([]KokugoTaskAttempt, error) {
	unitKey := KokugoUnitKey(stage, unitID)
	rows, err := db.QueryContext(ctx, `
		SELECT id, unit_key, task_id, answer_json, correct, feedback_json, created_at
		FROM kokugo_task_attempt
		WHERE unit_key = ?
		ORDER BY created_at DESC, id DESC`, unitKey)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []KokugoTaskAttempt
	for rows.Next() {
		var a KokugoTaskAttempt
		var answer string
		var correct sql.NullInt64
		var feedback sql.NullString
		if err := rows.Scan(&a.ID, &a.UnitKey, &a.TaskID, &answer, &correct, &feedback, &a.CreatedAt); err != nil {
			return nil, err
		}
		a.Answer = json.RawMessage(answer)
		if correct.Valid {
			v := correct.Int64 == 1
			a.Correct = &v
		}
		if feedback.Valid && feedback.String != "" {
			a.Feedback = json.RawMessage(feedback.String)
		}
		out = append(out, a)
	}
	if out == nil {
		out = []KokugoTaskAttempt{}
	}
	return out, rows.Err()
}

// ErrKokugoStaleWrite is returned when an optimistic concurrency check fails.
var ErrKokugoStaleWrite = errors.New("kokugo artifact stale write")

// ErrKokugoDraftRequired is returned when revision 1 is saved without revision 0.
var ErrKokugoDraftRequired = errors.New("kokugo draft required before revision")

// SaveKokugoArtifactOpts controls draft/revision invariants and concurrency.
type SaveKokugoArtifactOpts struct {
	// ExpectedVersion is required when a row already exists. The UPDATE is
	// atomic: WHERE version = ExpectedVersion; zero rows → ErrKokugoStaleWrite.
	// First insert must leave ExpectedVersion as 0 / unset.
	ExpectedVersion int
	// HasExpectedVersion distinguishes "not provided" from "version 0".
	HasExpectedVersion bool
}

// SaveKokugoArtifact upserts draft (0) or revision (1).
// Revision 1 requires an existing revision-0 draft (ADR-0005 draft→revise cycle).
// Concurrent updates use an atomic version compare-and-swap (not wall-clock).
func SaveKokugoArtifact(ctx context.Context, db *DB, stage, unitID string, revision int, body string, checklist json.RawMessage, opts SaveKokugoArtifactOpts) (KokugoArtifact, error) {
	if revision != 0 && revision != 1 {
		return KokugoArtifact{}, fmt.Errorf("invalid revision %d", revision)
	}
	unitKey := KokugoUnitKey(stage, unitID)
	if len(checklist) == 0 {
		checklist = json.RawMessage("[]")
	}

	// Draft-before-revision: full cycle requires a saved revision 0 first.
	// Check before Ensure so rejected rev=1 does not create a progress row.
	if revision == 1 {
		if _, err := GetKokugoArtifact(ctx, db, stage, unitID, 0); errors.Is(err, ErrNotFound) {
			return KokugoArtifact{}, ErrKokugoDraftRequired
		} else if err != nil {
			return KokugoArtifact{}, err
		}
	}

	if _, err := EnsureKokugoProgress(ctx, db, stage, unitID, "artifact"); err != nil {
		return KokugoArtifact{}, err
	}

	existing, err := GetKokugoArtifact(ctx, db, stage, unitID, revision)
	switch {
	case errors.Is(err, ErrNotFound):
		// First write: insert version=1. Reject clients that think a row exists.
		if opts.HasExpectedVersion {
			return KokugoArtifact{}, ErrKokugoStaleWrite
		}
		_, err = db.ExecContext(ctx, `
			INSERT INTO kokugo_artifact (unit_key, revision, body, checklist_json, version, updated_at)
			VALUES (?, ?, ?, ?, 1, CURRENT_TIMESTAMP)`,
			unitKey, revision, body, string(checklist))
		if err != nil {
			// Concurrent first-insert race: treat as stale so client reloads.
			return KokugoArtifact{}, fmt.Errorf("insert kokugo artifact: %w", err)
		}
	case err != nil:
		return KokugoArtifact{}, err
	default:
		// Existing row: require expected version and CAS update.
		if !opts.HasExpectedVersion {
			return KokugoArtifact{}, ErrKokugoStaleWrite
		}
		if opts.ExpectedVersion != existing.Version {
			return KokugoArtifact{}, ErrKokugoStaleWrite
		}
		res, err := db.ExecContext(ctx, `
			UPDATE kokugo_artifact
			SET body = ?, checklist_json = ?, version = version + 1, updated_at = CURRENT_TIMESTAMP
			WHERE unit_key = ? AND revision = ? AND version = ?`,
			body, string(checklist), unitKey, revision, opts.ExpectedVersion)
		if err != nil {
			return KokugoArtifact{}, fmt.Errorf("update kokugo artifact: %w", err)
		}
		n, _ := res.RowsAffected()
		if n != 1 {
			return KokugoArtifact{}, ErrKokugoStaleWrite
		}
	}
	return GetKokugoArtifact(ctx, db, stage, unitID, revision)
}

// GetKokugoArtifact returns one artifact revision.
func GetKokugoArtifact(ctx context.Context, db *DB, stage, unitID string, revision int) (KokugoArtifact, error) {
	unitKey := KokugoUnitKey(stage, unitID)
	var a KokugoArtifact
	var checklist string
	err := db.QueryRowContext(ctx, `
		SELECT unit_key, revision, body, checklist_json, version, created_at, updated_at
		FROM kokugo_artifact WHERE unit_key = ? AND revision = ?`, unitKey, revision).
		Scan(&a.UnitKey, &a.Revision, &a.Body, &checklist, &a.Version, &a.CreatedAt, &a.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return KokugoArtifact{}, ErrNotFound
	}
	if err != nil {
		return KokugoArtifact{}, err
	}
	a.Checklist = json.RawMessage(checklist)
	return a, nil
}

// ListKokugoArtifacts returns all artifacts for a unit.
func ListKokugoArtifacts(ctx context.Context, db *DB, stage, unitID string) ([]KokugoArtifact, error) {
	unitKey := KokugoUnitKey(stage, unitID)
	rows, err := db.QueryContext(ctx, `
		SELECT unit_key, revision, body, checklist_json, version, created_at, updated_at
		FROM kokugo_artifact WHERE unit_key = ? ORDER BY revision`, unitKey)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []KokugoArtifact
	for rows.Next() {
		var a KokugoArtifact
		var checklist string
		if err := rows.Scan(&a.UnitKey, &a.Revision, &a.Body, &checklist, &a.Version, &a.CreatedAt, &a.UpdatedAt); err != nil {
			return nil, err
		}
		a.Checklist = json.RawMessage(checklist)
		out = append(out, a)
	}
	if out == nil {
		out = []KokugoArtifact{}
	}
	return out, rows.Err()
}

// GetKokugoUnitState returns progress + attempts + artifacts (progress may be nil).
func GetKokugoUnitState(ctx context.Context, db *DB, stage, unitID string) (KokugoUnitState, error) {
	state := KokugoUnitState{Attempts: []KokugoTaskAttempt{}, Artifacts: []KokugoArtifact{}}
	p, err := GetKokugoProgress(ctx, db, stage, unitID)
	if err == nil {
		state.Progress = &p
	} else if !errors.Is(err, ErrNotFound) {
		return state, err
	}
	attempts, err := ListKokugoTaskAttempts(ctx, db, stage, unitID)
	if err != nil {
		return state, err
	}
	state.Attempts = attempts
	artifacts, err := ListKokugoArtifacts(ctx, db, stage, unitID)
	if err != nil {
		return state, err
	}
	state.Artifacts = artifacts
	return state, nil
}

// KokugoUnitKey builds the stable progress key.
func KokugoUnitKey(stage, unitID string) string {
	return strings.TrimSpace(stage) + "/" + strings.TrimSpace(unitID)
}

// ErrKokugoCycleIncomplete is returned when client requests completed status
// without the required full pedagogical cycle.
var ErrKokugoCycleIncomplete = errors.New("kokugo cycle incomplete")

// KokugoTaskIDs extracts task id strings from a unit map (loader JSON shape).
func KokugoTaskIDs(unit map[string]any) []string {
	tasks, _ := unit["tasks"].([]any)
	var ids []string
	for _, t := range tasks {
		m, ok := t.(map[string]any)
		if !ok {
			continue
		}
		id, _ := m["id"].(string)
		if id != "" {
			ids = append(ids, id)
		}
	}
	return ids
}

// UnitHasArtifact reports whether the unit defines an artifact writing step.
func UnitHasArtifact(unit map[string]any) bool {
	return unit["artifact"] != nil
}

// KokugoCycleComplete reports whether required task attempts and artifact
// revisions exist for completion (ADR-0005 full cycle).
//
// Requirements:
//   - every unit task id has at least one attempt
//   - if unit has artifact: both revision 0 and revision 1 rows exist
//   - if unit has no artifact: tasks alone are sufficient
func KokugoCycleComplete(ctx context.Context, db *DB, stage, unitID string, unit map[string]any) (bool, error) {
	ok, err := kokugoTasksComplete(ctx, db, stage, unitID, unit)
	if err != nil || !ok {
		return ok, err
	}
	if !UnitHasArtifact(unit) {
		return true, nil
	}
	if _, err := GetKokugoArtifact(ctx, db, stage, unitID, 0); err != nil {
		if errors.Is(err, ErrNotFound) {
			return false, nil
		}
		return false, err
	}
	if _, err := GetKokugoArtifact(ctx, db, stage, unitID, 1); err != nil {
		if errors.Is(err, ErrNotFound) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// KokugoCycleCompleteAssumingArtifactRevision is used while saving revision r
// inside a transaction: treat r as present after a successful write.
// For r==1, tasks + draft (rev0) must already exist; rev1 is the write in flight.
func KokugoCycleCompleteAssumingArtifactRevision(ctx context.Context, db *DB, stage, unitID string, unit map[string]any, writingRevision int) (bool, error) {
	ok, err := kokugoTasksComplete(ctx, db, stage, unitID, unit)
	if err != nil || !ok {
		return ok, err
	}
	if !UnitHasArtifact(unit) {
		return true, nil
	}
	if writingRevision == 1 {
		if _, err := GetKokugoArtifact(ctx, db, stage, unitID, 0); err != nil {
			if errors.Is(err, ErrNotFound) {
				return false, nil
			}
			return false, err
		}
		return true, nil
	}
	// writing draft alone never completes the cycle.
	return false, nil
}

func kokugoTasksComplete(ctx context.Context, db *DB, stage, unitID string, unit map[string]any) (bool, error) {
	taskIDs := KokugoTaskIDs(unit)
	if len(taskIDs) == 0 {
		return false, nil
	}
	attempts, err := ListKokugoTaskAttempts(ctx, db, stage, unitID)
	if err != nil {
		return false, err
	}
	seen := make(map[string]bool, len(attempts))
	for _, a := range attempts {
		seen[a.TaskID] = true
	}
	for _, id := range taskIDs {
		if !seen[id] {
			return false, nil
		}
	}
	return true, nil
}

// SaveKokugoArtifactAndProgress writes artifact and progress step in one transaction.
func SaveKokugoArtifactAndProgress(
	ctx context.Context,
	db *DB,
	stage, unitID string,
	revision int,
	body string,
	checklist json.RawMessage,
	opts SaveKokugoArtifactOpts,
	step, status string,
) (KokugoArtifact, KokugoUnitProgress, error) {
	if revision != 0 && revision != 1 {
		return KokugoArtifact{}, KokugoUnitProgress{}, fmt.Errorf("invalid revision %d", revision)
	}
	if status == "" {
		status = "in_progress"
	}
	if status != "in_progress" && status != "completed" {
		return KokugoArtifact{}, KokugoUnitProgress{}, fmt.Errorf("invalid status %q", status)
	}
	if step == "" {
		step = "artifact"
	}
	unitKey := KokugoUnitKey(stage, unitID)
	if len(checklist) == 0 {
		checklist = json.RawMessage("[]")
	}

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return KokugoArtifact{}, KokugoUnitProgress{}, err
	}
	defer func() { _ = tx.Rollback() }()

	// Draft-before-revision (check inside tx).
	if revision == 1 {
		var n int
		err := tx.QueryRowContext(ctx, `
			SELECT COUNT(*) FROM kokugo_artifact WHERE unit_key = ? AND revision = 0`, unitKey).Scan(&n)
		if err != nil {
			return KokugoArtifact{}, KokugoUnitProgress{}, err
		}
		if n == 0 {
			return KokugoArtifact{}, KokugoUnitProgress{}, ErrKokugoDraftRequired
		}
	}

	// Ensure progress row exists.
	_, err = tx.ExecContext(ctx, `
		INSERT INTO kokugo_unit_progress (unit_key, stage, unit_id, status, step)
		VALUES (?, ?, ?, 'in_progress', ?)
		ON CONFLICT(unit_key) DO NOTHING`,
		unitKey, stage, unitID, step)
	if err != nil {
		return KokugoArtifact{}, KokugoUnitProgress{}, fmt.Errorf("ensure progress: %w", err)
	}

	// Artifact upsert with CAS.
	var existingVersion sql.NullInt64
	err = tx.QueryRowContext(ctx, `
		SELECT version FROM kokugo_artifact WHERE unit_key = ? AND revision = ?`,
		unitKey, revision).Scan(&existingVersion)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		if opts.HasExpectedVersion {
			return KokugoArtifact{}, KokugoUnitProgress{}, ErrKokugoStaleWrite
		}
		_, err = tx.ExecContext(ctx, `
			INSERT INTO kokugo_artifact (unit_key, revision, body, checklist_json, version, updated_at)
			VALUES (?, ?, ?, ?, 1, CURRENT_TIMESTAMP)`,
			unitKey, revision, body, string(checklist))
		if err != nil {
			return KokugoArtifact{}, KokugoUnitProgress{}, fmt.Errorf("insert artifact: %w", err)
		}
	case err != nil:
		return KokugoArtifact{}, KokugoUnitProgress{}, err
	default:
		if !opts.HasExpectedVersion {
			return KokugoArtifact{}, KokugoUnitProgress{}, ErrKokugoStaleWrite
		}
		if opts.ExpectedVersion != int(existingVersion.Int64) {
			return KokugoArtifact{}, KokugoUnitProgress{}, ErrKokugoStaleWrite
		}
		res, err := tx.ExecContext(ctx, `
			UPDATE kokugo_artifact
			SET body = ?, checklist_json = ?, version = version + 1, updated_at = CURRENT_TIMESTAMP
			WHERE unit_key = ? AND revision = ? AND version = ?`,
			body, string(checklist), unitKey, revision, opts.ExpectedVersion)
		if err != nil {
			return KokugoArtifact{}, KokugoUnitProgress{}, fmt.Errorf("update artifact: %w", err)
		}
		n, _ := res.RowsAffected()
		if n != 1 {
			return KokugoArtifact{}, KokugoUnitProgress{}, ErrKokugoStaleWrite
		}
	}

	var completedAt any
	if status == "completed" {
		completedAt = time.Now().UTC().Format("2006-01-02 15:04:05")
	}
	_, err = tx.ExecContext(ctx, `
		INSERT INTO kokugo_unit_progress (unit_key, stage, unit_id, status, step, completed_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP)
		ON CONFLICT(unit_key) DO UPDATE SET
			status = excluded.status,
			step = excluded.step,
			updated_at = CURRENT_TIMESTAMP,
			completed_at = CASE
				WHEN excluded.status = 'completed' THEN COALESCE(kokugo_unit_progress.completed_at, excluded.completed_at)
				ELSE NULL
			END`,
		unitKey, stage, unitID, status, step, completedAt)
	if err != nil {
		return KokugoArtifact{}, KokugoUnitProgress{}, fmt.Errorf("update progress: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return KokugoArtifact{}, KokugoUnitProgress{}, err
	}
	art, err := GetKokugoArtifact(ctx, db, stage, unitID, revision)
	if err != nil {
		return KokugoArtifact{}, KokugoUnitProgress{}, err
	}
	prog, err := GetKokugoProgress(ctx, db, stage, unitID)
	if err != nil {
		return KokugoArtifact{}, KokugoUnitProgress{}, err
	}
	return art, prog, nil
}
