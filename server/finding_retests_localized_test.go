package server

import "testing"

// 说明。
// 说明。
// 说明。
// 说明。
// 说明。
// 说明。
//
// 说明。
// 说明。
// 说明。
// 说明。
func TestFindingRetestErrorsLocalized(t *testing.T) {
	for _, c := range []struct {
		label, msg string
	}{
		{"notes_too_long", errFindingRetestNotesTooLong},
		{"agent_missing", errFindingRetestAgentMissing},
		{"tool_required", errFindingRetestToolRequired},
		{"service_stopping", errFindingRetestServiceStopping},
	} {
		assertChineseMessage(t, c.label, c.msg)
	}
}
