package notify

import (
	"errors"
	"strings"
	"testing"
)

// 说明。
// 说明。
// 说明。
// 说明。

// 说明。
// 说明。
func TestFilterValidateLocalized(t *testing.T) {
	err := Filter{MinSeverity: "hgih"}.Validate()
	if err == nil {
		t.Fatal("测试文本 测试文本 测试文本 测试文本 测试文本")
	}
	assertKorean(t, "Filter.Validate", err.Error())
	// 说明。
	for _, tok := range []string{"low", "medium", "high", "critical"} {
		if !strings.Contains(err.Error(), tok) {
			t.Errorf("测试文本 测试文本 %q 测试文本 测试文本: %q", tok, err.Error())
		}
	}
	// 说明。
	if err := (Filter{MinSeverity: "high"}).Validate(); err != nil {
		t.Errorf("high 测试文本 测试文本 测试文本 测试文本: %v", err)
	}
	if err := (Filter{}).Validate(); err != nil {
		t.Errorf("测试文本 测试文本 测试文本 测试文本 测试文本 测试文本: %v", err)
	}
}

// 说明。
// 说明。
func TestPrepareConfigUpdateUnknownKindLocalized(t *testing.T) {
	_, err := PrepareConfigUpdate("definitely-not-a-channel", map[string]any{}, map[string]any{})
	if err == nil {
		t.Fatal("测试文本 测试文本 测试文本 测试文本 测试文本")
	}
	assertKorean(t, "PrepareConfigUpdate unknown kind", err.Error())
}

// 说明。
// 说明。
// 说明。
func TestDestinationChangedErrorLocalized(t *testing.T) {
	_, err := PrepareConfigUpdate("webhook",
		map[string]any{
			"url":     "https://old.example.com/hook",
			"headers": map[string]any{"Authorization": "Bearer real-token"},
		},
		map[string]any{"url": "https://attacker.example.com/hook"},
	)
	if err == nil {
		t.Fatal("测试文本 测试文本 测试文本 测试文本 测试文本 测试文本 测试文本 测试文本 测试文本")
	}
	var de *ErrDestinationChangedWithoutCredentials
	if !errors.As(err, &de) {
		t.Fatalf("测试文本 测试文本 测试文本 测试文本: %T", err)
	}
	assertKorean(t, "ErrDestinationChangedWithoutCredentials", err.Error())
	// 说明。
	if !strings.Contains(err.Error(), "url") || !strings.Contains(err.Error(), "headers") {
		t.Errorf("测试文本 测试文本 测试文本 测试文本: %q", err.Error())
	}
}

// 说明。
// 说明。
func TestRejectMaskedInContainersLocalized(t *testing.T) {
	_, err := PrepareConfigUpdate("webhook",
		map[string]any{},
		map[string]any{"headers": map[string]any{"Authorization": MaskedPrefix + ":…abc123"}},
	)
	if err == nil {
		t.Fatal("测试文本 测试文本 测试文本 测试文本 测试文本 测试文本 测试文本 测试文本")
	}
	assertKorean(t, "rejectMaskedInContainers", err.Error())
	// 说明。
	if !strings.Contains(err.Error(), MaskedPrefix) {
		t.Errorf("测试文本 测试文本 测试文本 测试文本: %q", err.Error())
	}
}
