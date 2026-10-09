package db

import (
	"testing"
	"unicode"
)

// 说明。
// 说明。
// 说明。
func assertRetestReasonChinese(t *testing.T, name, value string) {
	t.Helper()
	if value == "" {
		t.Fatalf("%s: 文案不能为空", name)
	}
	hasHan := false
	for _, r := range value {
		if unicode.Is(unicode.Hangul, r) {
			t.Errorf("%s: 文案包含韩文: %q", name, value)
		}
		if unicode.Is(unicode.Han, r) {
			hasHan = true
		}
	}
	if !hasHan {
		t.Errorf("%s: 文案缺少中文: %q", name, value)
	}
}

// 说明。
// 说明。
// 说明。
// 说明。
func TestFindingRetestReasonsLocalized(t *testing.T) {
	assertRetestReasonChinese(t, "retestNoConclusionReason", retestNoConclusionReason)
	assertRetestReasonChinese(t, "retestServiceRestartReason", retestServiceRestartReason)
}
