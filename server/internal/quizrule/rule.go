package quizrule

import "strings"

// Rule is one ordered classifier rule loaded from grammar corpus JSON.
// The first matching rule wins. A rule with Default=true always matches.
type Rule struct {
	ErrorClass             string   `json:"error_class"`
	IfAnswerEqualsAny      []string `json:"if_answer_equals_any,omitempty"`
	IfAnswerSuffixAny      []string `json:"if_answer_suffix_any,omitempty"`
	IfAnswerContainsAny    []string `json:"if_answer_contains_any,omitempty"`
	IfAnswerNotContainsAny []string `json:"if_answer_not_contains_any,omitempty"`
	IfAnswerDictionaryForm bool     `json:"if_answer_dictionary_form,omitempty"`
	Default                bool     `json:"default,omitempty"`
}

// Classify returns the first matching error_class, or "generic" when no
// rule matches. expected is reserved for future rule shapes.
func Classify(rules []Rule, _, answer string) string {
	for _, rule := range rules {
		if rule.ErrorClass == "" {
			continue
		}
		if matches(rule, answer) {
			return rule.ErrorClass
		}
	}
	return "generic"
}

func matches(rule Rule, answer string) bool {
	if rule.Default {
		return true
	}
	if len(rule.IfAnswerEqualsAny) > 0 && !equalsAny(answer, rule.IfAnswerEqualsAny) {
		return false
	}
	if len(rule.IfAnswerSuffixAny) > 0 && !hasSuffixAny(answer, rule.IfAnswerSuffixAny) {
		return false
	}
	if len(rule.IfAnswerContainsAny) > 0 && !containsAny(answer, rule.IfAnswerContainsAny) {
		return false
	}
	if len(rule.IfAnswerNotContainsAny) > 0 && containsAny(answer, rule.IfAnswerNotContainsAny) {
		return false
	}
	if rule.IfAnswerDictionaryForm && !endsInDictionaryForm(answer) {
		return false
	}
	return true
}

func equalsAny(s string, values []string) bool {
	for _, value := range values {
		if s == value {
			return true
		}
	}
	return false
}

func hasSuffixAny(s string, suffixes []string) bool {
	for _, suffix := range suffixes {
		if strings.HasSuffix(s, suffix) {
			return true
		}
	}
	return false
}

func containsAny(s string, substrs []string) bool {
	for _, sub := range substrs {
		if strings.Contains(s, sub) {
			return true
		}
	}
	return false
}

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
