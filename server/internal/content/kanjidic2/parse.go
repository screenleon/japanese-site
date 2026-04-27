// Package kanjidic2 parses kanjidic2-simplified JSON and emits Kanji rows.
// Format: https://github.com/scriptin/jmdict-simplified
package kanjidic2

import (
	"archive/tar"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
)

type Document struct {
	Version    string      `json:"version"`
	Characters []Character `json:"characters"`
}

type Character struct {
	Literal     string         `json:"literal"`
	Misc        Misc           `json:"misc"`
	ReadingMean ReadingMeaning `json:"readingMeaning"`
}

type Misc struct {
	Grade        *int   `json:"grade"`
	StrokeCounts []int  `json:"strokeCounts"`
	JLPTLevel    *int   `json:"jlptLevel"`   // old-JLPT 1..4
	Frequency    *int   `json:"frequency"`
}

type ReadingMeaning struct {
	Groups []ReadingMeaningGroup `json:"groups"`
}

type ReadingMeaningGroup struct {
	Readings []Reading `json:"readings"`
	Meanings []Meaning `json:"meanings"`
}

type Reading struct {
	Type  string `json:"type"`  // ja_on, ja_kun, ...
	Value string `json:"value"`
}

type Meaning struct {
	Lang  string `json:"lang"`
	Value string `json:"value"`
}

type Kanji struct {
	Character    string
	Onyomi       string // comma-joined ja_on
	Kunyomi      string // comma-joined ja_kun
	MeaningEN    string
	JLPTLevelOld int // 0 if absent
	Grade        int // 0 if absent
	StrokeCount  int // 0 if absent
}

func Parse(path string) ([]Kanji, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	r, err := openMaybeTGZ(f, path)
	if err != nil {
		return nil, err
	}

	var doc Document
	if err := json.NewDecoder(r).Decode(&doc); err != nil {
		return nil, fmt.Errorf("decode: %w", err)
	}

	out := make([]Kanji, 0, len(doc.Characters))
	for _, c := range doc.Characters {
		k := Kanji{Character: c.Literal}
		if c.Misc.JLPTLevel != nil {
			k.JLPTLevelOld = *c.Misc.JLPTLevel
		}
		if c.Misc.Grade != nil {
			k.Grade = *c.Misc.Grade
		}
		if len(c.Misc.StrokeCounts) > 0 {
			k.StrokeCount = c.Misc.StrokeCounts[0]
		}
		var on, kun, mean []string
		for _, g := range c.ReadingMean.Groups {
			for _, r := range g.Readings {
				switch r.Type {
				case "ja_on":
					on = append(on, r.Value)
				case "ja_kun":
					kun = append(kun, r.Value)
				}
			}
			for _, m := range g.Meanings {
				if m.Lang == "" || m.Lang == "en" {
					mean = append(mean, m.Value)
				}
			}
		}
		k.Onyomi = strings.Join(on, ",")
		k.Kunyomi = strings.Join(kun, ",")
		k.MeaningEN = strings.Join(mean, "; ")
		out = append(out, k)
	}
	return out, nil
}

// MapOldJLPT translates the legacy 4-level JLPT scale into the modern 5-level
// scale (best-effort). Returns "" when no level is set.
//   old 1 → N1
//   old 2 → N2
//   old 3 → N4   (old-3 split into modern N3 + N4; N4 is the safer floor)
//   old 4 → N5
// Use this only for a starting tag; the content-validator may refine later.
func MapOldJLPT(old int) string {
	switch old {
	case 4:
		return "N5"
	case 3:
		return "N4"
	case 2:
		return "N2"
	case 1:
		return "N1"
	default:
		return ""
	}
}

func openMaybeTGZ(f *os.File, path string) (io.Reader, error) {
	low := strings.ToLower(path)
	if !strings.HasSuffix(low, ".tgz") && !strings.HasSuffix(low, ".tar.gz") {
		return f, nil
	}
	gz, err := gzip.NewReader(f)
	if err != nil {
		return nil, fmt.Errorf("gzip: %w", err)
	}
	tr := tar.NewReader(gz)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			return nil, fmt.Errorf("no .json file inside %s", path)
		}
		if err != nil {
			return nil, err
		}
		if hdr.Typeflag == tar.TypeReg && strings.HasSuffix(strings.ToLower(hdr.Name), ".json") {
			return tr, nil
		}
	}
}
