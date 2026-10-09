package server

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/Autumn-27/artex/db"
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
// 说明。
// 说明。
// 说明。
// 说明。

// 说明。
// 说明。
// 说明。
// 说明。
func TestControlWorkNoRunningWorkErrorLocalized(t *testing.T) {
	e := NewEngine(nil)
	err := e.ControlWork(context.Background(), 42, "pause")
	if err == nil {
		t.Fatal("测试文本 测试文本 work 测试文本 测试文本 ControlWork 测试文本 nil 测试文本 测试文本")
	}
	if !errors.Is(err, errWorkControlConflict) {
		t.Fatalf("errors.Is(err, errWorkControlConflict) = false, err=%v", err)
	}
	assertChineseMessage(t, "control_work_no_running", err.Error())
}

// 说明。
// 说明。
// 说明。
func TestEngineUserFacingErrorsLocalized(t *testing.T) {
	for _, c := range []struct {
		label, msg string
	}{
		{"work_control_busy", fmt.Errorf(errWorkControlBusyFmt, errWorkControlConflict, int64(42), "pause").Error()},
		{"work_control_wait", fmt.Errorf(errWorkControlWaitFmt, int64(42), "pause", context.DeadlineExceeded).Error()},
		{"detached_worker_not_ready", errDetachedWorkerNotReady},
		{"detached_state_conflict", fmt.Errorf(errIntentCtrlStateConflictFmt, db.ErrIntentStateConflict).Error()},
	} {
		assertChineseMessage(t, c.label, c.msg)
	}
}
