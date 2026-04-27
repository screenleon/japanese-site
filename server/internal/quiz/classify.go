package quiz

import "strings"

// Classifier returns an error_class key for an incorrect answer. Returning
// the empty string means "no specific match — fall through to generic".
//
// Every grammar point that wants pointed feedback registers a Classifier in
// the classifiers map (see classify_n5.go, classify_n3.go, classify_n2.go).
// Grammar points without an entry produce error_class = "generic".
//
// When this map gets too long or rules become uniform string-matching, we
// can move the rules into the grammar JSON files (see ROADMAP.md → quiz
// classifier data-driven migration).
type Classifier func(expected, answer string) string

var classifiers = map[string]Classifier{}

func registerClassifier(grammarPoint string, c Classifier) {
	if _, exists := classifiers[grammarPoint]; exists {
		// Duplicate registration is a programming bug; panic at init time
		// rather than silently overwrite.
		panic("duplicate classifier for grammar point: " + grammarPoint)
	}
	classifiers[grammarPoint] = c
}

func classify(grammarPoint, expected, answer string) string {
	if c, ok := classifiers[grammarPoint]; ok {
		if cls := c(expected, answer); cls != "" {
			return cls
		}
	}
	return "generic"
}

// ── Helpers used by individual classifier functions ─────────────────────────

func hasSuffixAny(s string, suffixes ...string) bool {
	for _, suffix := range suffixes {
		if strings.HasSuffix(s, suffix) {
			return true
		}
	}
	return false
}

func containsAny(s string, substrs ...string) bool {
	for _, sub := range substrs {
		if strings.Contains(s, sub) {
			return true
		}
	}
	return false
}

// endsInDictionaryForm checks whether the last rune of s is one of the nine
// 五段/一段 dictionary-form tail kana. Used by te-form / masu-form classifiers.
func endsInDictionaryForm(s string) bool {
	r := []rune(s)
	if len(r) == 0 {
		return false
	}
	switch r[len(r)-1] {
	case 'う', 'く', 'ぐ', 'す', 'つ', 'ぬ', 'ぶ', 'む', 'る':
		return true
	}
	return false
}
