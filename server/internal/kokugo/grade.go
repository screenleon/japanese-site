package kokugo

import (
	"encoding/json"
	"fmt"
	"strings"
	"unicode/utf8"
)

// GradeResult is the deterministic grade for one task or artifact check.
type GradeResult struct {
	Correct      *bool          `json:"correct,omitempty"`
	Explanation  string         `json:"explanation_ja"`
	Detail       map[string]any `json:"detail,omitempty"`
}

// GradeTask grades a learner answer against a unit task (no LLM).
func GradeTask(unit map[string]any, taskID string, answer json.RawMessage) (GradeResult, error) {
	task, err := findTask(unit, taskID)
	if err != nil {
		return GradeResult{}, err
	}
	kind, _ := task["kind"].(string)
	payload, _ := task["payload"].(map[string]any)

	switch kind {
	case "predict":
		// Prediction is pedagogical priming — always accepted if a choice is present.
		var a struct {
			ChoiceID string `json:"choice_id"`
			FreeText string `json:"free_text"`
		}
		_ = json.Unmarshal(answer, &a)
		if a.ChoiceID == "" && strings.TrimSpace(a.FreeText) == "" {
			return GradeResult{}, fmt.Errorf("predict requires choice_id or free_text")
		}
		return GradeResult{
			Correct:     nil,
			Explanation: "予測を記録しました。本文を読んで考えを確かめましょう。",
			Detail:      map[string]any{"choice_id": a.ChoiceID},
		}, nil

	case "summary-choice":
		var a struct {
			ChoiceID string `json:"choice_id"`
		}
		if err := json.Unmarshal(answer, &a); err != nil || a.ChoiceID == "" {
			return GradeResult{}, fmt.Errorf("summary-choice requires choice_id")
		}
		correctID, _ := payload["correct_id"].(string)
		ok := a.ChoiceID == correctID
		expl := "正しい要約を選べました。"
		if !ok {
			expl = "要約の選択が本文の主張とずれています。筆者の提案と理由をもう一度確認しましょう。"
		}
		return GradeResult{
			Correct:     boolPtr(ok),
			Explanation: expl,
			Detail:      map[string]any{"expected": correctID, "got": a.ChoiceID},
		}, nil

	case "paragraph-role":
		var a struct {
			Roles []string `json:"roles"`
		}
		if err := json.Unmarshal(answer, &a); err != nil {
			return GradeResult{}, fmt.Errorf("paragraph-role requires roles[]")
		}
		goldAny, _ := payload["gold_by_paragraph_index"].([]any)
		gold := make([]string, len(goldAny))
		for i, g := range goldAny {
			gold[i], _ = g.(string)
		}
		ok := len(a.Roles) == len(gold)
		if ok {
			for i := range gold {
				if a.Roles[i] != gold[i] {
					ok = false
					break
				}
			}
		}
		expl := "各段落の役割を正しく整理できました。"
		if !ok {
			expl = "段落の役割に誤りがあります。問題→原因→提案→結論の流れを見直しましょう。"
		}
		return GradeResult{
			Correct:     boolPtr(ok),
			Explanation: expl,
			Detail:      map[string]any{"expected": gold, "got": a.Roles},
		}, nil

	case "evidence-highlight":
		var a struct {
			// Selected surface quotes from the passage.
			Quotes []string `json:"quotes"`
			// Optional free selection string (concatenated).
			Selected string `json:"selected"`
		}
		if err := json.Unmarshal(answer, &a); err != nil {
			return GradeResult{}, fmt.Errorf("evidence-highlight requires quotes or selected")
		}
		goldAny, _ := payload["gold_quotes"].([]any)
		var gold []string
		for _, g := range goldAny {
			if s, ok := g.(string); ok {
				gold = append(gold, s)
			}
		}
		haystack := strings.Join(a.Quotes, "\n")
		if a.Selected != "" {
			haystack = haystack + "\n" + a.Selected
		}
		// Normalize whitespace for fuzzy-ish containment.
		normHay := compactSpace(haystack)
		matched := 0
		for _, g := range gold {
			if strings.Contains(normHay, compactSpace(g)) {
				matched++
			}
		}
		ok := matched == len(gold) && len(gold) > 0
		// Also accept if user picked exactly one gold quote among options UX.
		if !ok && len(a.Quotes) == 1 {
			for _, g := range gold {
				if compactSpace(a.Quotes[0]) == compactSpace(g) {
					ok = true
					matched = 1
					break
				}
			}
		}
		expl := "根拠となる文を正しく選べました。"
		if !ok {
			expl = "根拠の引用が足りないかずれています。筆者の提案が書かれている文を探しましょう。"
		}
		return GradeResult{
			Correct:     boolPtr(ok),
			Explanation: expl,
			Detail:      map[string]any{"matched": matched, "gold_count": len(gold)},
		}, nil

	default:
		return GradeResult{}, fmt.Errorf("unsupported task kind %q", kind)
	}
}

