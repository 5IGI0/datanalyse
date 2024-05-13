package main

func SanitizeText() {

}

func OnlyAlphaNum(text string) string {
	var ret string
	for _, c := range text {
		if (c >= 'a' && c <= 'z') ||
			(c >= 'A' && c <= 'Z') ||
			(c >= '0' && c <= '9') {
			ret += string(c)
		}
	}
	return ret
}
