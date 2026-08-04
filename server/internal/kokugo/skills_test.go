package kokugo

import (
	"encoding/json"
	"testing"
)

func sampleUnits() []map[string]any {
	return []map[string]any{
		{
			"id":       "u1",
			"stage":    "e5-6",
			"title_ja": "ユニット1",
			"genre":    "story",
			"tasks": []any{
				map[string]any{"id": "predict-1", "skill": "reading.predict", "kind": "predict"},
				map[string]any{"id": "evidence-1", "skill": "reading.locate-evidence", "kind": "evidence-highlight"},
				map[string]any{"id": "summary-1", "skill": "reading.summary", "kind": "summary-choice"},
			},
			"artifact": map[string]any{"kind": "summary"},
		},
		{
			"id":       "u2",
			"stage":    "e5-6",
			"title_ja": "ユニット2",
			"genre":    "opinion",
			"tasks": []any{
				map[string]any{"id": "structure-1", "skill": "reading.structure", "kind": "paragraph-role"},
				map[string]any{"id": "summary-1", "skill": "reading.summary", "kind": "summary-choice"},
			},
			"artifact": map[string]any{"kind": "short-proposal"},
		},
	}
}

func TestBuildUnitSkillIndexMapsCanonicalTasks(t *testing.T) {
	// Behavior: BuildUnitSkillIndex maps canonical task skills onto unit keys.
	// Steps:
	// 1. Arrange two units with shared summary + distinct skills.
	// 2. Act BuildUnitSkillIndex.
	// 3. Assert task→skill map and multi-unit SkillUnits membership.
	idx := BuildUnitSkillIndex(sampleUnits())
	if len(idx.SkillUnits["reading.summary"]) != 2 {
		t.Fatalf("summary should touch 2 units, got %v", idx.SkillUnits["reading.summary"])
	}
	if idx.TaskSkill["e5-6/u1"]["evidence-1"] != "reading.locate-evidence" {
		t.Fatalf("task skill map missing evidence")
	}
}

func TestAggregateSkillMapLatestAttemptWins(t *testing.T) {
	// Behavior: AggregateSkillMap uses the highest attempt id as the latest grade.
	// Steps:
	// 1. Arrange index + two evidence attempts (wrong then correct).
	// 2. Act AggregateSkillMap.
	// 3. Assert evidence graded 1/1 practiced (latest correct wins).
	idx := BuildUnitSkillIndex(sampleUnits())
	falseV := false
	trueV := true
	m := AggregateSkillMap(idx, []AttemptRow{
		{UnitKey: "e5-6/u1", TaskID: "evidence-1", Correct: &falseV, ID: 2},
		{UnitKey: "e5-6/u1", TaskID: "evidence-1", Correct: &trueV, ID: 3},
	}, nil, nil)
	byID := map[string]SkillStat{}
	for _, s := range m.Skills {
		byID[s.Skill] = s
	}
	ev := byID["reading.locate-evidence"]
	if ev.Graded != 1 || ev.Correct != 1 || ev.Status != "practiced" {
		t.Fatalf("evidence latest-wins unexpected: %+v", ev)
	}
}

func TestAggregateSkillMapStatusAndWriting(t *testing.T) {
	// Behavior: aggregate status mixes graded accuracy, ungraded predict, and draft writing.
	// Steps:
	// 1. Arrange mixed attempts + one draft artifact on incomplete progress.
	// 2. Act AggregateSkillMap.
	// 3. Assert weak summary, practiced predict, practiced claim-reason, unseen structure.
	idx := BuildUnitSkillIndex(sampleUnits())
	falseV := false
	m := AggregateSkillMap(idx, []AttemptRow{
		{UnitKey: "e5-6/u1", TaskID: "predict-1", Correct: nil, ID: 1},
		{UnitKey: "e5-6/u1", TaskID: "summary-1", Correct: &falseV, ID: 4},
		{UnitKey: "e5-6/u2", TaskID: "summary-1", Correct: &falseV, ID: 5},
	}, []ArtifactRow{
		{UnitKey: "e5-6/u1", Revision: 0, Body: "下書き", ChecklistAllTrue: false},
	}, []ProgressRow{
		{UnitKey: "e5-6/u1", Stage: "e5-6", UnitID: "u1", Status: "in_progress"},
	})
	byID := map[string]SkillStat{}
	for _, s := range m.Skills {
		byID[s.Skill] = s
	}
	if sum := byID["reading.summary"]; sum.Graded != 2 || sum.Correct != 0 || sum.Status != "weak" {
		t.Fatalf("summary unexpected: %+v", sum)
	}
	if p := byID["reading.predict"]; p.Status != "practiced" || p.Graded != 0 {
		t.Fatalf("predict unexpected: %+v", p)
	}
	if byID["writing.claim-reason"].Practiced != 1 {
		t.Fatalf("claim-reason unexpected: %+v", byID["writing.claim-reason"])
	}
	if byID["reading.structure"].Status != "unseen" {
		t.Fatalf("structure should be unseen: %+v", byID["reading.structure"])
	}
	if len(m.ReviewQueue) == 0 {
		t.Fatal("expected non-empty review queue for weak/unseen skills")
	}
}

func TestChecklistAllTrue(t *testing.T) {
	// Behavior: ChecklistAllTrue requires a non-empty all-true bool array.
	// Steps:
	// 1. Arrange three JSON checklist payloads.
	// 2. Act ChecklistAllTrue on each.
	// 3. Assert only [true,true] is accepted.
	if !ChecklistAllTrue(json.RawMessage(`[true,true]`)) {
		t.Fatal("expected all true")
	}
	if ChecklistAllTrue(json.RawMessage(`[true,false]`)) {
		t.Fatal("expected false on mixed")
	}
	if ChecklistAllTrue(json.RawMessage(`[]`)) {
		t.Fatal("empty is false")
	}
}

