package utils

import "strings"

func Ternary[T any](condition bool, a, b T) T {
	if condition {
		return a
	}
	return b
}

func PadEmptyLines(s string, height int) string {
	lines := strings.Split(s, "\n")
	for len(lines) < height {
		lines = append(lines, "")
	}
	return strings.Join(lines, "\n")
}

type ChainReplacer struct {
	str string
}

func NewReplacer(s string) *ChainReplacer {
	return &ChainReplacer{str: s}
}

func (cr *ChainReplacer) Replace(old, new string) *ChainReplacer {
	cr.str = strings.ReplaceAll(cr.str, old, new)
	return cr
}

func (cr *ChainReplacer) String() string {
	return cr.str
}
