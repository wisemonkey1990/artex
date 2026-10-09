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
// 说明。
func TestIntentInterventionErrorConstantsLocalized(t *testing.T) {
	cases := []struct {
		name string
		msg  string
	}{
		{"request_too_large", errIntentRequestTooLarge},
		{"message_empty", errIntentMessageEmpty},
		{"message_too_long", errIntentMessageTooLong},
		{"bad_request_id", errIntentBadRequestID},
		{"task_deleting", errIntentTaskDeleting},
		{"task_paused", errIntentTaskPaused},
		{"task_queued", errIntentTaskQueued},
		{"task_terminal", errIntentTaskTerminal},
		{"task_settling", errIntentTaskSettling},
		{"inherited_readonly", errIntentInheritedReadonly},
		{"not_paused", errIntentNotPaused},
	}
	for _, c := range cases {
		assertChineseMessage(t, c.name, c.msg)
	}
}

// 说明。
// 说明。
// 说明。
// 说明。
func TestSendWorkerMessageInputValidationLocalized(t *testing.T) {
	s := &Server{m: &Manager{tasks: map[string]*Task{"t1": {ID: "t1"}}}}

	call := func(body string) *httptest.ResponseRecorder {
		req := httptest.NewRequest(http.MethodPost, "/api/tasks/t1/intents/1/message", strings.NewReader(body))
		req.SetPathValue("id", "t1")
		req.SetPathValue("iid", "1")
		rec := httptest.NewRecorder()
		s.sendWorkerMessage(rec, req)
		return rec
	}
	errBody := func(t *testing.T, rec *httptest.ResponseRecorder) string {
		t.Helper()
		var resp struct {
			Error string `json:"error"`
		}
		if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
			t.Fatalf("测试文本 JSON 测试文本 测试文本: %v (测试文本 %q)", err, rec.Body.String())
		}
		return resp.Error
	}

	cases := []struct {
		name string
		body string
		code int
		want string
	}{
		// 说明。
		{"request_too_large", `{"message":"` + strings.Repeat("a", maxWorkerMessageBytes+1024) + `"}`, http.StatusRequestEntityTooLarge, errIntentRequestTooLarge},
		{"message_empty", `{"message":"  "}`, http.StatusBadRequest, errIntentMessageEmpty},
		{"message_too_long", `{"message":"` + strings.Repeat("测试文本", 4001) + `"}`, http.StatusBadRequest, errIntentMessageTooLong},
		{"bad_request_id", `{"message":"测试文本","request_id":"测试文本 测试文本"}`, http.StatusBadRequest, errIntentBadRequestID},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			rec := call(c.body)
			if rec.Code != c.code {
				t.Fatalf("测试文本 测试文本 = %d, 测试文本 = %d (测试文本 %q)", rec.Code, c.code, rec.Body.String())
			}
			if got := errBody(t, rec); got != c.want {
				t.Fatalf("测试文本 测试文本 = %q, 测试文本 = %q", got, c.want)
			}
			assertChineseMessage(t, c.name, errBody(t, rec))
		})
	}
}
