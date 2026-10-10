package report

import (
	"encoding/json"
	"strings"
	"testing"
	"time"
	"unicode"

	"github.com/Autumn-27/artex/db"
)

// containsHangul detects residual Korean text in a Chinese-localized report.
func containsHangul(s string) (rune, bool) {
	for _, r := range s {
		if unicode.Is(unicode.Hangul, r) {
			return r, true
		}
	}
	return 0, false
}

func sampleFindingNode(t *testing.T) *db.Node {
	t.Helper()
	payload, err := json.Marshal(map[string]any{
		"vulnclass": "SQL Injection",
		"name":      "登录表单 SQL 注入",
		"severity":  "high",
		"summary":   "已确认登录参数中存在基于错误回显的 SQL 注入。",
		"evidence":  map[string]string{"poc": "' OR '1'='1"},
	})
	if err != nil {
		t.Fatalf("payload 序列化失败: %v", err)
	}
	return &db.Node{ID: 1, Kind: "finding", Payload: payload}
}

func sampleDBFinding() *db.DBFinding {
	return &db.DBFinding{
		ID:              123,
		VulnClass:       "SQL Injection",
		Name:            "登录表单 SQL 注入",
		Severity:        "high",
		Summary:         "已确认登录参数中存在基于错误回显的 SQL 注入。",
		Evidence:        "' OR '1'='1",
		Status:          "confirmed",
		Report:          "详细分析正文。",
		CreatedAt:       time.Date(2026, 10, 4, 2, 30, 0, 0, time.UTC),
		EvidenceVersion: 2,
		TaskDescription: "演示目标渗透测试",
		TrafficBindings: []db.FindingTrafficBinding{
			{
				ID:   9,
				Role: "request",
				Note: "注入 payload 发送位置",
				Snapshot: db.TrafficEvidenceSnapshot{
					Method: "POST",
					URL:    "https://sandbox.local/login",
					Status: 200,
				},
			},
		},
	}
}

// TestReportMarkdownUsesSimplifiedChinese verifies the report outline uses Chinese and contains no Korean.
func TestReportMarkdownUsesSimplifiedChinese(t *testing.T) {
	out := Markdown(Input{
		Title:       "演示任务",
		Goal:        "全面检查沙箱",
		GeneratedAt: time.Date(2026, 10, 4, 2, 30, 0, 0, time.UTC),
		AssetCounts: map[string]int{"host": 2, "url": 5},
		Findings:    []*db.Node{sampleFindingNode(t)},
	})
	if r, ok := containsHangul(out); ok {
		t.Fatalf("报告骨架中仍有韩文字符 %q:\n%s", string(r), out)
	}
	for _, want := range []string{
		"# 渗透测试报告：演示任务",
		"- **任务目标**: 全面检查沙箱",
		"- **生成时间**: ",
		"## 摘要",
		"- 已确认漏洞：**1** 项",
		"## 漏洞",
		"**PoC / 证据：**",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("报告骨架中缺少 %q", want)
		}
	}
}

// TestReportMarkdownEmptyUsesChinese verifies the empty-report message.
func TestReportMarkdownEmptyUsesChinese(t *testing.T) {
	out := Markdown(Input{Title: "空任务", GeneratedAt: time.Now()})
	if r, ok := containsHangul(out); ok {
		t.Fatalf("空报告中仍有韩文字符 %q:\n%s", string(r), out)
	}
	if !strings.Contains(out, "_本次未发现已确认漏洞。_") {
		t.Errorf("空报告提示文案无效:\n%s", out)
	}
}

// TestFindingsMarkdownUsesSimplifiedChinese 는 漏洞摘要和详情 Markdown 导出骨架及严重程度标签使用中文且不含韩文。
func TestFindingsMarkdownUsesSimplifiedChinese(t *testing.T) {
	f := sampleDBFinding()
	out := FindingsMarkdown([]*db.DBFinding{f}, time.Date(2026, 10, 4, 2, 30, 0, 0, time.UTC))
	if r, ok := containsHangul(out); ok {
		t.Fatalf("漏洞摘要报告中仍有韩文字符 %q:\n%s", string(r), out)
	}
	for _, want := range []string{
		"# 漏洞摘要报告",
		"- **生成时间**: ",
		"- **漏洞总数**：1 项",
		"## 摘要",
		"| 严重程度 | 数量 |",
		"| 严重 |", "| 高危 |", "| 中危 |", "| 低危 |", // severity labels match the UI
		"## 漏洞详情",
		"- **类型**: SQL Injection",
		"- **状态**：已确认",
		"- **所属任务**: 演示目标渗透测试",
		"- **发现时间**: ",
		"**证据：**",
		"**详细报告：**",
		"## 相关流量证据",
		"证据版本：2，绑定数量：1。",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("摘要报告中缺少 %q", want)
		}
	}

	single := SingleFindingMarkdown(f, time.Date(2026, 10, 4, 2, 30, 0, 0, time.UTC))
	if r, ok := containsHangul(single); ok {
		t.Fatalf("单项报告中仍有韩文字符 %q:\n%s", string(r), single)
	}
	for _, want := range []string{
		"- **严重程度**：高危",
		"## 概述",
		"## 证据",
		"## 详细报告",
		"[请求内容](evidence/123/9/request.http)",
		"[响应内容](evidence/123/9/response.http)",
	} {
		if !strings.Contains(single, want) {
			t.Errorf("单项报告中缺少 %q", want)
		}
	}
}

// TestFindingsCSVHeaderUsesChinese verifies the CSV header uses Chinese. Data values remain unchanged.
func TestFindingsCSVHeaderUsesChinese(t *testing.T) {
	raw := FindingsCSV([]*db.DBFinding{sampleDBFinding()})
	out := strings.TrimPrefix(string(raw), "\xEF\xBB\xBF")
	header := strings.SplitN(out, "\n", 2)[0]
	if r, ok := containsHangul(header); ok {
		t.Fatalf("CSV 表头中仍有韩文字符 %q: %s", string(r), header)
	}
	for _, col := range []string{"ID", "名称", "类型", "严重程度", "状态", "所属任务", "发现时间", "概述", "流量证据数量", "流量证据 ID"} {
		if !strings.Contains(header, col) {
			t.Errorf("CSV 表头缺少 %q 列: %s", col, header)
		}
	}
}
