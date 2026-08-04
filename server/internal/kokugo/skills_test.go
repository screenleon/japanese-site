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
