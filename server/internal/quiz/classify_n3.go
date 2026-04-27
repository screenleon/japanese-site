package quiz

func init() {
	registerClassifier("conditional-tara", classifyTara)
	registerClassifier("conditional-ba", classifyBa)
	registerClassifier("conditional-nara", classifyNara)
	registerClassifier("concession-temo", classifyTemo)
	registerClassifier("contrast-noni", classifyNoni)
}

// 〜たら expected; user typed another conditional or wrong conjugation.
func classifyTara(_, answer string) string {
	switch {
	case hasSuffixAny(answer, "ければ", "れば", "ば"):
		return "used-ba"
	case hasSuffixAny(answer, "なら"):
		return "used-nara"
	case hasSuffixAny(answer, "と"):
		return "used-to"
	case !containsAny(answer, "たら", "だら"):
		return "wrong-conjugation"
	}
	return ""
}

func classifyBa(_, answer string) string {
	switch {
	case hasSuffixAny(answer, "たら", "だら"):
		return "used-tara"
	case !containsAny(answer, "ば"):
		return "wrong-conjugation"
	}
	return ""
}

func classifyNara(_, answer string) string {
	switch {
	case hasSuffixAny(answer, "たら", "だら"):
		return "used-tara"
	case hasSuffixAny(answer, "ければ", "れば", "ば"):
		return "used-ba"
	}
	return ""
}

func classifyTemo(_, answer string) string {
	switch {
	case hasSuffixAny(answer, "たら", "だら"):
		return "used-tara"
	case hasSuffixAny(answer, "のに"):
		return "used-noni"
	case !containsAny(answer, "ても", "でも"):
		return "wrong-conjugation"
	}
	return ""
}

func classifyNoni(_, answer string) string {
	switch {
	case containsAny(answer, "ても", "でも"):
		return "used-temo"
	case containsAny(answer, "けど") || answer == "が":
		return "used-kedo"
	}
	return ""
}
