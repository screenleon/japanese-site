package quizrule

import "testing"

func TestClassify(t *testing.T) {
	cases := []struct {
		name   string
		rules  []Rule
		answer string
		want   string
	}{
		{
			name: "suffix",
			rules: []Rule{
				{ErrorClass: "used-ba", IfAnswerSuffixAny: []string{"ければ", "れば", "ば"}},
			},
			answer: "降れば",
			want:   "used-ba",
		},
		{
			name: "contains",
			rules: []Rule{
				{ErrorClass: "used-temo", IfAnswerContainsAny: []string{"ても", "でも"}},
			},
			answer: "降っても",
			want:   "used-temo",
		},
		{
			name: "equals",
			rules: []Rule{
				{ErrorClass: "used-ga", IfAnswerEqualsAny: []string{"が"}},
			},
			answer: "が",
			want:   "used-ga",
		},
		{
			name: "dictionary form",
			rules: []Rule{
				{ErrorClass: "used-plain", IfAnswerDictionaryForm: true},
			},
			answer: "食べる",
			want:   "used-plain",
		},
		{
			name: "not contains",
			rules: []Rule{
				{ErrorClass: "wrong-conjugation", IfAnswerNotContainsAny: []string{"ても", "でも"}},
			},
			answer: "降ったら",
			want:   "wrong-conjugation",
		},
		{
			name: "default",
			rules: []Rule{
				{ErrorClass: "wrong-tense", Default: true},
			},
			answer: "anything",
			want:   "wrong-tense",
		},
		{
			name:   "generic fallback",
			rules:  nil,
			answer: "anything",
			want:   "generic",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := Classify(tc.rules, "", tc.answer)
			if got != tc.want {
				t.Fatalf("Classify() = %q, want %q", got, tc.want)
			}
		})
	}
}

func TestValidateRulesRejectsMalformedRules(t *testing.T) {
	cases := []struct {
		name  string
		rules []Rule
	}{
		{
			name:  "missing error class",
			rules: []Rule{{IfAnswerEqualsAny: []string{"が"}}},
		},
		{
			name:  "missing predicate",
			rules: []Rule{{ErrorClass: "used-ga"}},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if err := ValidateRules(tc.rules); err == nil {
				t.Fatal("expected validation error, got nil")
			}
		})
	}
}

func TestValidateRulesAcceptsDefault(t *testing.T) {
	err := ValidateRules([]Rule{{Default: true, ErrorClass: "generic"}})
	if err != nil {
		t.Fatalf("ValidateRules(default): %v", err)
	}
}
