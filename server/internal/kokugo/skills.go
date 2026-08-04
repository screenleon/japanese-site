// Package kokugo — skill map aggregation for Track B (JS-136).
// Derives skill progress from unit L1 JSON + task attempts / artifacts.
// Does not write to cloze SRS tables.
package kokugo

import (
	"encoding/json"
	"sort"
	"strings"
)

// Skill IDs match web/src/kokugoTypes.ts KOKUGO_SKILLS / ADR-0005.
var AllSkills = []string{
	"reading.predict",
	"reading.locate-evidence",
	"reading.structure",
	"reading.summary",
	"writing.claim-reason",
	"writing.revision",
}

// SkillLabelJa is the learner-facing Japanese label for each skill.
var SkillLabelJa = map[string]string{
	"reading.predict":         "予測する",
	"reading.locate-evidence": "根拠を探す",
	"reading.structure":       "段落の役割",
	"reading.summary":         "要約する",
	"writing.claim-reason":    "主張と理由",
	"writing.revision":        "改稿する",
}

// AttemptRow is a minimal attempt projection for skill aggregation.
type AttemptRow struct {
	UnitKey string
	TaskID  string
	// Correct is nil when the attempt was ungraded (e.g. predict).
	Correct *bool
	// CreatedAt is used only for latest-wins ordering (ISO-ish string, lexical OK for our timestamps).
	CreatedAt string
	ID        int64
}

// ArtifactRow is a minimal artifact projection for writing skills.
type ArtifactRow struct {
	UnitKey  string
	Revision int
	Body     string
	// ChecklistAllTrue is true when every checklist item is checked (JSON bool array).
	ChecklistAllTrue bool
}

// ProgressRow is a minimal progress projection.
type ProgressRow struct {
	UnitKey string
	Stage   string
	UnitID  string
	Status  string
}

// UnitSkillIndex maps unit_key → task_id → skill, plus writing skill presence.
type UnitSkillIndex struct {
	// TaskSkill[unitKey][taskID] = skill
	TaskSkill map[string]map[string]string
	// UnitMeta[unitKey] = summary for review queue
	UnitMeta map[string]UnitMeta
	// WritingSkills[unitKey] lists writing skills the unit exercises (artifact present).
	WritingSkills map[string][]string
	// SkillUnits[skill] = unit keys that exercise the skill
	SkillUnits map[string][]string
}

// UnitMeta is enough to recommend a unit without loading full JSON again.
type UnitMeta struct {
	Stage   string
	UnitID  string
	TitleJa string
	Genre   string
}

// SkillStat is one skill's aggregate for the skill map API.
type SkillStat struct {
	Skill         string   `json:"skill"`
	LabelJa       string   `json:"label_ja"`
	Status        string   `json:"status"` // unseen | practiced | weak | strong
	Graded        int      `json:"graded"`
	Correct       int      `json:"correct"`
	// Accuracy is nil when Graded == 0 (including ungraded-only skills like predict).
	Accuracy      *float64 `json:"accuracy,omitempty"`
	// Practiced counts any attempt / writing artifact touch (includes ungraded).
	Practiced     int      `json:"practiced"`
	UnitsTouching []string `json:"units_touching"`
}

// ReviewItem is a unit recommended to practice weak skills.
type ReviewItem struct {
	Stage         string   `json:"stage"`
	UnitID        string   `json:"unit_id"`
	TitleJa       string   `json:"title_ja"`
	Genre         string   `json:"genre"`
	TargetSkills  []string `json:"target_skills"`
	UnitCompleted bool     `json:"unit_completed"`
}

// SkillMap is the GET /api/kokugo/skills payload.
type SkillMap struct {
	Skills      []SkillStat  `json:"skills"`
	ReviewQueue []ReviewItem `json:"review_queue"`
}

