// Package jmdict parses the jmdict-simplified JSON format and emits Vocab rows.
//
// Format reference: https://github.com/scriptin/jmdict-simplified
//
// Each top-level "word" entry has zero or more kanji forms, one or more kana
// readings, and one or more senses. We pick the first common kanji form (or
// the first kana reading if no kanji form exists) as the headword, the first
// kana reading as the reading, and join sense-level glosses for English.
package jmdict

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
	Version string `json:"version"`
	Words   []Word `json:"words"`
}

type Word struct {
	ID    string  `json:"id"`
	Kanji []Form  `json:"kanji"`
	Kana  []Form  `json:"kana"`
	Sense []Sense `json:"sense"`
}

type Form struct {
	Text   string `json:"text"`
	Common bool   `json:"common"`
}

type Sense struct {
	PartOfSpeech []string `json:"partOfSpeech"`
	Gloss        []Gloss  `json:"gloss"`
}

type Gloss struct {
	Lang string `json:"lang"`
	Text string `json:"text"`
}

// Vocab is the importer's row-shape view of a JMdict word.
type Vocab struct {
	JMDictID  string
	Headword  string
	Reading   string
	POS       string
	GlossEN   string
	Common    bool
}

// Parse reads a jmdict-simplified .json or .json.tgz file and returns parsed
// Vocab rows.
func Parse(path string) ([]Vocab, error) {
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
	dec := json.NewDecoder(r)
	if err := dec.Decode(&doc); err != nil {
		return nil, fmt.Errorf("decode jmdict: %w", err)
	}

	out := make([]Vocab, 0, len(doc.Words))
	for _, w := range doc.Words {
		v, ok := wordToVocab(w)
		if !ok {
			continue
		}
		out = append(out, v)
	}
	return out, nil
}

func wordToVocab(w Word) (Vocab, bool) {
	headword, common := pickHeadword(w)
	if headword == "" {
		return Vocab{}, false
	}
	reading := ""
	if len(w.Kana) > 0 {
		reading = w.Kana[0].Text
	}
	if reading == "" {
		return Vocab{}, false
	}
	pos := joinPOS(w.Sense)
	gloss := joinGloss(w.Sense, "eng")
	return Vocab{
		JMDictID: w.ID,
		Headword: headword,
		Reading:  reading,
		POS:      pos,
		GlossEN:  gloss,
		Common:   common,
	}, true
}

func pickHeadword(w Word) (string, bool) {
	for _, k := range w.Kanji {
		if k.Common {
			return k.Text, true
		}
	}
	if len(w.Kanji) > 0 {
		return w.Kanji[0].Text, false
	}
	if len(w.Kana) > 0 {
		return w.Kana[0].Text, w.Kana[0].Common
	}
	return "", false
}

func joinPOS(senses []Sense) string {
	seen := map[string]bool{}
	parts := []string{}
	for _, s := range senses {
		for _, p := range s.PartOfSpeech {
			if seen[p] {
				continue
			}
			seen[p] = true
			parts = append(parts, p)
		}
	}
	return strings.Join(parts, ",")
}

func joinGloss(senses []Sense, lang string) string {
	parts := []string{}
	for _, s := range senses {
		var senseParts []string
		for _, g := range s.Gloss {
			if g.Lang != "" && g.Lang != lang {
				continue
			}
			senseParts = append(senseParts, g.Text)
		}
		if len(senseParts) > 0 {
			parts = append(parts, strings.Join(senseParts, "; "))
		}
	}
	return strings.Join(parts, " | ")
}

// openMaybeTGZ wraps the reader with gzip + tar if the file ends in .tgz / .tar.gz.
// jmdict-simplified releases ship a single .json inside the tarball.
func openMaybeTGZ(f *os.File, path string) (io.Reader, error) {
	low := strings.ToLower(path)
	if !strings.HasSuffix(low, ".tgz") && !strings.HasSuffix(low, ".tar.gz") {
		return f, nil
	}
	gz, err := gzip.NewReader(f)
	if err != nil {
		return nil, fmt.Errorf("gzip open: %w", err)
	}
	tr := tar.NewReader(gz)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			return nil, fmt.Errorf("no .json file inside %s", path)
		}
		if err != nil {
			return nil, fmt.Errorf("tar read: %w", err)
		}
		if hdr.Typeflag != tar.TypeReg {
			continue
		}
		if strings.HasSuffix(strings.ToLower(hdr.Name), ".json") {
			return tr, nil
		}
	}
}
