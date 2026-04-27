// Package tatoeba parses Tatoeba's per-language sentence exports.
// Source: https://downloads.tatoeba.org/exports/per_language/
//
// jpn_sentences.tsv format (tab-separated):
//   <id>\t<lang>\t<text>
package tatoeba

import (
	"bufio"
	"compress/bzip2"
	"fmt"
	"io"
	"os"
	"strings"
)

type Sentence struct {
	TatoebaID int64
	TextJA    string
}

func Parse(path string) ([]Sentence, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var r io.Reader = f
	if strings.HasSuffix(strings.ToLower(path), ".bz2") {
		r = bzip2.NewReader(f)
	}

	scanner := bufio.NewScanner(r)
	scanner.Buffer(make([]byte, 1024*1024), 4*1024*1024)

	var out []Sentence
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, "\t", 3)
		if len(parts) < 3 {
			continue
		}
		var id int64
		if _, err := fmt.Sscanf(parts[0], "%d", &id); err != nil {
			continue
		}
		if parts[1] != "jpn" {
			continue
		}
		text := strings.TrimSpace(parts[2])
		if text == "" {
			continue
		}
		out = append(out, Sentence{TatoebaID: id, TextJA: text})
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scan: %w", err)
	}
	return out, nil
}
