package notify

import (
	"strings"
	"testing"
)

// hanFreeItems 는 CJK 한자가 없는 항목(ASCII 제목·유형·요약)만 만든다.
// 라벨이 중국어로 되돌아가면 렌더 결과 전체에 한자가 생기므로, 데이터가
// 한자 0 일 때만 "출력 전체에 한자 0" 단언이 라벨 회귀를 정확히 잡아낸다.
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

// TestMarkdownTitleLocalized 는 markdown 계열 채널(dingtalk·wecom)과 제목을
// 공유하는 feishu·html·telegram·webhook 이 함께 쓰는 markdownTitle 이 한국어인지
// 검사한다. 하나라도 중국어로 되돌아가면 이 다섯 채널의 제목이 전부 혼재된다.
func TestMarkdownTitleLocalized(t *testing.T) {
	// PR#1 把 markdownTitle(markdown.go) 从韩语改成中文，这里反向断言。
	// 다건(汇总): "漏洞摘要 · 共 N 项"
	got := markdownTitle(Message{Batch: true, Items: hanFreeItems(3)})
	assertChineseMessage(t, "markdownTitle(batch)", got)
	if !strings.Contains(got, "漏洞摘要") || !strings.Contains(got, "共 3 项") {
		t.Errorf("다건 제목이 '漏洞摘要 · 共 3 项' 형태여야 합니다, 받은 값 %q", got)
	}
	// 항목 없음: "漏洞通知"
	empty := markdownTitle(Message{})
	assertChineseMessage(t, "markdownTitle(empty)", empty)
	if empty != "漏洞通知" {
		t.Errorf("빈 메시지 제목은 '漏洞通知' 이어야 합니다, 받은 값 %q", empty)
	}
}

// TestMarkdownBatchIntroLocalized 는 다이제스트 머리말(시간창·건수·초과 안내)이
// 한국어인지 검사한다. 항목 데이터가 한자 0 이므로 출력에 한자가 보이면 머리말
// 문구가 중국어로 회귀한 것이다.
func TestMarkdownBatchIntroLocalized(t *testing.T) {
	// PR#1 把 markdownBatchIntro(markdown.go) 从韩语改成中文，这里反向断言：
	// ASCII 数据下输出不应再含韩文谚字。
	items := hanFreeItems(3)

	// 시간창 있음: "最近 N 分钟新增漏洞 N 项"
	win := markdownBatchIntro(Message{Batch: true, WindowMinutes: 30}, items, 3)
	if hasHangul(win) {
		t.Errorf("시간창 머리말에 한글이 남았습니다: %q", win)
	}
	if !strings.Contains(win, "最近 30 分钟") || !strings.Contains(win, "新增漏洞 3 项") {
		t.Errorf("시간창 머리말이 '最近 30 分钟新增漏洞 3 项' 형태여야 합니다, 받은 값 %q", win)
	}

	// 시간창 없음: "新增漏洞 N 项"(分钟 표기 없음)
	noWin := markdownBatchIntro(Message{Batch: true}, items, 3)
	if hasHangul(noWin) {
		t.Errorf("머리말에 한글이 남았습니다: %q", noWin)
	}
	if !strings.Contains(noWin, "新增漏洞 3 项") || strings.Contains(noWin, "分钟") {
		t.Errorf("시간창 없는 머리말은 '新增漏洞 3 项'(分钟 표기 없음)이어야 합니다, 받은 값 %q", noWin)
	}

	// 일부만 담겼을 때: "（此消息仅显示前 N 项，其余 N 项将在下一条消息中发送）"
	trunc := markdownBatchIntro(Message{Batch: true, WindowMinutes: 30}, items, 5)
	if hasHangul(trunc) {
		t.Errorf("초과 안내에 한글이 남았습니다: %q", trunc)
	}
	for _, want := range []string{"仅显示前 3 项", "其余 2 项", "下一条消息"} {
		if !strings.Contains(trunc, want) {
			t.Errorf("초과 안내에 %q 가 있어야 합니다, 받은 값 %q", want, trunc)
		}
	}
}

// TestWriteItemLabelsLocalized 는 단건 상세 렌더(markdown 3채널 공유 경로)의
// 필드 라벨(상태 변경·유형·자산·개요)과 상세 링크가 한국어인지 검사한다.
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

	// PR#1 把 writeItem(markdown.go) 的字段标签从韩语改成中文，这里反向断言。
	assertChineseMessage(t, "writeItem(single)", got)
	for _, want := range []string{
		"**状态变更**: 待处理 → 已修复",
		"**类型**: SQLi",
		"**资产**: a.example.com",
		"**概述**: SQL injection via q param",
		"[查看详情](https://platform.example/finding/1)",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("단건 렌더에 %q 가 있어야 합니다:\n%s", want, got)
		}
	}
}

// TestMarkdownBodyFooterLocalized 는 다건 본문 끝의 플랫폼 입구 링크가
// 한국어인지 검사한다(HomeURL 이 있을 때만 붙는다).
func TestMarkdownBodyFooterLocalized(t *testing.T) {
	m := Message{Batch: true, HomeURL: "https://platform.example", Items: hanFreeItems(2)}
	body, kept := markdownBody(m, 0)
	if kept != 2 {
		t.Fatalf("한도 없음(0)이면 2건 모두 담겨야 합니다, 받은 값 %d", kept)
	}
	if !strings.Contains(body, "[在平台中查看全部](https://platform.example)") {
		t.Errorf("본문 끝에 '在平台中查看全部' 링크가 있어야 합니다:\n%s", body)
	}
}
