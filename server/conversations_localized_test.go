package server

import (
	"encoding/json"
	"fmt"
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
func TestConversationErrorConstantsLocalized(t *testing.T) {
	cases := []struct {
		name string
		msg  string
	}{
		{"request_too_large", convErrRequestTooLarge},
		{"agent_key_empty", convErrAgentKeyEmpty},
		{"agent_key_too_long", fmt.Sprintf(convErrAgentKeyTooLong, maxConversationAgentKeyRunes)},
		{"agent_not_found", convErrAgentNotFound},
		{"llm_profile", convErrLLMProfile},
		{"title_too_long", fmt.Sprintf(convErrTitleTooLong, maxConversationTitleRunes)},
		{"title_or_pinned", convErrTitleOrPinned},
		{"title_empty", convErrTitleEmpty},
		{"bad_conv_id", convErrBadConvID},
		{"ids_count", fmt.Sprintf(convErrIDsCount, maxConversationDeleteBatch)},
		{"message_empty", convErrMessageEmpty},
		{"busy", convErrBusy},
	}
	for _, c := range cases {
		assertChineseMessage(t, c.name, c.msg)
	}
}

// 说明。
// 说明。
// 说明。
// 说明。
func TestDecodeConversationRequestTooLargeLocalized(t *testing.T) {
	// 说明。
	body := `{"title":"` + strings.Repeat("a", maxConversationRequestBytes+1024) + `"}`
	req := httptest.NewRequest(http.MethodPost, "/api/conversations", strings.NewReader(body))
	rec := httptest.NewRecorder()
	var dst struct {
		Title string `json:"title"`
	}
	if decodeConversationRequest(rec, req, &dst) {
		t.Fatal("测试文本 测试文本 测试文本 decodeConversationRequest 测试文本 true 测试文本 测试文本")
	}
	if rec.Code != http.StatusRequestEntityTooLarge {
		t.Fatalf("测试文本 测试文本 = %d, 测试文本 = %d", rec.Code, http.StatusRequestEntityTooLarge)
	}
	var resp struct {
		Error string `json:"error"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("测试文本 JSON 测试文本 测试文本: %v (测试文本=%q)", err, rec.Body.String())
	}
	if resp.Error != convErrRequestTooLarge {
		t.Fatalf("测试文本 error = %q, 测试文本 = %q", resp.Error, convErrRequestTooLarge)
	}
	assertChineseMessage(t, "decode.too_large.response", resp.Error)
}

// 说明。
// 说明。
// 说明。
func TestConversationDefaultTitlesLocalized(t *testing.T) {
	assertChineseMessage(t, "default_title", convDefaultTitle)
	assertChineseMessage(t, "attachment_title", convAttachmentTitle)
	if convDefaultTitle == convAttachmentTitle {
		t.Fatal("测试文本 测试文本 测试文本 测试文本 测试文本 测试文本 测试文本 测试文本")
	}
}

// 说明。
// 说明。
// 说明。
func TestIsDefaultConversationTitle(t *testing.T) {
	if !isDefaultConversationTitle("") {
		t.Fatal("测试文本 测试文本 测试文本 测试文本 测试文本 测试文本")
	}
	if !isDefaultConversationTitle(convDefaultTitle) {
		t.Fatalf("测试文本 测试文本 %q 测试文本 测试文本 测试文本 测试文本 测试文本", convDefaultTitle)
	}
	if isDefaultConversationTitle("测试文本 测试文本 测试文本") {
		t.Fatal("测试文本 测试文本 测试文本 测试文本 测试文本 测试文本 测试文本 测试文本")
	}
}

// 说明。
// 说明。
// 说明。
func TestConversationRetestReasonsLocalized(t *testing.T) {
	assertChineseMessage(t, "retest_failed_to_start", convRetestFailedToStart)
	assertChineseMessage(t, "retest_status_read_failed", convRetestStatusReadFailed)
	assertChineseMessage(t, "retest_stopped_or_closed", convRetestStoppedOrClosed)
}

// 说明。
// 说明。
// 说明。
func TestTranscriptErrorSummaryLocalized(t *testing.T) {
	cases := []struct {
		name  string
		label string
		errm  string
		want  string
	}{
		{"chat_turn", "", "connection reset", "(错误：connection reset)"},
		{"main_agent", "测试文本 测试文本", "connection reset", "(测试文本 测试文本 错误：connection reset)"},
	}
	for _, c := range cases {
		got := transcriptErrorSummary(c.label, c.errm)
		if got != c.want {
			t.Fatalf("%s: transcriptErrorSummary = %q, 测试文本 = %q", c.name, got, c.want)
		}
		if !strings.Contains(got, c.errm) {
			t.Fatalf("%s: err 测试文本 测试文本 测试文本: %q", c.name, got)
		}
		// 说明。
		if strings.ContainsAny(got, "（）：") {
			t.Fatalf("%s: 测试文本 测试文本 测试文本 测试文本: %q", c.name, got)
		}
		// 说明。
		// 说明。
		wrapper := strings.ReplaceAll(got, c.errm, "")
		assertChineseMessage(t, c.name+".wrapper", wrapper)
	}
}
