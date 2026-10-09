package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Autumn-27/artex/db"
)

// 说明。
// 说明。
// 说明。

// 说明。
// 说明。
// 说明。
// 说明。
// 说明。
func TestTaskControlErrorConstantsLocalized(t *testing.T) {
	cases := map[string]string{
		"task_deleting":          errTaskCtrlDeleting,
		"task_terminal_pause":    errTaskCtrlTerminalPause,
		"task_already_paused":    errTaskCtrlAlreadyPaused,
		"batch_size":             errTaskCtrlBatchSizeFmt,
		"intent_inherited":       errIntentCtrlInheritedReadonly,
		"intent_only_running":    errIntentCtrlOnlyRunningPause,
		"intent_only_paused":     errIntentCtrlOnlyPausedResume,
		"intent_state_conflict":  errIntentCtrlStateConflictFmt,
		"intent_only_deletable":  errIntentCtrlOnlyDeletable,
		"intent_reason_required": errIntentCtrlReasonRequired,
	}
	for label, msg := range cases {
		assertChineseMessage(t, label, msg)
	}
}

// 说明。
// 说明。
// 说明。
// 说明。
// 说明。
// 说明。
func TestTaskControlApplyResponsesLocalized(t *testing.T) {
	// 说明。
	terminalServer := func() (*Server, *Task) {
		tk := &Task{ID: "1", Status: "done"}
		return &Server{m: &Manager{tasks: map[string]*Task{"1": tk}}, engine: &Engine{}}, tk
	}
	// 说明。
	pausedServer := func() (*Server, *Task) {
		tk := &Task{ID: "1", Status: "running", Paused: true}
		return &Server{m: &Manager{tasks: map[string]*Task{"1": tk}}, engine: &Engine{}}, tk
	}
	// 说明。
	deletingServer := func() (*Server, *Task) {
		tk := &Task{ID: "1", Status: "running"}
		s := &Server{m: &Manager{tasks: map[string]*Task{"1": tk}}, engine: &Engine{}}
		s.engine.deleting.Store("1", true)
		return s, tk
	}

	cases := []struct {
		name  string
		setup func() (*Server, *Task)
		want  string
	}{
		{"terminal", terminalServer, errTaskCtrlTerminalPause},
		{"already-paused", pausedServer, errTaskCtrlAlreadyPaused},
		{"deleting", deletingServer, errTaskCtrlDeleting},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			s, tk := c.setup()
			_, err := s.applyTaskControl(tk, "pause")
			if err == nil {
				t.Fatalf("测试文本 测试文本 nil 测试文本 测试文本")
			}
			if err.Error() != c.want {
				t.Fatalf("测试文本 测试文本 = %q, 测试文本 = %q", err.Error(), c.want)
			}
			assertChineseMessage(t, c.name, err.Error())
		})
	}
}

// 说明。
// 说明。
// 说明。
// 说明。
// 说明。
func TestControlTasksBatchResponsesLocalized(t *testing.T) {
	t.Run("size-error", func(t *testing.T) {
		s := &Server{}
		req := httptest.NewRequest(http.MethodPost, "/api/tasks/control",
			strings.NewReader(`{"task_ids":[],"action":"pause"}`))
		rec := httptest.NewRecorder()
		s.controlTasksBatch(rec, req)
		if rec.Code != http.StatusBadRequest {
			t.Fatalf("测试文本 测试文本 = %d, 测试文本 = 400 (测试文本 %q)", rec.Code, rec.Body.String())
		}
		got := decodeErrorField(t, rec.Body.Bytes())
		want := fmt.Sprintf(errTaskCtrlBatchSizeFmt, maxBatchControlIDs)
		if got != want {
			t.Fatalf("测试文本 测试文本 = %q, 测试文本 = %q", got, want)
		}
		if !strings.Contains(got, "100") {
			t.Fatalf("测试文本 测试文本 100 测试文本 测试文本 测试文本 测试文本: %q", got)
		}
		assertChineseMessage(t, "size-error", got)
	})

	t.Run("item-error-surfaces", func(t *testing.T) {
		tk := &Task{ID: "1", Status: "done"}
		s := &Server{m: &Manager{tasks: map[string]*Task{"1": tk}}, engine: &Engine{}}
		req := httptest.NewRequest(http.MethodPost, "/api/tasks/control",
			strings.NewReader(`{"task_ids":["1"],"action":"pause"}`))
		rec := httptest.NewRecorder()
		s.controlTasksBatch(rec, req)
		if rec.Code != http.StatusOK {
			t.Fatalf("测试文本 测试文本 = %d, 测试文本 = 200 (测试文本 %q)", rec.Code, rec.Body.String())
		}
		var out struct {
			Items []struct {
				ID    string `json:"id"`
				OK    bool   `json:"ok"`
				Error string `json:"error"`
			} `json:"items"`
		}
		if err := json.Unmarshal(rec.Body.Bytes(), &out); err != nil {
			t.Fatalf("测试文本 JSON 测试文本 测试文本: %v (测试文本 %q)", err, rec.Body.String())
		}
		if len(out.Items) != 1 {
			t.Fatalf("items 测试文本 = %d, 测试文本 = 1", len(out.Items))
		}
		if out.Items[0].OK || out.Items[0].Error != errTaskCtrlTerminalPause {
			t.Fatalf("items[0] = %+v, error 测试文本 = %q", out.Items[0], errTaskCtrlTerminalPause)
		}
		assertChineseMessage(t, "item-error", out.Items[0].Error)
	})
}

// 说明。
// 说明。
// 说明。
func TestIntentStateConflictWrapsSentinel(t *testing.T) {
	err := fmt.Errorf(errIntentCtrlStateConflictFmt, db.ErrIntentStateConflict)
	if !errors.Is(err, db.ErrIntentStateConflict) {
		t.Fatalf("errors.Is(.., ErrIntentStateConflict) = false, %%w 测试文本 测试文本: %q", err.Error())
	}
	if !strings.Contains(err.Error(), "paused") {
		t.Fatalf("enum paused 测试文本 测试文本 测试文本 测试文本: %q", err.Error())
	}
	assertChineseMessage(t, "state-conflict", err.Error())
}
