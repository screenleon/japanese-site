// Package kokugo loads L1 国語 units from the corpus filesystem (JS-131/132).
// Units are not seeded into SQLite; only progress/attempts/artifacts are.
package kokugo

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// UnitSummary is a list-row projection of a KokugoUnit.
type UnitSummary struct {
	ID               string `json:"id"`
	Stage            string `json:"stage"`
	TitleJa          string `json:"title_ja"`
	Genre            string `json:"genre"`
	EstimatedMinutes int    `json:"estimated_minutes"`
	TaskCount        int    `json:"task_count"`
	HasArtifact      bool   `json:"has_artifact"`
}

// Loader reads units from corpus/kokugo/<stage>/<id>.json.
type Loader struct {
	// Root is the path to data/corpus/kokugo (not data/corpus).
	Root string
}

// List walks stage directories and returns unit summaries sorted by stage, id.
func (l Loader) List() ([]UnitSummary, error) {
	if l.Root == "" {
		return nil, fmt.Errorf("kokugo root empty")
	}
	entries, err := os.ReadDir(l.Root)
	if err != nil {
		if os.IsNotExist(err) {
			return []UnitSummary{}, nil
		}
		return nil, err
	}
	var out []UnitSummary
	for _, e := range entries {
		if !e.IsDir() || strings.HasPrefix(e.Name(), ".") {
			continue
		}
		stageDir := filepath.Join(l.Root, e.Name())
		files, err := os.ReadDir(stageDir)
		if err != nil {
			return nil, err
		}
		for _, f := range files {
			if f.IsDir() || !strings.HasSuffix(f.Name(), ".json") {
				continue
			}
			path := filepath.Join(stageDir, f.Name())
			sum, err := loadSummary(path)
			if err != nil {
				return nil, fmt.Errorf("%s: %w", path, err)
			}
			out = append(out, sum)
		}
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Stage != out[j].Stage {
			return out[i].Stage < out[j].Stage
		}
		return out[i].ID < out[j].ID
	})
	return out, nil
}

// Get returns the raw unit JSON object for stage/id.
func (l Loader) Get(stage, id string) (json.RawMessage, error) {
	if err := validateSegment(stage); err != nil {
		return nil, err
	}
	if err := validateSegment(id); err != nil {
		return nil, err
	}
	path := filepath.Join(l.Root, stage, id+".json")
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, ErrUnitNotFound
		}
		return nil, err
	}
	// Validate minimal shape + id/stage match.
	var probe struct {
		ID    string `json:"id"`
		Stage string `json:"stage"`
	}
	if err := json.Unmarshal(data, &probe); err != nil {
		return nil, fmt.Errorf("parse unit: %w", err)
	}
	if probe.ID != id || probe.Stage != stage {
		return nil, fmt.Errorf("unit id/stage mismatch: file=%s/%s json=%s/%s", stage, id, probe.Stage, probe.ID)
	}
	return json.RawMessage(data), nil
}

// GetMap unmarshals the unit into a generic map for grading.
func (l Loader) GetMap(stage, id string) (map[string]any, error) {
	raw, err := l.Get(stage, id)
	if err != nil {
		return nil, err
	}
	var m map[string]any
	if err := json.Unmarshal(raw, &m); err != nil {
		return nil, err
	}
	return m, nil
}

// ListMaps returns full unit maps for every L1 unit (JS-136 skill index).
// Order matches List() (stage, id).
func (l Loader) ListMaps() ([]map[string]any, error) {
	summaries, err := l.List()
	if err != nil {
		return nil, err
	}
	out := make([]map[string]any, 0, len(summaries))
	for _, s := range summaries {
		m, err := l.GetMap(s.Stage, s.ID)
		if err != nil {
			return nil, fmt.Errorf("%s/%s: %w", s.Stage, s.ID, err)
		}
		out = append(out, m)
	}
	return out, nil
}

func loadSummary(path string) (UnitSummary, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return UnitSummary{}, err
	}
	var u struct {
		ID               string `json:"id"`
		Stage            string `json:"stage"`
		TitleJa          string `json:"title_ja"`
		Genre            string `json:"genre"`
		EstimatedMinutes int    `json:"estimated_minutes"`
		Tasks            []any  `json:"tasks"`
		Artifact         any    `json:"artifact"`
	}
	if err := json.Unmarshal(data, &u); err != nil {
		return UnitSummary{}, err
	}
	if u.ID == "" || u.Stage == "" {
		return UnitSummary{}, fmt.Errorf("missing id/stage")
	}
	return UnitSummary{
		ID:               u.ID,
		Stage:            u.Stage,
		TitleJa:          u.TitleJa,
		Genre:            u.Genre,
		EstimatedMinutes: u.EstimatedMinutes,
		TaskCount:        len(u.Tasks),
		HasArtifact:      u.Artifact != nil,
	}, nil
}

func validateSegment(s string) error {
	if s == "" || strings.Contains(s, "..") || strings.ContainsAny(s, `/\`) {
		return fmt.Errorf("invalid path segment %q", s)
	}
	return nil
}

// ErrUnitNotFound means the unit JSON is missing.
var ErrUnitNotFound = fmt.Errorf("kokugo unit not found")
