package external

import (
	"encoding/json"
	"fmt"
	"os"
)

const SchemaVersion = 1

type Manifest struct {
	SchemaVersion int       `json:"schema_version"`
	Datasets      []Dataset `json:"datasets"`
}

type Dataset struct {
	Name     string `json:"name"`
	Version  string `json:"version"`
	URL      string `json:"url"`
	SHA256   string `json:"sha256"`
	License  string `json:"license"`
	SourceID string `json:"source_id"`
	Level    string `json:"level,omitempty"` // optional: JLPT level for overlay datasets
	Note     string `json:"_note,omitempty"`
}

// FilterBySourceID returns datasets whose source_id matches.
func (m Manifest) FilterBySourceID(id string) []Dataset {
	var out []Dataset
	for _, d := range m.Datasets {
		if d.SourceID == id {
			out = append(out, d)
		}
	}
	return out
}

func (d Dataset) Validate() error {
	if d.Name == "" {
		return fmt.Errorf("dataset.name is required")
	}
	if d.URL == "" {
		return fmt.Errorf("dataset %q: url is required", d.Name)
	}
	if d.License == "" {
		return fmt.Errorf("dataset %q: license is required", d.Name)
	}
	if d.SourceID == "" {
		return fmt.Errorf("dataset %q: source_id is required", d.Name)
	}
	if d.Version == "" {
		return fmt.Errorf("dataset %q: version is required", d.Name)
	}
	return nil
}

func Load(path string) (Manifest, error) {
	body, err := os.ReadFile(path)
	if err != nil {
		return Manifest{}, fmt.Errorf("read %s: %w", path, err)
	}
	var m Manifest
	if err := json.Unmarshal(body, &m); err != nil {
		return Manifest{}, fmt.Errorf("parse %s: %w", path, err)
	}
	if m.SchemaVersion != SchemaVersion {
		return Manifest{}, fmt.Errorf("unsupported schema_version %d (want %d)", m.SchemaVersion, SchemaVersion)
	}
	for _, d := range m.Datasets {
		if err := d.Validate(); err != nil {
			return Manifest{}, err
		}
	}
	return m, nil
}

func (m Manifest) Find(name string) (Dataset, bool) {
	for _, d := range m.Datasets {
		if d.Name == name {
			return d, true
		}
	}
	return Dataset{}, false
}

func Save(path string, m Manifest) error {
	body, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return err
	}
	body = append(body, '\n')
	return os.WriteFile(path, body, 0o644)
}
