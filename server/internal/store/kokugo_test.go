package store

import (
	"context"
	"encoding/json"
	"errors"
	"path/filepath"
	"testing"
)

func openKokugoTestDB(t *testing.T) *DB {
	t.Helper()
	db, err := Open(filepath.Join(t.TempDir(), "kokugo.sqlite"))
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })
	if err := Migrate(db); err != nil {
		t.Fatalf("Migrate: %v", err)
	}
	return db
}

func TestKokugoProgressRoundTrip(t *testing.T) {
	db := openKokugoTestDB(t)
	ctx := context.Background()

	p, err := EnsureKokugoProgress(ctx, db, "e5-6", "library-use", "predict")
	if err != nil {
		t.Fatalf("Ensure: %v", err)
	}
	if p.UnitKey != "e5-6/library-use" || p.Status != "in_progress" {
		t.Fatalf("unexpected progress: %+v", p)
	}

	p2, err := UpdateKokugoProgress(ctx, db, "e5-6", "library-use", "read", "in_progress")
	if err != nil {
		t.Fatalf("Update: %v", err)
	}
	if p2.Step != "read" {
		t.Fatalf("step=%q", p2.Step)
	}

	answer := json.RawMessage(`{"choice_id":"b"}`)
	correct := true
	fb := json.RawMessage(`{"ok":true}`)
	attempt, err := SaveKokugoTaskAttempt(ctx, db, "e5-6", "library-use", "predict-1", answer, &correct, fb)
	if err != nil {
		t.Fatalf("SaveAttempt: %v", err)
	}
	if attempt.TaskID != "predict-1" || attempt.Correct == nil || !*attempt.Correct {
		t.Fatalf("attempt: %+v", attempt)
	}

	art, err := SaveKokugoArtifact(ctx, db, "e5-6", "library-use", 0, "提案本文です。理由もあります。", json.RawMessage(`["ok"]`), SaveKokugoArtifactOpts{})
	if err != nil {
		t.Fatalf("SaveArtifact: %v", err)
	}
	if art.Revision != 0 || art.Body == "" || art.Version != 1 {
		t.Fatalf("artifact: %+v", art)
	}

	// revision without draft rejected
	if _, err := SaveKokugoArtifact(ctx, db, "e5-6", "other", 1, "改稿のみ", json.RawMessage(`[]`), SaveKokugoArtifactOpts{}); !errors.Is(err, ErrKokugoDraftRequired) {
		t.Fatalf("want draft required, got %v", err)
	}

	// update without expected version rejected
	if _, err := SaveKokugoArtifact(ctx, db, "e5-6", "library-use", 0, "blank clobber", json.RawMessage(`[]`), SaveKokugoArtifactOpts{}); !errors.Is(err, ErrKokugoStaleWrite) {
		t.Fatalf("want stale without token, got %v", err)
	}

	// same-token double write: first CAS wins, second gets stale (atomic version)
	art2, err := SaveKokugoArtifact(ctx, db, "e5-6", "library-use", 0, "改稿前の草案", json.RawMessage(`[]`), SaveKokugoArtifactOpts{
		HasExpectedVersion: true,
		ExpectedVersion:    art.Version,
	})
	if err != nil {
		t.Fatalf("SaveArtifact2: %v", err)
	}
	if art2.Body != "改稿前の草案" || art2.Version != 2 {
		t.Fatalf("body/version=%q/%d", art2.Body, art2.Version)
	}
	if _, err := SaveKokugoArtifact(ctx, db, "e5-6", "library-use", 0, "lost update", json.RawMessage(`[]`), SaveKokugoArtifactOpts{
		HasExpectedVersion: true,
		ExpectedVersion:    art.Version, // stale token from before first CAS
	}); !errors.Is(err, ErrKokugoStaleWrite) {
		t.Fatalf("want stale on reused token, got %v", err)
	}
	// Body must still be the winner of the first CAS
	got, err := GetKokugoArtifact(ctx, db, "e5-6", "library-use", 0)
	if err != nil || got.Body != "改稿前の草案" || got.Version != 2 {
		t.Fatalf("after stale: %+v err=%v", got, err)
	}

	// revision 1 after draft
	rev, err := SaveKokugoArtifact(ctx, db, "e5-6", "library-use", 1, "改稿後の本文です。", json.RawMessage(`[]`), SaveKokugoArtifactOpts{})
	if err != nil {
		t.Fatalf("revision: %v", err)
	}
	if rev.Revision != 1 || rev.Version != 1 {
		t.Fatalf("rev=%+v", rev)
	}

	p3, err := UpdateKokugoProgress(ctx, db, "e5-6", "library-use", "done", "completed")
	if err != nil {
		t.Fatalf("complete: %v", err)
	}
	if p3.Status != "completed" || p3.CompletedAt == nil {
		t.Fatalf("expected completed: %+v", p3)
	}

	state, err := GetKokugoUnitState(ctx, db, "e5-6", "library-use")
	if err != nil {
		t.Fatalf("state: %v", err)
	}
	if state.Progress == nil || len(state.Attempts) != 1 || len(state.Artifacts) != 2 {
		t.Fatalf("state: %+v", state)
	}

	list, err := ListKokugoProgress(ctx, db)
	if err != nil || len(list) != 1 {
		t.Fatalf("list: %v %#v", err, list)
	}

	_, err = GetKokugoProgress(ctx, db, "e5-6", "missing")
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("want ErrNotFound, got %v", err)
	}
}