// BuildUnitSkillIndex walks loaded unit maps (from Loader.GetMap-style objects).
// units: slice of full unit JSON maps; each must have id, stage, title_ja, tasks, optional artifact.
func BuildUnitSkillIndex(units []map[string]any) UnitSkillIndex {
	idx := UnitSkillIndex{
		TaskSkill:     map[string]map[string]string{},
		UnitMeta:      map[string]UnitMeta{},
		WritingSkills: map[string][]string{},
		SkillUnits:    map[string][]string{},
	}
	for _, skill := range AllSkills {
		idx.SkillUnits[skill] = []string{}
	}
	seenSkillUnit := map[string]map[string]bool{}
	for _, skill := range AllSkills {
		seenSkillUnit[skill] = map[string]bool{}
	}

	for _, u := range units {
		id, _ := u["id"].(string)
		stage, _ := u["stage"].(string)
		title, _ := u["title_ja"].(string)
		genre, _ := u["genre"].(string)
		if id == "" || stage == "" {
			continue
		}
		unitKey := stage + "/" + id
		idx.UnitMeta[unitKey] = UnitMeta{Stage: stage, UnitID: id, TitleJa: title, Genre: genre}
		idx.TaskSkill[unitKey] = map[string]string{}

		if tasks, ok := u["tasks"].([]any); ok {
			for _, t := range tasks {
				tm, ok := t.(map[string]any)
				if !ok {
					continue
				}
				taskID, _ := tm["id"].(string)
				skill, _ := tm["skill"].(string)
				if taskID == "" || skill == "" {
					continue
				}
				idx.TaskSkill[unitKey][taskID] = skill
				if !seenSkillUnit[skill][unitKey] {
					seenSkillUnit[skill][unitKey] = true
					idx.SkillUnits[skill] = append(idx.SkillUnits[skill], unitKey)
				}
			}
		}
		if u["artifact"] != nil {
			ws := []string{"writing.claim-reason", "writing.revision"}
			idx.WritingSkills[unitKey] = ws
			for _, skill := range ws {
				if !seenSkillUnit[skill][unitKey] {
					seenSkillUnit[skill][unitKey] = true
					idx.SkillUnits[skill] = append(idx.SkillUnits[skill], unitKey)
				}
			}
		}
	}
	for skill := range idx.SkillUnits {
		sort.Strings(idx.SkillUnits[skill])
	}
	return idx
}

// AggregateSkillMap builds skill stats + a small review queue.
// latest-wins per (unit_key, task_id) for graded accuracy.
func AggregateSkillMap(
	idx UnitSkillIndex,
	attempts []AttemptRow,
	artifacts []ArtifactRow,
	progress []ProgressRow,
) SkillMap {
	// Latest attempt per unit_key|task_id (prefer higher ID, then CreatedAt).
	type key struct{ uk, tid string }
	latest := map[key]AttemptRow{}
	for _, a := range attempts {
		k := key{a.UnitKey, a.TaskID}
		prev, ok := latest[k]
		if !ok || a.ID > prev.ID || (a.ID == prev.ID && a.CreatedAt > prev.CreatedAt) {
			latest[k] = a
		}
	}

	// Per-skill counters.
	type counters struct {
		practiced int
		graded    int
		correct   int
	}
	bySkill := map[string]*counters{}
	for _, s := range AllSkills {
		bySkill[s] = &counters{}
	}

	// Reading skills from attempts.
	for _, a := range latest {
		taskMap := idx.TaskSkill[a.UnitKey]
		if taskMap == nil {
			continue
		}
		skill := taskMap[a.TaskID]
		if skill == "" {
			continue
		}
		c := bySkill[skill]
		if c == nil {
			continue
		}
		c.practiced++
		if a.Correct != nil {
			c.graded++
			if *a.Correct {
				c.correct++
			}
		}
	}

	// Writing skills from artifacts (latest body presence per unit).
	type writeState struct {
		hasDraft    bool
		hasRevision bool
		revOK       bool
	}
	writes := map[string]*writeState{}
	for _, art := range artifacts {
		ws := writes[art.UnitKey]
		if ws == nil {
			ws = &writeState{}
			writes[art.UnitKey] = ws
		}
		if art.Revision == 0 && strings.TrimSpace(art.Body) != "" {
			ws.hasDraft = true
		}
		if art.Revision == 1 && strings.TrimSpace(art.Body) != "" {
			ws.hasRevision = true
			ws.revOK = art.ChecklistAllTrue
		}
	}
	for unitKey, ws := range writes {
		if len(idx.WritingSkills[unitKey]) == 0 {
			continue
		}
		if ws.hasDraft {
			c := bySkill["writing.claim-reason"]
			c.practiced++
			c.graded++
			// Draft itself is progressive; count as correct if non-empty (practiced).
			c.correct++
		}
		if ws.hasRevision {
			c := bySkill["writing.revision"]
			c.practiced++
			c.graded++
			if ws.revOK {
				c.correct++
			}
		}
	}

	completed := map[string]bool{}
	for _, p := range progress {
		if p.Status == "completed" {
			completed[p.UnitKey] = true
		}
	}

	skills := make([]SkillStat, 0, len(AllSkills))
	weakSkills := []string{}
	for _, skill := range AllSkills {
		c := bySkill[skill]
		st := SkillStat{
			Skill:         skill,
			LabelJa:       SkillLabelJa[skill],
			Graded:        c.graded,
			Correct:       c.correct,
			Practiced:     c.practiced,
			UnitsTouching: append([]string{}, idx.SkillUnits[skill]...),
		}
		st.Status = skillStatus(c.practiced, c.graded, c.correct)
		if c.graded > 0 {
			acc := float64(c.correct) / float64(c.graded)
			st.Accuracy = &acc
		}
		if st.Status == "weak" || st.Status == "unseen" {
			weakSkills = append(weakSkills, skill)
		}
		skills = append(skills, st)
	}

	review := buildReviewQueue(idx, weakSkills, completed, 5)
	return SkillMap{Skills: skills, ReviewQueue: review}
}

