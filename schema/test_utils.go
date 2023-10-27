package schema

import (
	"strings"
)

func visibleWhitespace(v string) string {
	v = strings.ReplaceAll(v, "\n", "\\n")
	v = strings.ReplaceAll(v, "\t", "\\t")
	v = strings.ReplaceAll(v, " ", "\\s")
	return v
}

func stripLines(a string) string {
	lines := strings.Split(a, "\n")
	for i := 0; i < len(lines); i++ {
		lines[i] = strings.TrimSpace(lines[i])
		if lines[i] == "" {
			lines = append(lines[:i], lines[i+1:]...)
			i--
		}
	}
	return strings.Join(lines, " ")
}

func removeDuplicateSpaces(a string) string {
	b := a
	for {
		a = strings.ReplaceAll(a, "  ", " ")
		if a == b {
			return b
		}
		b = a
	}
}

func normalizeParenthesis(a string) string {
	a = strings.ReplaceAll(a, "( ", "(")
	a = strings.ReplaceAll(a, " )", ")")
	a = strings.ReplaceAll(a, " (", "(")
	a = strings.ReplaceAll(a, ") ", ")")
	return a
}

func normalizeSemicolons(a string) string {
	a = strings.ReplaceAll(a, "; ", ";")
	a = strings.ReplaceAll(a, " ;", ";")
	a = strings.ReplaceAll(a, "\n;", ";\n")
	return a
}

func sqlStatementsAreEqual(a, b string) bool {
	if strings.TrimSpace(a) == strings.TrimSpace(b) {
		return true
	}
	a, b = stripLines(a), stripLines(b)
	if a == b {
		return true
	}
	a, b = removeDuplicateSpaces(a), removeDuplicateSpaces(b)
	if a == b {
		return true
	}
	a, b = normalizeParenthesis(a), normalizeParenthesis(b)
	if a == b {
		return true
	}
	a, b = normalizeSemicolons(a), normalizeSemicolons(b)
	return a == b
}
