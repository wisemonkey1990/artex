package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// 说明。
// 说明。

// 说明。
// 说明。
// 说明。
func TestDeepenFindingBodyTooLargeLocalized(t *testing.T) {
	s := &Server{}
	body := `{"description":"` + strings.Repeat("a", 33<<10) + `"}`
	req := httptest.NewRequest(http.MethodPost, "/api/exploration/findings/1/deepen", strings.NewReader(body))
	req.SetPathValue("id", "1")
	rec := httptest.NewRecorder()
	s.deepenFinding(rec, req)

	if rec.Code != http.StatusRequestEntityTooLarge {
		t.Fatalf("测试文本 测试文本 = %d, 测试文本 = %d (测试文本 %q)", rec.Code, http.StatusRequestEntityTooLarge, rec.Body.String())
	}
	var resp struct {
		Error string `json:"error"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("测试文本 JSON 测试文本 测试文本: %v (测试文本 %q)", err, rec.Body.String())
	}
	const want = "请求正文过大"
	if resp.Error != want {
		t.Fatalf("测试文本 测试文本 = %q, 测试文本 = %q", resp.Error, want)
	}
	assertChineseMessage(t, "body_too_large", resp.Error)
}

// 说明。
// 说明。
// 说明。
// 说明。
func TestFindingFollowUpAuditSummaryLocalized(t *testing.T) {
	assertChineseMessage(t, "follow_up_summary", auditFindingFollowUpSummary)
}
