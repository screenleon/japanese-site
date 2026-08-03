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

func TestGradeArtifactLength(t *testing.T) {
	unit := loadLibraryUse(t)
	// min 80 max 120
	short := GradeArtifact(unit, "短い", []bool{true, true, true})
	if short.Correct == nil || *short.Correct {
		t.Fatalf("short should fail: %+v", short)
	}
	body := "私は入口に「今週のおすすめ」コーナーを置くことを提案します。何を読めばよいか迷う生徒が減り、本を手に取るきっかけが増えるからです。さらに教科ごとの棚表示も分かりやすくします。"
	good := GradeArtifact(unit, body, []bool{true, true, true})
	if good.Correct == nil || !*good.Correct {
		t.Fatalf("good should pass: %+v", good)
	}
}

func TestGradeArtifactBoundaries(t *testing.T) {
	// Behavior-named coverage for inclusive min/max, just-outside bounds,
	// unchecked checklist, and checklist length mismatch (library-use: 80–120, 3 checks).
	unit := loadLibraryUse(t)
	checksOK := []bool{true, true, true}

	// Exactly minimum (80 runes) — pad a base string to 80.
	minBody := padRunes("あ", 80)
	if n := len([]rune(minBody)); n != 80 {
		t.Fatalf("min pad len %d", n)
	}
	got := GradeArtifact(unit, minBody, checksOK)
	if got.Correct == nil || !*got.Correct {
		t.Fatalf("exactly min should pass: %+v", got)
	}

	// Exactly maximum (120 runes).
	maxBody := padRunes("い", 120)
	got = GradeArtifact(unit, maxBody, checksOK)
	if got.Correct == nil || !*got.Correct {
		t.Fatalf("exactly max should pass: %+v", got)
	}

	// Just outside: 79 and 121.
	got = GradeArtifact(unit, padRunes("う", 79), checksOK)
	if got.Correct == nil || *got.Correct {
		t.Fatalf("min-1 should fail: %+v", got)
	}
	got = GradeArtifact(unit, padRunes("え", 121), checksOK)
	if got.Correct == nil || *got.Correct {
		t.Fatalf("max+1 should fail: %+v", got)
	}

	// Unchecked checklist entry (length in range).
	got = GradeArtifact(unit, minBody, []bool{true, false, true})
	if got.Correct == nil || *got.Correct {
		t.Fatalf("unchecked item should fail: %+v", got)
	}

	// Checklist length mismatch.
	got = GradeArtifact(unit, minBody, []bool{true, true})
	if got.Correct == nil || *got.Correct {
		t.Fatalf("checklist length mismatch should fail: %+v", got)
	}
	got = GradeArtifact(unit, minBody, []bool{true, true, true, true})
	if got.Correct == nil || *got.Correct {
		t.Fatalf("checklist too long should fail: %+v", got)
	}

	// Omitted / empty checklist is not a free pass when unit has checklist items.
	got = GradeArtifact(unit, minBody, nil)
	if got.Correct == nil || *got.Correct {
		t.Fatalf("nil checklist should fail: %+v", got)
	}
	got = GradeArtifact(unit, minBody, []bool{})
	if got.Correct == nil || *got.Correct {
		t.Fatalf("empty checklist should fail: %+v", got)
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
