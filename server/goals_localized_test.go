package server

import "testing"

// 说明。
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
func TestGoalResumeErrorsLocalized(t *testing.T) {
	for _, c := range []struct {
		label, msg string
	}{
		{"resume_terminal", errGoalResumeTerminal},
		{"resume_not_paused", errGoalResumeNotPaused},
	} {
		assertChineseMessage(t, c.label, c.msg)
	}
}
