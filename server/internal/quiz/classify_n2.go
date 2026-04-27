package quiz

func init() {
	registerClassifier("wakeda", classifyWakeda)
	registerClassifier("monoda", classifyMonoda)
	registerClassifier("dokoroka", classifyDokoroka)
	registerClassifier("nikanshite", classifyNikanshite)
	registerClassifier("tsutsu", classifyTsutsu)
}

func classifyWakeda(_, answer string) string {
	switch {
	case containsAny(answer, "はず"):
		return "used-hazu"
	case containsAny(answer, "こと"):
		return "used-koto"
	}
	return ""
}

func classifyMonoda(_, answer string) string {
	switch {
	case containsAny(answer, "こと"):
		return "used-koto"
	case containsAny(answer, "べき"):
		return "used-beki"
	}
	return ""
}

func classifyDokoroka(_, answer string) string {
	switch {
	case containsAny(answer, "ばかり"):
		return "used-bakarika"
	case containsAny(answer, "くらい", "ぐらい"):
		return "used-kurai"
	}
	return ""
}

func classifyNikanshite(_, answer string) string {
	switch {
	case containsAny(answer, "とって"):
		return "used-nitotte"
	case containsAny(answer, "対して"):
		return "used-nitaishite"
	}
	return ""
}

func classifyTsutsu(_, answer string) string {
	switch {
	case containsAny(answer, "ながら"):
		return "used-nagara"
	case containsAny(answer, "ている", "ています"):
		return "used-teiru"
	case !containsAny(answer, "つつ"):
		return "wrong-conjugation"
	}
	return ""
}
