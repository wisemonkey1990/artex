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

// per-run 과 task-timeout 은 의미가 달라야 한다(특히 planner): per-run 은 "이번 라운드만
// 끝난다"이고 task-timeout 은 "작업 전체가 끝난다"이다. 상수 매핑이 바뀌어 섞이면 안 된다.
func TestWrapupDefaultsRouting(t *testing.T) {
	if WrapupDefault("worker") != settleWrapUpPrompt {
		t.Error("worker per-run 기본값이 settleWrapUpPrompt 가 아니다")
	}
	if WrapupDefault("planner") != plannerWrapUpDefault {
		t.Error("planner per-run 기본값이 plannerWrapUpDefault 가 아니다")
	}
	if WrapupDefault("mainagent") != mainAgentWrapUpDefault {
		t.Error("mainagent per-run 기본값이 mainAgentWrapUpDefault 가 아니다")
	}
	// 미등록 키(커스텀 에이전트)는 generic 으로 떨어진다.
	if WrapupDefault("unknown-agent") != genericWrapUpDefault {
		t.Error("미등록 키가 genericWrapUpDefault 로 떨어지지 않는다")
	}
	// task-timeout 은 worker/planner 에만 있고, 그 외는 빈 문자열(호출부가 per-run 으로 회귀).
	if TaskTimeoutWrapupDefault("worker") != workerTaskTimeoutDefault {
		t.Error("worker task-timeout 기본값이 workerTaskTimeoutDefault 가 아니다")
	}
	if TaskTimeoutWrapupDefault("planner") != plannerTaskTimeoutDefault {
		t.Error("planner task-timeout 기본값이 plannerTaskTimeoutDefault 가 아니다")
	}
	if TaskTimeoutWrapupDefault("mainagent") != "" {
		t.Error("mainagent 은 task-timeout 문구가 없어야 한다(빈 문자열)")
	}
	// per-run 과 task-timeout 문구가 동일하면 의미 구분이 사라진 것이다.
	if workerTaskTimeoutDefault == settleWrapUpPrompt {
		t.Error("worker 의 per-run 과 task-timeout 문구가 동일하다")
	}
	if plannerTaskTimeoutDefault == plannerWrapUpDefault {
		t.Error("planner 의 per-run 과 task-timeout 문구가 동일하다")
	}
}