func skillStatus(practiced, graded, correct int) string {
	if practiced == 0 {
		return "unseen"
	}
	if graded == 0 {
		// Ungraded-only (predict) or writing without grade signal.
		return "practiced"
	}
	acc := float64(correct) / float64(graded)
	if graded >= 2 && acc >= 0.8 {
		return "strong"
	}
	if acc < 0.6 {
		return "weak"
	}
	return "practiced"
}

func buildReviewQueue(idx UnitSkillIndex, weakSkills []string, completed map[string]bool, limit int) []ReviewItem {
	if limit <= 0 {
		return []ReviewItem{}
	}
	// Score incomplete units by how many weak skills they cover.
	type scored struct {
		unitKey string
		skills  []string
		score   int
	}
	var candidates []scored
	for unitKey, meta := range idx.UnitMeta {
		_ = meta
		if completed[unitKey] {
			continue
		}
		var hit []string
		for _, skill := range weakSkills {
			for _, uk := range idx.SkillUnits[skill] {
				if uk == unitKey {
					hit = append(hit, skill)
					break
				}
			}
		}
		if len(hit) == 0 {
			continue
		}
		// Prefer weak (had graded failures) over pure unseen: order weakSkills already groups them.
		candidates = append(candidates, scored{unitKey: unitKey, skills: hit, score: len(hit)})
	}
	sort.Slice(candidates, func(i, j int) bool {
		if candidates[i].score != candidates[j].score {
			return candidates[i].score > candidates[j].score
		}
		return candidates[i].unitKey < candidates[j].unitKey
	})
	out := make([]ReviewItem, 0, limit)
	for _, c := range candidates {
		if len(out) >= limit {
			break
		}
		meta := idx.UnitMeta[c.unitKey]
		out = append(out, ReviewItem{
			Stage:         meta.Stage,
			UnitID:        meta.UnitID,
			TitleJa:       meta.TitleJa,
			Genre:         meta.Genre,
			TargetSkills:  c.skills,
			UnitCompleted: false,
		})
	}
	if out == nil {
		out = []ReviewItem{}
	}
	return out
}

// ChecklistAllTrue parses a checklist JSON array of bools (or truthy values).
func ChecklistAllTrue(raw json.RawMessage) bool {
	if len(raw) == 0 {
		return false
	}
	var arr []any
	if err := json.Unmarshal(raw, &arr); err != nil || len(arr) == 0 {
		return false
	}
	for _, v := range arr {
		switch t := v.(type) {
		case bool:
			if !t {
				return false
			}
		default:
			return false
		}
	}
	return true
}
