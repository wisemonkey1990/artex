package agent

import (
	"strings"
	"testing"
	)

// Wrap-up prompts shown to users must request Simplified Chinese summaries.
func TestWrapupPromptsUseSimplifiedChinese(t *testing.T) {
	userFacing := map[string]string{
		"settleWrapUpPrompt":     settleWrapUpPrompt,
		"mainAgentWrapUpDefault": mainAgentWrapUpDefault,
		"genericWrapUpDefault":   genericWrapUpDefault,
		"workerTaskTimeoutDefault": workerTaskTimeoutDefault,
	}
	for name, prompt := range userFacing {
		if !strings.Contains(prompt, "简体中文") {
			t.Errorf("%s: must request Simplified Chinese output, got %q", name, prompt)
		}
	}
	mustContain := func(name, prompt string, subs ...string) {
		for _, sub := range subs {
			if !strings.Contains(prompt, sub) {
				t.Errorf("%s: expected instruction %q", name, sub)
			}
		}
	}
	mustContain("settleWrapUpPrompt", settleWrapUpPrompt, "insert_assets", "record_fact", "report_finding", "一句", "纯文本")
	mustContain("workerTaskTimeoutDefault", workerTaskTimeoutDefault, "insert_assets", "record_fact", "report_finding", "一句", "纯文本")
	mustContain("genericWrapUpDefault", genericWrapUpDefault, "一句", "纯文本")
	mustContain("mainAgentWrapUpDefault", mainAgentWrapUpDefault, "一句", "纯文本")
	mustContain("plannerWrapUpDefault", plannerWrapUpDefault, "add_intent", "prove_goal", "TodoWrite")
	mustContain("plannerTaskTimeoutDefault", plannerTaskTimeoutDefault, "prove_goal")
}

// 说明。
// 说明。
func TestWrapupDefaultsRouting(t *testing.T) {
	if WrapupDefault("worker") != settleWrapUpPrompt {
		t.Error("worker per-run 测试文本 settleWrapUpPrompt 测试文本 测试文本")
	}
	if WrapupDefault("planner") != plannerWrapUpDefault {
		t.Error("planner per-run 测试文本 plannerWrapUpDefault 测试文本 测试文本")
	}
	if WrapupDefault("mainagent") != mainAgentWrapUpDefault {
		t.Error("mainagent per-run 测试文本 mainAgentWrapUpDefault 测试文本 测试文本")
	}
	// 说明。
	if WrapupDefault("unknown-agent") != genericWrapUpDefault {
		t.Error("测试文本 测试文本 genericWrapUpDefault 测试文本 测试文本 测试文本")
	}
	// 说明。
	if TaskTimeoutWrapupDefault("worker") != workerTaskTimeoutDefault {
		t.Error("worker task-timeout 测试文本 workerTaskTimeoutDefault 测试文本 测试文本")
	}
	if TaskTimeoutWrapupDefault("planner") != plannerTaskTimeoutDefault {
		t.Error("planner task-timeout 测试文本 plannerTaskTimeoutDefault 测试文本 测试文本")
	}
	if TaskTimeoutWrapupDefault("mainagent") != "" {
		t.Error("mainagent 测试文本 task-timeout 测试文本 测试文本 测试文本(测试文本 测试文本)")
	}
	// 说明。
	if workerTaskTimeoutDefault == settleWrapUpPrompt {
		t.Error("worker 测试文本 per-run 测试文本 task-timeout 测试文本 测试文本")
	}
	if plannerTaskTimeoutDefault == plannerWrapUpDefault {
		t.Error("planner 测试文本 per-run 测试文本 task-timeout 测试文本 测试文本")
	}
}
