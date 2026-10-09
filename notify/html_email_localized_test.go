package notify

import (
	"strings"
	"testing"
)

// 说明。
// 说明。
// 说明。
//
// 说明。
// 说明。
// 说明。

// 说明。
// 说明。
func TestHTMLItemLabelsLocalized(t *testing.T) {
	m := Message{
		HomeURL: "https://example.com/panel",
		Items: []Item{{
			Name:       "sqli-login", // Title() = Name
			VulnClass:  "injection",
			Severity:   "high",
			Summary:    "login form is injectable",
			Assets:     []string{"host-a.example.com"},
			DetailURL:  "https://example.com/f/1",
			FromStatus: "pending",
			ToStatus:   "fixed",
		}},
	}
	out := htmlBody(m, 0)
	if hasHan(out) {
		t.Errorf("测试文本 测试文本 测试文本 测试文本 测试文本 测试文本:\n%s", out)
	}
	for _, want := range []string{
		"测试文本 测试文本", "测试文本", "测试文本", "概述", "查看详情", "测试文本 测试文本 测试文本",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("测试文本 测试文本 %q 测试文本 测试文本:\n%s", want, out)
		}
	}
}

// 说明。
// 说明。
func TestHTMLBatchIntroLocalized(t *testing.T) {
	items := []Item{{Name: "a", Severity: "high"}, {Name: "b", Severity: "low"}}

	withWindow := htmlBatchIntro(Message{Batch: true, WindowMinutes: 30, Items: items})
	assertKorean(t, "htmlBatchIntro(测试文本)", withWindow)
	for _, want := range []string{"测试文本 30测试文本", "测试文本 测试文本", "2测试文本"} {
		if !strings.Contains(withWindow, want) {
			t.Errorf("测试文本 测试文本 %q 测试文本 测试文本: %q", want, withWindow)
		}
	}

	noWindow := htmlBatchIntro(Message{Batch: true, WindowMinutes: 0, Items: items})
	assertKorean(t, "htmlBatchIntro(测试文本 测试文本)", noWindow)
	if !strings.Contains(noWindow, "测试文本 测试文本 2测试文本") {
		t.Errorf("测试文本 测试文本 测试文本 %q 测试文本 测试文本 测试文本: %q", "测试文本 测试文本 2测试文本", noWindow)
	}
	if strings.Contains(noWindow, "测试文本") {
		t.Errorf("测试文本 测试文本 '测试文本' 测试文本 测试文本: %q", noWindow)
	}
}

// 说明。
// 说明。
// 说明。
func TestEmailValidateLocalized(t *testing.T) {
	cases := []struct {
		name   string
		cfg    map[string]any
		substr string
	}{
		{"测试文本 测试文本 测试文本", map[string]any{"port": float64(25), "from": "a@b.c", "to": []any{"d@e.f"}}, "SMTP"},
		{"测试文本 测试文本 测试文本", map[string]any{"host": "h"}, "测试文本"},
		{"测试文本 测试文本", map[string]any{"host": "h", "port": float64(25), "to": []any{"d@e.f"}}, "测试文本"},
		{"测试文本 测试文本", map[string]any{"host": "h", "port": float64(25), "from": "a@b.c"}, "测试文本"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := (emailChannel{}).Validate(tc.cfg)
			if err == nil {
				t.Fatalf("测试文本 测试文本 测试文本: %v", tc.cfg)
			}
			assertKorean(t, "Validate("+tc.name+")", err.Error())
			if !strings.Contains(err.Error(), tc.substr) {
				t.Errorf("测试文本 测试文本 %q 测试文本 测试文本 测试文本, 测试文本 测试文本 %q", tc.substr, err.Error())
			}
		})
	}
}
