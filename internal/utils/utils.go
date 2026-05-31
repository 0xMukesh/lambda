package utils

func IsLetter(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || c == '_'
}

func IsWhitespace(c byte) bool {
	return c == '\n' || c == ' ' || c == '\t' || c == 0
}
