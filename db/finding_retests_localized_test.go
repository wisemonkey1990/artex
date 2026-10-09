package db

import (
	"testing"
	"unicode"
)

// 说明。
// 说明。
// 说明。
func assertRetestReasonKorean(t *testing.T, name, s string) {
	t.Helper()
	if s == "" {
		t.Fatalf("%s: 测试文本 测试文本", name)
	}
	hasHangul := false
	for _, r := range s {
		if unicode.Is(unicode.Han, r) {
			t.Fatalf("%s: 测试文本 测试文本 测试文本 测试文本: %q", name, s)
		}
		if unicode.Is(unicode.Hangul, r) {
			hasHangul = true
		}
	}
	if !hasHangul {
		t.Fatalf("%s: 测试文本 测试文本: %q", name, s)
	}
}

// 说明。
// 说明。
// 说明。
// 说明。
func TestFindingRetestReasonsLocalized(t *testing.T) {
	assertRetestReasonKorean(t, "retestNoConclusionReason", retestNoConclusionReason)
	assertRetestReasonKorean(t, "retestServiceRestartReason", retestServiceRestartReason)
}
