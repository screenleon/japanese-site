package quizrule

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestCorpusClassifierRules(t *testing.T) {
	cases := []struct {
		level        string
		grammarPoint string
		answer       string
		want         string
	}{
		// N5
		{"N5", "te-form", "食べます", "used-masu"},
		{"N5", "te-form", "食べる", "used-dictionary"},
		{"N5", "te-form", "食べた", "wrong-conjugation"},
		{"N5", "masu-form", "食べる", "used-plain"},
		{"N5", "copula-desu", "だ", "used-da-instead"},
		{"N5", "wa-particle", "が", "used-ga"},
		{"N5", "ga-particle", "は", "used-wa"},
		// N3
		{"N3", "conditional-tara", "降れば", "used-ba"},
		{"N3", "conditional-tara", "降るなら", "used-nara"},
		{"N3", "conditional-tara", "降ると", "used-to"},
		{"N3", "conditional-tara", "降ります", "wrong-conjugation"},
		{"N3", "conditional-ba", "降ったら", "used-tara"},
		{"N3", "conditional-nara", "降ったら", "used-tara"},
		{"N3", "conditional-nara", "降れば", "used-ba"},
		{"N3", "concession-temo", "降ったら", "used-tara"},
		{"N3", "concession-temo", "降ったのに", "used-noni"},
		{"N3", "contrast-noni", "降っても", "used-temo"},
		{"N3", "contrast-noni", "が", "used-kedo"},
		// N2
		{"N2", "wakeda", "はずだ", "used-hazu"},
		{"N2", "monoda", "ことだ", "used-koto"},
		{"N2", "dokoroka", "ばかりか", "used-bakarika"},
		{"N2", "nikanshite", "にとって", "used-nitotte"},
		{"N2", "tsutsu", "聞きながら", "used-nagara"},
		{"N2", "tsutsu", "聞いている", "used-teiru"},
		// Unknown/unmatched falls through.
		{"N2", "wakeda", "neither hazu nor koto", "generic"},
	}

	for _, tc := range cases {
		t.Run(tc.grammarPoint+"/"+tc.answer, func(t *testing.T) {
			rules := loadCorpusRules(t, tc.level, tc.grammarPoint)
			got := Classify(rules, "", tc.answer)
			if got != tc.want {
				t.Fatalf("Classify(%q, %q) = %q, want %q", tc.grammarPoint, tc.answer, got, tc.want)
			}
		})
	}
}

func loadCorpusRules(t *testing.T, level, slug string) []Rule {
	t.Helper()
	path := filepath.Join("..", "..", "data", "corpus", "grammar", level, slug+".json")
	body, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	var gp struct {
		ClassifierRules []Rule `json:"classifier_rules"`
	}
	if err := json.Unmarshal(body, &gp); err != nil {
		t.Fatalf("decode %s: %v", path, err)
	}
	return gp.ClassifierRules
}
