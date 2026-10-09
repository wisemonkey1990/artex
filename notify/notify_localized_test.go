package notify

import (
	"strings"
	"testing"
)

// 说明。
// 说明。
func hasHan(s string) bool {
	for _, r := range s {
		if r >= 0x4E00 && r <= 0x9FFF {
			return true
		}
	}
	return false
}

// 说明。
func hasHangul(s string) bool {
	for _, r := range s {
		if r >= 0xAC00 && r <= 0xD7A3 {
			return true
		}
	}
	return false
}

// 说明。
func assertKorean(t *testing.T, where, got string) {
	t.Helper()
	if hasHan(got) {
		t.Errorf("%s: 测试文本 测试文本 测试文本 测试文本: %q", where, got)
	}
	if !hasHangul(got) {
		t.Errorf("%s: 测试文本 测试文本: %q", where, got)
	}
}

// 说明。
// 说明。
// 说明。
func TestSeverityLabelLocalized(t *testing.T) {
	// 说明。
	want := map[string]string{
		"critical": "测试文本",
		"high":     "高危",
		"medium":   "测试文本",
		"low":      "测试文本",
	}
	for sev, label := range want {
		got := SeverityLabel(sev)
		assertKorean(t, "SeverityLabel("+sev+")", got)
		if !strings.Contains(got, label) {
			t.Errorf("SeverityLabel(%q)=%q, %q 测试文本 测试文本 测试文本", sev, got, label)
		}
	}
	// 说明。
	if got := SeverityLabel("made_up"); got != "made_up" {
		t.Errorf("测试文本 测试文本 测试文本 测试文本 测试文本 测试文本 测试文本, 测试文本 测试文本 %q", got)
	}
}

// 说明。
// 说明。
func TestStatusLabelLocalized(t *testing.T) {
	want := map[string]string{
		"pending":        "测试文本 测试文本",
		"in_progress":    "测试文本 测试文本",
		"confirmed":      "测试文本",
		"resolved":       "测试文本",
		"fixed":          "测试文本",
		"false_positive": "测试文本",
		"ignored":        "测试文本",
		"duplicate":      "测试文本",
		"risk_accepted":  "测试文本 测试文本",
	}
	for status, label := range want {
		got := StatusLabel(status)
		assertKorean(t, "StatusLabel("+status+")", got)
		if got != label {
			t.Errorf("StatusLabel(%q)=%q, %q 测试文本 测试文本", status, got, label)
		}
	}
	// 说明。
	if got := StatusLabel("weird_status"); got != "weird_status" {
		t.Errorf("测试文本 测试文本 测试文本 测试文本 测试文本 测试文本 测试文本, 测试文本 测试文本 %q", got)
	}
}

// 说明。
// 说明。
func TestItemTitlePlaceholderLocalized(t *testing.T) {
	got := Item{}.Title()
	assertKorean(t, "Item{}.Title()", got)
	// 说明。
	if n := (Item{Name: "测试文本 SQLi"}).Title(); n != "测试文本 SQLi" {
		t.Errorf("测试文本 测试文本 测试文本 测试文本, 测试文本 测试文本 %q", n)
	}
	if v := (Item{VulnClass: "XSS"}).Title(); v != "XSS" {
		t.Errorf("测试文本 测试文本 测试文本 测试文本, 测试文本 测试文本 %q", v)
	}
}

// 说明。
// 说明。
func TestAssetLineLocalized(t *testing.T) {
	// 说明。
	if got := assetLine([]string{"a.example.com", "b.example.com"}, 3); got != "a.example.com, b.example.com" {
		t.Errorf("测试文本 ASCII 测试文本 测试文本, 测试文本 测试文本 %q", got)
	}
	// 说明。
	got := assetLine([]string{"a", "b", "c", "d", "e"}, 2)
	if hasHan(got) {
		t.Errorf("测试文本 测试文本 测试文本 测试文本 测试文本: %q", got)
	}
	if !strings.Contains(got, "测试文本 5测试文本") {
		t.Errorf("测试文本 测试文本 5 测试文本 '测试文本 5测试文本'测试文本 测试文本 测试文本, 测试文本 测试文本 %q", got)
	}
}
