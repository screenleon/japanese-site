package kokugo

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func loadLibraryUse(t *testing.T) map[string]any {
	t.Helper()
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}
	// server/internal/kokugo → repo root
	root := filepath.Clean(filepath.Join(filepath.Dir(file), "../../.."))
	path := filepath.Join(root, "data/corpus/kokugo/e5-6/library-use.json")
	// When tests run from module root japanese-site, path is server/data/...
	if _, err := os.Stat(path); err != nil {
		path = filepath.Join(root, "server/data/corpus/kokugo/e5-6/library-use.json")
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read fixture: %v", err)
	}
	var m map[string]any
	if err := json.Unmarshal(data, &m); err != nil {
		t.Fatalf("json: %v", err)
	}
	return m
}

func TestGradeSummaryAndEvidence(t *testing.T) {
	unit := loadLibraryUse(t)

	ok, err := GradeTask(unit, "summary-1", json.RawMessage(`{"choice_id":"b"}`))
	if err != nil || ok.Correct == nil || !*ok.Correct {
		t.Fatalf("summary correct: %+v err=%v", ok, err)
	}
	bad, err := GradeTask(unit, "summary-1", json.RawMessage(`{"choice_id":"a"}`))
	if err != nil || bad.Correct == nil || *bad.Correct {
		t.Fatalf("summary wrong: %+v err=%v", bad, err)
	}

	ev, err := GradeTask(unit, "evidence-1", json.RawMessage(`{"quotes":["まず探しやすさを改善する必要があります。"]}`))
	if err != nil || ev.Correct == nil || !*ev.Correct {
		t.Fatalf("evidence: %+v err=%v", ev, err)
	}

	roles, err := GradeTask(unit, "structure-1", json.RawMessage(`{"roles":["問題","原因","提案","結論"]}`))
	if err != nil || roles.Correct == nil || !*roles.Correct {
		t.Fatalf("roles: %+v err=%v", roles, err)
	}

	pred, err := GradeTask(unit, "predict-1", json.RawMessage(`{"choice_id":"b"}`))
	if err != nil || pred.Correct != nil {
		t.Fatalf("predict should be ungraded: %+v err=%v", pred, err)
	}
}

func TestGradeArtifactProgressiveWriting(t *testing.T) {
	// library-use uses min_chars=0 max_chars=0: only non-empty is required on draft.
	unit := loadLibraryUse(t)
	checksOK := []bool{true, true, true}

	got := GradeArtifact(unit, "   ", checksOK, 0)
	if got.Correct == nil || *got.Correct {
		t.Fatalf("empty draft should fail: %+v", got)
	}
	got = GradeArtifact(unit, "", checksOK, 1)
	if got.Correct == nil || *got.Correct {
		t.Fatalf("empty revision should fail: %+v", got)
	}

	got = GradeArtifact(unit, "提案です。", nil, 0)
	if got.Correct == nil || !*got.Correct {
		t.Fatalf("short draft without checklist should pass: %+v", got)
	}

	long := padRunes("あ", 500)
	got = GradeArtifact(unit, long, nil, 0)
	if got.Correct == nil || !*got.Correct {
		t.Fatalf("long draft should pass with max=0: %+v", got)
	}

	got = GradeArtifact(unit, "提案です。理由もあります。", []bool{true, false, true}, 1)
	if got.Correct == nil || *got.Correct {
		t.Fatalf("revision with unchecked item should fail: %+v", got)
	}
	got = GradeArtifact(unit, "提案です。理由もあります。", checksOK, 1)
	if got.Correct == nil || !*got.Correct {
		t.Fatalf("revision with full checklist should pass: %+v", got)
	}
	got = GradeArtifact(unit, "提案です。", nil, 1)
	if got.Correct == nil || *got.Correct {
		t.Fatalf("revision with nil checklist should fail: %+v", got)
	}
}

func TestGradeArtifactOptionalBounds(t *testing.T) {
	unit := map[string]any{
		"artifact": map[string]any{
			"kind":      "short-proposal",
			"min_chars": float64(5),
			"max_chars": float64(10),
			"checklist": []any{"a"},
		},
	}
	got := GradeArtifact(unit, "abcd", []bool{true}, 0)
	if got.Correct == nil || *got.Correct {
		t.Fatalf("below min should fail when min>0: %+v", got)
	}
	got = GradeArtifact(unit, "abcde", []bool{true}, 0)
	if got.Correct == nil || !*got.Correct {
		t.Fatalf("at min draft should pass: %+v", got)
	}
	got = GradeArtifact(unit, padRunes("x", 11), []bool{true}, 0)
	if got.Correct == nil || *got.Correct {
		t.Fatalf("above max should fail when max>0: %+v", got)
	}
}

func padRunes(ch string, n int) string {
	r := []rune(ch)
	if len(r) == 0 {
		r = []rune("あ")
	}
	out := make([]rune, n)
	for i := 0; i < n; i++ {
		out[i] = r[0]
	}
	return string(out)
}
