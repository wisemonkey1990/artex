package server

import (
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
// 说明。
func TestGoalErrorConstantsLocalized(t *testing.T) {
	cases := map[string]string{
		"task_deleting_add":    errGoalTaskDeletingAdd,
		"task_deleting_edit":   errGoalTaskDeletingEdit,
		"task_deleting_delete": errGoalTaskDeletingDelete,
		"text_empty":           errGoalTextEmpty,
		"not_found":            errGoalNotFound,
		"read_after_add":       errGoalReadAfterAdd,
		"read_after_edit":      errGoalReadAfterEdit,
	}
	for label, msg := range cases {
		assertChineseMessage(t, label, msg)
	}
}

// 说明。
// 说明。
// 说明。
// 说明。
func TestGoalHandlersResponsesLocalized(t *testing.T) {
	// 说明。
	deletingServer := func() *Server {
		s := &Server{m: &Manager{tasks: map[string]*Task{"t1": {ID: "t1"}}}, engine: &Engine{}}
		s.engine.deleting.Store("t1", true)
		return s
	}
	// 说明。
	liveServer := func() *Server {
		return &Server{m: &Manager{tasks: map[string]*Task{"t1": {ID: "t1"}}}, engine: &Engine{}}
	}

	newReq := func(body, gid string) *http.Request {
		req := httptest.NewRequest(http.MethodPost, "/api/tasks/t1/goals", strings.NewReader(body))
		req.SetPathValue("id", "t1")
		if gid != "" {
			req.SetPathValue("gid", gid)
		}
		return req
	}

	cases := []struct {
		name    string
		handler func(*Server) http.HandlerFunc
		server  func() *Server
		body    string
		gid     string
		code    int
		want    string
	}{
		{
			name:    "add/task-deleting",
			handler: func(s *Server) http.HandlerFunc { return s.addGoal },
			server:  deletingServer,
			body:    `{"text":"SQL 测试文本 测试文本"}`,
			code:    http.StatusConflict,
			want:    errGoalTaskDeletingAdd,
		},
		{
			name:    "add/text-empty",
			handler: func(s *Server) http.HandlerFunc { return s.addGoal },
			server:  liveServer,
			body:    `{"text":"   "}`,
			code:    http.StatusBadRequest,
			want:    errGoalTextEmpty,
		},
		{
			name:    "edit/task-deleting",
			handler: func(s *Server) http.HandlerFunc { return s.editGoal },
			server:  deletingServer,
			body:    `{"text":"测试文本 测试文本"}`,
			gid:     "5",
			code:    http.StatusConflict,
			want:    errGoalTaskDeletingEdit,
		},
		{
			name:    "edit/text-empty",
			handler: func(s *Server) http.HandlerFunc { return s.editGoal },
			server:  liveServer,
			body:    `{"text":"  "}`,
			gid:     "5",
			code:    http.StatusBadRequest,
			want:    errGoalTextEmpty,
		},
		{
			name:    "delete/task-deleting",
			handler: func(s *Server) http.HandlerFunc { return s.deleteGoal },
			server:  deletingServer,
			body:    ``,
			gid:     "5",
			code:    http.StatusConflict,
			want:    errGoalTaskDeletingDelete,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			s := c.server()
			rec := httptest.NewRecorder()
			c.handler(s)(rec, newReq(c.body, c.gid))
			if rec.Code != c.code {
				t.Fatalf("测试文本 测试文本 = %d, 测试文本 = %d (测试文本 %q)", rec.Code, c.code, rec.Body.String())
			}
			got := decodeErrorField(t, rec.Body.Bytes())
			if got != c.want {
				t.Fatalf("测试文本 测试文本 = %q, 测试文本 = %q", got, c.want)
			}
			assertChineseMessage(t, c.name, got)
		})
	}
}
