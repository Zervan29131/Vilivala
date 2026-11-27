package model

import "testing"

func TestCountWordsAndReadingMinutes(t *testing.T) {
	if got := CountWords("你好 世界\nGo!"); got != 7 {
		t.Fatalf("expected 7 non-whitespace runes, got %d", got)
	}
	if got := ReadingMinutes(0); got != 1 {
		t.Fatalf("expected min 1 minute, got %d", got)
	}
	if got := ReadingMinutes(801); got != 3 {
		t.Fatalf("expected 3 minutes for 801 chars, got %d", got)
	}
}
