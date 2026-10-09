package notify

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"testing"
)

// 说明。
// 说明。
// 说明。
func asciiStatusChangeItem() Item {
	return Item{
		FindingID:  7,
		Name:       "login-flaw",
		VulnClass:  "SQLi",
		Severity:   "high",
		Summary:    "SQL injection via q param",
		Assets:     []string{"a.example.com"},
		DetailURL:  "https://platform.example/finding/7",
		FromStatus: "pending",
		ToStatus:   "fixed",
	}
}

func mustJSON(t *testing.T, v any) string {
	t.Helper()
	b, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("测试文本 测试文本 测试文本: %v", err)
	}
	return string(b)
}

// 说明。
// 说明。
// 说明。
// 说明。
func TestFeishuItemLinesLocalized(t *testing.T) {
	got := feishuItemLines(asciiStatusChangeItem())
	assertKorean(t, "feishuItemLines", got)
	for _, want := range []string{
		"**测试文本 测试文本**: 测试文本 测试文本 → 测试文本",
		"**测试文本**: SQLi",
		"**测试文本**: a.example.com",
		"**测试文本**: SQL injection via q param",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("飞书 测试文本 测试文本 %q 测试文本 测试文本 测试文本:\n%s", want, got)
		}
	}
}

// 说明。
// 说明。
func TestFeishuCardButtonsLocalized(t *testing.T) {
	// 说明。
	single, _ := feishuCard(Message{Items: []Item{asciiStatusChangeItem()}})
	singleJSON := mustJSON(t, single)
	if hasHan(singleJSON) {
		t.Errorf("测试文本 飞书 测试文本 测试文本 测试文本 测试文本:\n%s", singleJSON)
	}
	if !strings.Contains(singleJSON, "查看详情") {
		t.Errorf("测试文本 飞书 测试文本 '测试文本 测试文本' 测试文本 测试文本 测试文本:\n%s", singleJSON)
	}

	// 说明。
	batch, _ := feishuCard(Message{Batch: true, HomeURL: "https://platform.example", Items: hanFreeItems(2)})
	batchJSON := mustJSON(t, batch)
	if hasHan(batchJSON) {
		t.Errorf("测试文本 飞书 测试文本 测试文本 测试文本 测试文本:\n%s", batchJSON)
	}
	if !strings.Contains(batchJSON, "测试文本 测试文本 测试文本") {
		t.Errorf("测试文本 飞书 测试文本 '测试文本 测试文本 测试文本' 测试文本 测试文本 测试文本:\n%s", batchJSON)
	}
}

// 说明。
// 说明。
func TestDingTalkActionCardButtonLocalized(t *testing.T) {
	var singleTitle string
	srv := capturePost(t, `{"errcode":0,"errmsg":"ok"}`, func(_ *testing.T, body map[string]any, _ *http.Request) {
		card, _ := body["actionCard"].(map[string]any)
		singleTitle, _ = card["singleTitle"].(string)
	})
	m := Message{Items: []Item{asciiStatusChangeItem()}}
	if _, err := (dingTalkChannel{}).Send(context.Background(), map[string]any{"webhook": srv.URL}, m); err != nil {
		t.Fatalf("测试文本递 测试文本: %v", err)
	}
	if singleTitle != "查看详情" {
		t.Errorf("钉钉 ActionCard singleTitle 测试文本 '测试文本 测试文本' 测试文本 测试文本, 测试文本 测试文本 %q", singleTitle)
	}
}

