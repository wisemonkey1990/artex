package server

import (
	"errors"
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

func TestUpdateMessageConstantsLocalized(t *testing.T) {
	assertChineseMessage(t, "updateMsgPreparing", updateMsgPreparing)
	assertChineseMessage(t, "updateMsgFailed", updateMsgFailed)
	assertChineseMessage(t, "updateMsgStaged", updateMsgStaged)
	assertChineseMessage(t, "updateErrInProgress", updateErrInProgress)
	assertChineseMessage(t, "updateErrRollbackInProgress", updateErrRollbackInProgress)

	// 说明。
	notRelease := fmt.Sprintf(updateErrNotReleaseFmt, "v0.0.0-dev")
	assertChineseMessage(t, "updateErrNotReleaseFmt", notRelease)
	if !strings.Contains(notRelease, "v0.0.0-dev") {
		t.Fatalf("updateErrNotReleaseFmt: 测试文本 测试文本 测试文本 测试文本: %q", notRelease)
	}
	latest := fmt.Sprintf(updateErrAlreadyLatestFmt, "v1.2.3")
	assertChineseMessage(t, "updateErrAlreadyLatestFmt", latest)
	if !strings.Contains(latest, "v1.2.3") {
		t.Fatalf("updateErrAlreadyLatestFmt: 测试文本 测试文本 测试文本 测试文本: %q", latest)
	}
}

// 说明。
// 说明。
func TestUpdateProgressMessagesLocalized(t *testing.T) {
	h := &updateHub{subs: map[chan updateProgress]struct{}{}}

	if !h.begin("v1.2.3") {
		t.Fatal("begin 测试文本 测试文本 测试文本 true 测试文本 测试文本")
	}
	cur, running := h.snapshot()
	if !running {
		t.Fatal("begin 测试文本 running 测试文本 测试文本")
	}
	if cur.Message != updateMsgPreparing {
		t.Fatalf("测试文本 测试文本 测试文本 测试文本: %q", cur.Message)
	}
	if h.begin("v1.2.4") {
		t.Fatal("测试文本 测试文本 begin 测试文本 false 测试文本 测试文本")
	}

	h.finish(errors.New("测试文本 测试文本"))
	cur, running = h.snapshot()
	if running {
		t.Fatal("finish 测试文本 running 测试文本 测试文本 测试文本")
	}
	if cur.Message != updateMsgFailed {
		t.Fatalf("测试文本 测试文本 测试文本: %q", cur.Message)
	}

	h.finish(nil)
	cur, _ = h.snapshot()
	if cur.Message != updateMsgStaged {
		t.Fatalf("测试文本 测试文本 测试文本 测试文本: %q", cur.Message)
	}
}

// 说明。
func TestUpdateRollbackInProgressLocalized(t *testing.T) {
	updHub.mu.Lock()
	prev := updHub.running
	updHub.running = true
	updHub.mu.Unlock()
	defer func() {
		updHub.mu.Lock()
		updHub.running = prev
		updHub.mu.Unlock()
	}()

	s := &Server{}
	req := httptest.NewRequest(http.MethodPost, "/api/update/rollback", nil)
	rec := httptest.NewRecorder()
	s.updateRollback(rec, req)

	if rec.Code != http.StatusConflict {
		t.Fatalf("测试文本 测试文本 409 测试文本, 测试文本 %d (测试文本 %s)", rec.Code, rec.Body.Bytes())
	}
	msg := decodeErrorField(t, rec.Body.Bytes())
	assertChineseMessage(t, "updateRollback 测试文本 测试文本 409", msg)
	if msg != updateErrRollbackInProgress {
		t.Fatalf("测试文本 测试文本 测试文本 测试文本: %q", msg)
	}
}
