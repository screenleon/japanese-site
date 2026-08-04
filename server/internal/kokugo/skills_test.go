package kokugo

import (
	"encoding/json"
	"testing"
)

func TestBuildUnitSkillIndexAndAggregate(t *testing.T) {
	units := []map[string]any{
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
	idx := BuildUnitSkillIndex(units)
	if len(idx.SkillUnits["reading.summary"]) != 2 {
		t.Fatalf("summary should touch 2 units, got %v", idx.SkillUnits["reading.summary"])
	}
	if idx.TaskSkill["e5-6/u1"]["evidence-1"] != "reading.locate-evidence" {
		t.Fatalf("task skill map missing evidence")
	}

	falseV := false
	trueV := true
	attempts := []AttemptRow{
		{UnitKey: "e5-6/u1", TaskID: "predict-1", Correct: nil, ID: 1},
		{UnitKey: "e5-6/u1", TaskID: "evidence-1", Correct: &falseV, ID: 2},
		{UnitKey: "e5-6/u1", TaskID: "evidence-1", Correct: &trueV, ID: 3}, // latest wins
		{UnitKey: "e5-6/u1", TaskID: "summary-1", Correct: &falseV, ID: 4},
		{UnitKey: "e5-6/u2", TaskID: "summary-1", Correct: &falseV, ID: 5},
	}
	artifacts := []ArtifactRow{
		{UnitKey: "e5-6/u1", Revision: 0, Body: "下書き", ChecklistAllTrue: false},
	}
	progress := []ProgressRow{
		{UnitKey: "e5-6/u1", Stage: "e5-6", UnitID: "u1", Status: "in_progress"},
	}

	m := AggregateSkillMap(idx, attempts, artifacts, progress)
	if len(m.Skills) != len(AllSkills) {
		t.Fatalf("expected %d skills, got %d", len(AllSkills), len(m.Skills))
	}

	byID := map[string]SkillStat{}
	for _, s := range m.Skills {
		byID[s.Skill] = s
	}

	// evidence: latest correct → graded 1 correct 1 → practiced (single graded high acc but <2 → practiced)
	ev := byID["reading.locate-evidence"]
	if ev.Graded != 1 || ev.Correct != 1 || ev.Status != "practiced" {
		t.Fatalf("evidence unexpected: %+v", ev)
	}

	// summary: 2 wrong → weak
	sum := byID["reading.summary"]
	if sum.Graded != 2 || sum.Correct != 0 || sum.Status != "weak" {
		t.Fatalf("summary unexpected: %+v", sum)
	}

	// predict ungraded → practiced
	if byID["reading.predict"].Status != "practiced" || byID["reading.predict"].Graded != 0 {
		t.Fatalf("predict unexpected: %+v", byID["reading.predict"])
	}

	// writing claim-reason practiced via draft
	if byID["writing.claim-reason"].Practiced != 1 {
		t.Fatalf("claim-reason unexpected: %+v", byID["writing.claim-reason"])
	}

	// structure unseen
	if byID["reading.structure"].Status != "unseen" {
		t.Fatalf("structure should be unseen: %+v", byID["reading.structure"])
	}

	// review queue should recommend incomplete units covering weak/unseen skills
	if len(m.ReviewQueue) == 0 {
		t.Fatal("expected non-empty review queue")
	}
	for _, item := range m.ReviewQueue {
		if item.UnitCompleted {
			t.Fatalf("review item should be incomplete: %+v", item)
		}
		if len(item.TargetSkills) == 0 {
			t.Fatalf("review item missing target skills: %+v", item)
		}
	}
}

func TestChecklistAllTrue(t *testing.T) {
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
	// critic-F001 / qa-tester-F001: unknown skill must not panic and must not
	// enter the skill map index.
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
	// Must not panic.
	idx := BuildUnitSkillIndex(units)
	if _, ok := idx.TaskSkill["e5-6/bad"]["bad-1"]; ok {
		t.Fatalf("unknown skill must not be indexed: %+v", idx.TaskSkill["e5-6/bad"])
	}
	if idx.TaskSkill["e5-6/bad"]["ok-1"] != "reading.summary" {
		t.Fatalf("canonical skill lost: %+v", idx.TaskSkill["e5-6/bad"])
	}
	for skill, unitsTouching := range idx.SkillUnits {
		if !isKnownSkill(skill) {
			t.Fatalf("SkillUnits has unknown key %q: %v", skill, unitsTouching)
		}
		for _, uk := range unitsTouching {
			if uk == "e5-6/bad" && skill == "reading.telepathy" {
				t.Fatal("unknown skill should not appear in SkillUnits")
			}
		}
	}
	// Aggregate over unknown skill in attempts must also stay stable.
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
	// qa-tester-F001: mutation-sensitive status boundaries.
	cases := []struct {
		name                string
		practiced, graded, correct int
		want                string
	}{
		{"unseen", 0, 0, 0, "unseen"},
		{"ungraded practiced", 1, 0, 0, "practiced"},
		{"weak low accuracy", 2, 2, 0, "weak"}, // acc 0
		{"weak boundary 0.5", 2, 2, 1, "weak"},  // acc 0.5 < 0.6
		{"practiced mid accuracy", 2, 2, 2, "strong"}, // acc 1.0, graded>=2
		{"practiced single high", 1, 1, 1, "practiced"}, // need graded>=2 for strong
		{"practiced accuracy 0.6", 5, 5, 3, "practiced"}, // 0.6 not < 0.6, not >=0.8
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
	// qa-tester-F001: completed units skipped; limit capped; higher score first.
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
	// Unit a covers 2 weak skills → should rank first.
	if q[0].UnitID != "a" {
		t.Fatalf("want unit a first (score 2), got %+v", q)
	}
	// limit 0 → empty
	if got := buildReviewQueue(idx, weak, completed, 0); len(got) != 0 {
		t.Fatalf("limit 0 → empty, got %+v", got)
	}
}
