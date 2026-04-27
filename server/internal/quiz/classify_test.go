package quiz

import "testing"

// TestClassify covers a few representative cases per grammar family. Adding
// a new grammar point's classifier should usually be paired with a couple of
// rows here so future refactors don't silently regress the rules.
func TestClassify(t *testing.T) {
	cases := []struct {
		grammarPoint string
		answer       string
		want         string
	}{
		// N5
		{"te-form", "食べます", "used-masu"},
		{"te-form", "食べる", "used-dictionary"},
		{"te-form", "食べた", "wrong-conjugation"},
		{"masu-form", "食べる", "used-plain"},
		{"copula-desu", "だ", "used-da-instead"},
		{"wa-particle", "が", "used-ga"},
		{"ga-particle", "は", "used-wa"},
		// N3
		{"conditional-tara", "降れば", "used-ba"},
		{"conditional-tara", "降るなら", "used-nara"},
		{"conditional-tara", "降ると", "used-to"},
		{"conditional-tara", "降ります", "wrong-conjugation"},
		{"conditional-ba", "降ったら", "used-tara"},
		{"conditional-nara", "降ったら", "used-tara"},
		{"conditional-nara", "降れば", "used-ba"},
		{"concession-temo", "降ったら", "used-tara"},
		{"concession-temo", "降ったのに", "used-noni"},
		{"contrast-noni", "降っても", "used-temo"},
		{"contrast-noni", "が", "used-kedo"},
		// N2
		{"wakeda", "はずだ", "used-hazu"},
		{"monoda", "ことだ", "used-koto"},
		{"dokoroka", "ばかりか", "used-bakarika"},
		{"nikanshite", "にとって", "used-nitotte"},
		{"tsutsu", "聞きながら", "used-nagara"},
		{"tsutsu", "聞いている", "used-teiru"},
		// Unknown grammar point or unmatched answer falls through to generic.
		{"unknown-point", "anything", "generic"},
		{"wakeda", "neither hazu nor koto", "generic"},
	}
	for _, c := range cases {
		t.Run(c.grammarPoint+"/"+c.answer, func(t *testing.T) {
			got := classify(c.grammarPoint, "", c.answer)
			if got != c.want {
				t.Errorf("classify(%q, %q) = %q, want %q", c.grammarPoint, c.answer, got, c.want)
			}
		})
	}
}