// 说明。
// 说明。
func TestChinaPlatformValidateLocalized(t *testing.T) {
	for _, kind := range []string{KindDingTalk, KindFeishu, KindWeCom} {
		ch, ok := Get(kind)
		if !ok {
			t.Fatalf("%s 测试文本 测试文本 测试文本 测试文本", kind)
		}

		// 说明。
		missing := ch.Validate(map[string]any{})
		if missing == nil {
			t.Fatalf("%s: 测试文本 测试文本 测试文本 测试文本 测试文本", kind)
		}
		assertKorean(t, kind+" missing", missing.Error())
		if !strings.Contains(missing.Error(), "Webhook") || !strings.Contains(missing.Error(), "测试文本") {
			t.Errorf("%s: 测试文本 测试文本 'Webhook 测试文本 测试文本' 测试文本 测试文本, 测试文本 测试文本 %q", kind, missing)
		}

		// 说明。
		bad := ch.Validate(map[string]any{"webhook": "ftp://x"})
		if bad == nil {
			t.Fatalf("%s: ftp 测试文本 测试文本 测试文本 测试文本", kind)
		}
		if hasHan(bad.Error()) {
			t.Errorf("%s: 测试文本 测试文本 测试文本 测试文本 测试文本 测试文本: %q", kind, bad)
		}
		if !strings.Contains(bad.Error(), "Webhook 地址无效") {
			t.Errorf("%s: 测试文本 测试文本 测试文本 'Webhook 测试文本 测试文本 测试文本' 测试文本 测试文本 测试文本, 测试文本 测试文本 %q", kind, bad)
		}
	}
}

// 说明。
// 说明。
// 说明。
// 说明。
func TestChinaPlatformSendErrorsLocalized(t *testing.T) {
	// 说明。
	bizCases := []struct {
		kind     string
		resp     string
		wantSubs []string
	}{
		{KindDingTalk, `{"errcode":310000,"errmsg":"keyword not matched"}`, []string{"DingTalk测试文本 测试文本 测试文本", "310000"}},
		{KindFeishu, `{"code":19021,"msg":"sign error"}`, []string{"Feishu测试文本 测试文本 测试文本", "19021"}},
		{KindWeCom, `{"errcode":45009,"errmsg":"freq out of limit"}`, []string{"WeCom 测试文本 测试文本", "45009"}},
		{KindWeCom, `{"errcode":93000,"errmsg":"invalid webhook"}`, []string{"WeCom测试文本 测试文本 测试文本", "93000"}},
	}
	for _, tc := range bizCases {
		srv := capturePost(t, tc.resp, nil)
		ch, ok := Get(tc.kind)
		if !ok {
			t.Fatalf("%s 测试文本 测试文本 测试文本 测试文本", tc.kind)
		}
		_, err := ch.Send(context.Background(), map[string]any{"webhook": srv.URL}, singleMsg())
		if err == nil {
			t.Fatalf("%s: %s 测试文本 测试文本 测试文本", tc.kind, tc.resp)
		}
		if hasHan(err.Error()) {
			t.Errorf("%s: 测试文本 测试文本 测试文本 测试文本 测试文本: %q", tc.kind, err)
		}
		for _, sub := range tc.wantSubs {
			if !strings.Contains(err.Error(), sub) {
				t.Errorf("%s: 测试文本 测试文本 %q 测试文本 测试文本 测试文本, 测试文本 测试文本 %q", tc.kind, sub, err)
			}
		}
	}

	// 说明。
	parseCases := []struct {
		kind     string
		platform string
	}{
		{KindDingTalk, "DingTalk"},
		{KindFeishu, "Feishu"},
		{KindWeCom, "WeCom"},
	}
	for _, pc := range parseCases {
		srv := capturePost(t, `not-json`, nil)
		ch, ok := Get(pc.kind)
		if !ok {
			t.Fatalf("%s 测试文本 测试文本 测试文本 测试文本", pc.kind)
		}
		_, err := ch.Send(context.Background(), map[string]any{"webhook": srv.URL}, singleMsg())
		if err == nil {
			t.Fatalf("%s: JSON 测试文本 测试文本 测试文本 测试文本 测试文本", pc.kind)
		}
		want := pc.platform + " 测试文本 测试文本 测试文本"
		if !strings.Contains(err.Error(), want) {
			t.Errorf("%s: 测试文本 测试文本 测试文本 %q 测试文本 测试文本 测试文本, 测试文本 测试文本 %q", pc.kind, want, err)
		}
	}
}