func TestBuildUnitSkillIndexIgnoresUnknownSkill(t *testing.T) {
	// Behavior: unknown task skills are skipped (no panic, no index entry).
	// Steps:
	// 1. Arrange a unit with one canonical skill, one unknown skill, empty skill, non-object task.
	// 2. Act BuildUnitSkillIndex (must not panic).
	// 3. Assert only the canonical task is indexed; AggregateSkillMap stays stable.
	units := []map[string]any{
		{
			"id":       "bad",
			"stage":    "e5-6",
			"title_ja": "不正スキル",
			"genre":    "story",
			"tasks": []any{
				map[string]any{"id": "ok-1", "skill": "reading.summary", "kind": "summary-choice"},
				map[string]any{"id": "bad-1", "skill": "reading.telepathy", "kind": "summary-choice"},
				map[string]any{"id": "empty-skill", "skill": "", "kind": "predict"},
				"not-an-object",
			},
		},
	}
	idx := BuildUnitSkillIndex(units)
	if _, ok := idx.TaskSkill["e5-6/bad"]["bad-1"]; ok {
		t.Fatalf("unknown skill must not be indexed: %+v", idx.TaskSkill["e5-6/bad"])
	}
	if idx.TaskSkill["e5-6/bad"]["ok-1"] != "reading.summary" {
		t.Fatalf("canonical skill lost: %+v", idx.TaskSkill["e5-6/bad"])
	}
	for skill := range idx.SkillUnits {
		if !isKnownSkill(skill) {
			t.Fatalf("SkillUnits has unknown key %q", skill)
		}
	}
	falseV := false
	m := AggregateSkillMap(idx, []AttemptRow{
		{UnitKey: "e5-6/bad", TaskID: "bad-1", Correct: &falseV, ID: 1},
		{UnitKey: "e5-6/bad", TaskID: "ok-1", Correct: &falseV, ID: 2},
	}, nil, nil)
	if len(m.Skills) != len(AllSkills) {
		t.Fatalf("want %d skills, got %d", len(AllSkills), len(m.Skills))
	}
}

func TestSkillStatusThresholds(t *testing.T) {
	// Behavior: skillStatus applies unseen / practiced / weak / strong thresholds.
	// Steps:
	// 1. Arrange practiced/graded/correct triples at each boundary.
	// 2. Act skillStatus for each triple.
	// 3. Assert the documented status string.
	cases := []struct {
		name                       string
		practiced, graded, correct int
		want                       string
	}{
		{"unseen", 0, 0, 0, "unseen"},
		{"ungraded practiced", 1, 0, 0, "practiced"},
		{"weak low accuracy", 2, 2, 0, "weak"},
		{"weak boundary 0.5", 2, 2, 1, "weak"},
		{"strong full accuracy", 2, 2, 2, "strong"},
		{"practiced single high", 1, 1, 1, "practiced"},
		{"practiced accuracy 0.6", 5, 5, 3, "practiced"},
		{"strong exact 0.8", 5, 5, 4, "strong"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := skillStatus(tc.practiced, tc.graded, tc.correct)
			if got != tc.want {
				t.Fatalf("skillStatus(%d,%d,%d)=%q want %q",
					tc.practiced, tc.graded, tc.correct, got, tc.want)
			}
		})
	}
}

func TestBuildReviewQueueCompletedLimitOrder(t *testing.T) {
	// Behavior: review queue skips completed units, ranks by weak-skill hit count, and respects limit.
	// Steps:
	// 1. Arrange four units with overlapping weak skills; mark one completed.
	// 2. Act buildReviewQueue with limit 2.
	// 3. Assert completed excluded, highest-score unit first, and limit 0 is empty.
	idx := BuildUnitSkillIndex([]map[string]any{
		{
			"id": "a", "stage": "e5-6", "title_ja": "A", "genre": "story",
			"tasks": []any{
				map[string]any{"id": "s1", "skill": "reading.summary"},
				map[string]any{"id": "e1", "skill": "reading.locate-evidence"},
			},
		},
		{
			"id": "b", "stage": "e5-6", "title_ja": "B", "genre": "opinion",
			"tasks": []any{
				map[string]any{"id": "s1", "skill": "reading.summary"},
			},
		},
		{
			"id": "c", "stage": "e5-6", "title_ja": "C", "genre": "poetry",
			"tasks": []any{
				map[string]any{"id": "e1", "skill": "reading.locate-evidence"},
			},
		},
		{
			"id": "done", "stage": "e5-6", "title_ja": "Done", "genre": "expository",
			"tasks": []any{
				map[string]any{"id": "s1", "skill": "reading.summary"},
				map[string]any{"id": "e1", "skill": "reading.locate-evidence"},
			},
		},
	})
	weak := []string{"reading.summary", "reading.locate-evidence"}
	completed := map[string]bool{"e5-6/done": true}

	q := buildReviewQueue(idx, weak, completed, 2)
	if len(q) != 2 {
		t.Fatalf("limit 2 → want 2 items, got %d (%+v)", len(q), q)
	}
	for _, item := range q {
		if item.UnitID == "done" {
			t.Fatalf("completed unit must be excluded: %+v", item)
		}
	}
	if q[0].UnitID != "a" {
		t.Fatalf("want unit a first (score 2), got %+v", q)
	}
	if got := buildReviewQueue(idx, weak, completed, 0); len(got) != 0 {
		t.Fatalf("limit 0 → empty, got %+v", got)
	}
}
