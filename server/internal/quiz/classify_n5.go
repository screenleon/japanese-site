package quiz

func init() {
	registerClassifier("te-form", classifyTeForm)
	registerClassifier("masu-form", classifyMasuForm)
	registerClassifier("copula-desu", classifyCopulaDesu)
	registerClassifier("wa-particle", classifyWaParticle)
	registerClassifier("ga-particle", classifyGaParticle)
}

func classifyTeForm(_, answer string) string {
	switch {
	case hasSuffixAny(answer, "ます", "ました"):
		return "used-masu"
	case endsInDictionaryForm(answer):
		return "used-dictionary"
	default:
		return "wrong-conjugation"
	}
}

func classifyMasuForm(_, answer string) string {
	if endsInDictionaryForm(answer) {
		return "used-plain"
	}
	return "wrong-tense"
}

func classifyCopulaDesu(_, answer string) string {
	if answer == "だ" {
		return "used-da-instead"
	}
	return "wrong-tense"
}

func classifyWaParticle(_, answer string) string {
	switch answer {
	case "が":
		return "used-ga"
	case "を":
		return "used-wo"
	}
	return ""
}

func classifyGaParticle(_, answer string) string {
	switch answer {
	case "は":
		return "used-wa"
	case "を":
		return "used-wo"
	}
	return ""
}