// GradeArtifact validates a draft (revision 0) or revision (revision 1) body.
//
// Progressive-writing policy:
//   - Empty body always fails.
//   - min_chars / max_chars are optional: 0 means "no limit" so learners can grow
//     from short notes. When set (>0), bounds still apply for units that want them.
//   - Checklist is hard only on revision 1 (completion). Draft may save without ticks.
func GradeArtifact(unit map[string]any, body string, checklistChecked []bool, revision int) GradeResult {
	art, _ := unit["artifact"].(map[string]any)
	if art == nil {
		return GradeResult{Correct: boolPtr(true), Explanation: "この単元に作品課題はありません。"}
	}
	minChars, _ := asInt(art["min_chars"])
	maxChars, _ := asInt(art["max_chars"])
	n := utf8.RuneCountInString(strings.TrimSpace(body))
	if n == 0 {
		return GradeResult{
			Correct:     boolPtr(false),
			Explanation: "本文が空です。少しでも書いてから保存してください。",
			Detail: map[string]any{
				"chars":    0,
				"min":      minChars,
				"max":      maxChars,
				"len_ok":   false,
				"check_ok": true,
				"revision": revision,
			},
		}
	}
	lenOK := true
	if minChars > 0 && n < minChars {
		lenOK = false
	}
	if maxChars > 0 && n > maxChars {
		lenOK = false
	}
	checks, _ := art["checklist"].([]any)
	// Draft (rev 0): checklist is self-reminder only.
	// Revision (rev 1): every declared checklist item must be present and true.
	checkOK := true
	if revision >= 1 && len(checks) > 0 {
		if len(checklistChecked) != len(checks) {
			checkOK = false
		} else {
			for _, c := range checklistChecked {
				if !c {
					checkOK = false
					break
				}
			}
		}
	}
	ok := lenOK && checkOK
	expl := "保存できました。"
	if revision >= 1 && ok {
		expl = "改稿の条件を満たしています。"
	}
	if !lenOK {
		if minChars > 0 && maxChars > 0 {
			expl = fmt.Sprintf("字数が範囲外です（現在 %d 字、目安 %d〜%d 字）。", n, minChars, maxChars)
		} else if minChars > 0 {
			expl = fmt.Sprintf("字数が少なすぎます（現在 %d 字、目安 %d 字以上）。", n, minChars)
		} else {
			expl = fmt.Sprintf("字数が多すぎます（現在 %d 字、目安 %d 字以下）。", n, maxChars)
		}
	} else if !checkOK {
		expl = "チェックリストの項目をすべて確認してから保存してください。"
	} else if revision == 0 {
		expl = "下書きを保存しました。何度でも書き直してかまいません。"
	}
	return GradeResult{
		Correct:     boolPtr(ok),
		Explanation: expl,
		Detail: map[string]any{
			"chars":    n,
			"min":      minChars,
			"max":      maxChars,
			"len_ok":   lenOK,
			"check_ok": checkOK,
			"revision": revision,
		},
	}
}

func findTask(unit map[string]any, taskID string) (map[string]any, error) {
	tasks, _ := unit["tasks"].([]any)
	for _, t := range tasks {
		m, ok := t.(map[string]any)
		if !ok {
			continue
		}
		if id, _ := m["id"].(string); id == taskID {
			return m, nil
		}
	}
	return nil, fmt.Errorf("task %q not found", taskID)
}

func boolPtr(v bool) *bool { return &v }

func compactSpace(s string) string {
	return strings.Join(strings.Fields(s), "")
}

func asInt(v any) (int, bool) {
	switch n := v.(type) {
	case float64:
		return int(n), true
	case int:
		return n, true
	case json.Number:
		i, err := n.Int64()
		return int(i), err == nil
	default:
		return 0, false
	}
}
