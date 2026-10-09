package sidequestion

import (
	"errors"
	"fmt"
	"testing"
	"unicode"
)

// assertKorean fails if msg is empty, carries a CJK Han ideograph (= leftover
// untranslated Chinese), or has no Hangul at all. ASCII field names and the
// "1–4000"-style ranges are fine; only Han marks an unlocalized string.
func assertKorean(t *testing.T, label, msg string) {
	t.Helper()
	if msg == "" {
		t.Fatalf("%s: 测试文本 测试文本", label)
	}
	hangul := false
	for _, r := range msg {
		if unicode.Is(unicode.Han, r) {
			t.Fatalf("%s: 测试文本 测试文本 测试文本 测试文本: %q", label, msg)
		}
		if unicode.Is(unicode.Hangul, r) {
			hangul = true
		}
	}
	if !hangul {
		t.Fatalf("%s: 测试文本 测试文本: %q", label, msg)
	}
}

// TestSideQuestionOutputsLocalized guards F21: every user-facing side-question
// 说明。
// 说明。
// revert to Chinese here is a user-visible regression and must fail the build.
func TestSideQuestionOutputsLocalized(t *testing.T) {
	assertKorean(t, "ErrContextBudget", ErrContextBudget.Error())
	assertKorean(t, "errSideModelInterrupted", errSideModelInterrupted.Error())
	assertKorean(t, "errSideNoAnswer", errSideNoAnswer.Error())
	assertKorean(t, "msgSideToolUnavailable", msgSideToolUnavailable)
	assertKorean(t, "errSideSummaryCallCap", errSideSummaryCallCap.Error())
	assertKorean(t, "sideSummaryFailedPrefix", sideSummaryFailedPrefix)
	assertKorean(t, "errSideSummaryIncomplete", errSideSummaryIncomplete.Error())
	assertKorean(t, "errSideSummaryOverBudget", errSideSummaryOverBudget.Error())
	assertKorean(t, "errSideHistoryCursor", errSideHistoryCursor.Error())
	assertKorean(t, "errSideCompactionStalled", errSideCompactionStalled.Error())
}

// TestSideQuestionBudgetSentinelPreserved: localizing the message must not break
// callers that branch on the ErrContextBudget sentinel with errors.Is (see
// context_test.go, which relies on it after the recovery path gives up).
func TestSideQuestionBudgetSentinelPreserved(t *testing.T) {
	if !errors.Is(fmt.Errorf("prepare: %w", ErrContextBudget), ErrContextBudget) {
		t.Fatal("ErrContextBudget 测试文本 测试文本(errors.Is)测试文本 测试文本")
	}
}

// TestSideQuestionBrainPreserved: the agent-brain prompts (the answering
// instruction and the summary instruction) must stay in their benchmarked
// 说明。
// 说明。
func TestSideQuestionBrainPreserved(t *testing.T) {
	for _, c := range []struct{ label, text string }{
		{"instruction", instruction},
		{"summaryInstruction", summaryInstruction},
	} {
		han := false
		for _, r := range c.text {
			if unicode.Is(unicode.Han, r) {
				han = true
				break
			}
		}
		if !han {
			t.Errorf("%s: 测试文本 测试文本 测试文本 测试文本 测试文本 测试文本(测试文本 测试文本 测试文本)", c.label)
		}
	}
}
