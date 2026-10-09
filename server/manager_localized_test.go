package server

import "testing"

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
func TestSetWorkersErrorLocalized(t *testing.T) {
	err := (&Manager{}).SetWorkers(0)
	if err == nil {
		t.Fatal("测试文本 测试文本 测试文本 0 测试文本 测试文本 测试文本")
	}
	assertChineseMessage(t, "set_workers_nonpositive", err.Error())
}
