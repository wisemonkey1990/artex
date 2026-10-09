package notify

import (
	"strings"
	"testing"
)

// 说明。
// 说明。
// 说明。
func hanFreeItems(n int) []Item {
	items := make([]Item, 0, n)
	for i := 0; i < n; i++ {
		items = append(items, Item{
			FindingID: int64(i + 1),
			Name:      "login-flaw",
			VulnClass: "SQLi",
			Severity:  "high",
			Summary:   "SQL injection via q param",
			Assets:    []string{"a.example.com"},
		})
	}
	return items
}

// 说明。
// 说明。
// 说明。
func TestMarkdownTitleLocalized(t *testing.T) {
	// 说明。
	got := markdownTitle(Message{Batch: true, Items: hanFreeItems(3)})
	assertKorean(t, "markdownTitle(batch)", got)
	if !strings.Contains(got, "测试文本 测试文本") || !strings.Contains(got, "测试文本 3测试文本") {
		t.Errorf("测试文本 测试文本 '测试文本 测试文本 · 测试文本 3测试文本' 测试文本 测试文本, 测试文本 测试文本 %q", got)
	}
	// 说明。
	empty := markdownTitle(Message{})
	assertKorean(t, "markdownTitle(empty)", empty)
	if empty != "测试文本 测试文本" {
		t.Errorf("测试文本 测试文本 测试文本 '测试文本 测试文本' 测试文本 测试文本, 测试文本 测试文本 %q", empty)
	}
}

// 说明。
// 说明。
// 说明。
func TestMarkdownBatchIntroLocalized(t *testing.T) {
	items := hanFreeItems(3)

	// 说明。
	win := markdownBatchIntro(Message{Batch: true, WindowMinutes: 30}, items, 3)
	if hasHan(win) {
		t.Errorf("测试文本 测试文本 测试文本 测试文本 测试文本: %q", win)
	}
	if !strings.Contains(win, "测试文本 30测试文本") || !strings.Contains(win, "新增漏洞 3 项") {
		t.Errorf("测试文本 测试文本 '测试文本 30测试文本 测试文本 测试文本 3测试文本' 测试文本 测试文本, 测试文本 测试文本 %q", win)
	}

	// 说明。
	noWin := markdownBatchIntro(Message{Batch: true}, items, 3)
	if hasHan(noWin) {
		t.Errorf("测试文本 测试文本 测试文本 测试文本: %q", noWin)
	}
	if !strings.Contains(noWin, "新增漏洞 3 项") || strings.Contains(noWin, "测试文本") {
		t.Errorf("测试文本 测试文本 测试文本 '测试文本 测试文本 3测试文本'(测试文本 测试文本 测试文本)测试文本 测试文本, 测试文本 测试文本 %q", noWin)
	}

	// 说明。
	trunc := markdownBatchIntro(Message{Batch: true, WindowMinutes: 30}, items, 5)
	if hasHan(trunc) {
		t.Errorf("测试文本 测试文本 测试文本 测试文本 测试文本: %q", trunc)
	}
	for _, want := range []string{"测试文本 3测试文本", "测试文本 2测试文本", "测试文本 测试文本"} {
		if !strings.Contains(trunc, want) {
			t.Errorf("测试文本 测试文本 %q 测试文本 测试文本 测试文本, 测试文本 测试文本 %q", want, trunc)
		}
	}
}

// 说明。
// 说明。
func TestWriteItemLabelsLocalized(t *testing.T) {
	it := Item{
		Name:       "login-flaw",
		VulnClass:  "SQLi",
		Severity:   "high",
		Summary:    "SQL injection via q param",
		Assets:     []string{"a.example.com"},
		DetailURL:  "https://platform.example/finding/1",
		FromStatus: "pending",
		ToStatus:   "fixed",
	}
	var b strings.Builder
	writeItem(&b, it, "", true)
	got := b.String()

	assertKorean(t, "writeItem(single)", got)
	for _, want := range []string{
		"**测试文本 测试文本**: 测试文本 测试文本 → 测试文本",
		"**测试文本**: SQLi",
		"**测试文本**: a.example.com",
		"**测试文本**: SQL injection via q param",
		"[测试文本 测试文本](https://platform.example/finding/1)",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("测试文本 测试文本 %q 测试文本 测试文本 测试文本:\n%s", want, got)
		}
	}
}

// 说明。
// 说明。
func TestMarkdownBodyFooterLocalized(t *testing.T) {
	m := Message{Batch: true, HomeURL: "https://platform.example", Items: hanFreeItems(2)}
	body, kept := markdownBody(m, 0)
	if kept != 2 {
		t.Fatalf("测试文本 测试文本(0)测试文本 2测试文本 测试文本 测试文本 测试文本, 测试文本 测试文本 %d", kept)
	}
	if !strings.Contains(body, "[测试文本 测试文本 测试文本](https://platform.example)") {
		t.Errorf("测试文本 测试文本 '测试文本 测试文本 测试文本' 测试文本 测试文本 测试文本:\n%s", body)
	}
}
