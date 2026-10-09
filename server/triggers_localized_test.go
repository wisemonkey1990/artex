package server

import "testing"

// 说明。
// 说明。

// 说明。
func TestTriggerErrorConstantsLocalized(t *testing.T) {
	cases := []struct {
		name string
		msg  string
	}{
		{"no_condition", errTriggerNoCondition},
		{"tool_set_empty", errTriggerToolSetEmpty},
		{"custom_only", errTriggerCustomOnly},
	}
	for _, c := range cases {
		assertChineseMessage(t, c.name, c.msg)
	}
}

// 说明。
// 说明。
// 说明。
func TestValidateTriggerLocalized(t *testing.T) {
	if msg := validateTrigger(&triggerReq{}); msg != errTriggerNoCondition {
		t.Fatalf("测试文本 测试文本 测试文本 = %q, 测试文本 = %q", msg, errTriggerNoCondition)
	}
	if msg := validateTrigger(&triggerReq{OnToolCall: true}); msg != errTriggerToolSetEmpty {
		t.Fatalf("测试文本 测试文本 测试文本 = %q, 测试文本 = %q", msg, errTriggerToolSetEmpty)
	}
	if msg := validateTrigger(&triggerReq{OnFinding: true}); msg != "" {
		t.Fatalf("测试文本 测试文本 测试文本 = %q, 测试文本 = 测试文本 测试文本", msg)
	}
	if msg := validateTrigger(&triggerReq{OnToolCall: true, ToolNames: []string{"nmap"}}); msg != "" {
		t.Fatalf("测试文本 测试文本 测试文本 测试文本 = %q, 测试文本 = 测试文本 测试文本", msg)
	}
}
