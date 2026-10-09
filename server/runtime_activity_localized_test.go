package server

import (
	"fmt"
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
// 说明。
func TestRuntimeActivitySummariesLocalized(t *testing.T) {
	// 说明。
	plain := map[string]string{
		"goalless_task_done": goallessTaskDoneSummary,
		"goal_breakdown_r0":  goalBreakdownRound0Summary,
		"queued_no_llm":      queuedNoLLMSummary,
		"queued_fifo":        queuedFIFOSummary,
	}
	for label, msg := range plain {
		assertChineseMessage(t, label, msg)
	}

	// 说明。
	// 说明。
	fmtCases := map[string]string{
		"planner_round":       plannerRoundSummaryFmt,
		"timeout_final_round": timeoutFinalRoundSummaryFmt,
		"queued_conc_limit":   queuedConcurrencyLimitSummaryFmt,
	}
	for label, f := range fmtCases {
		if !strings.Contains(f, "%d") {
			t.Fatalf("%s: 测试文本 测试文本 %%d 测试文本 测试文本: %q", label, f)
		}
		got := fmt.Sprintf(f, 7)
		assertChineseMessage(t, label, got)
		if !strings.Contains(got, "7") {
			t.Fatalf("%s: 测试文本 测试文本 测试文本 测试文本 测试文本: %q", label, got)
		}
	}
}
