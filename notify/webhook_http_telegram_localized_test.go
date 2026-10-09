package notify

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

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
// 说明。
func TestTelegramItemLabelsLocalized(t *testing.T) {
	m := Message{
		HomeURL: "https://example.com/panel",
		Items: []Item{{
			Name:       "sqli-login", // Title() = Name
			VulnClass:  "injection",
			Severity:   "high",
			Summary:    "login form is injectable",
			Assets:     []string{"host-a.example.com"},
			DetailURL:  "https://example.com/f/1",
			FromStatus: "pending",
			ToStatus:   "fixed",
		}},
	}
	out, kept := telegramHTML(m)
	if kept != 1 {
		t.Fatalf("测试文本 测试文本 1 测试文本 测试文本 测试文本, 测试文本 测试文本 %d", kept)
	}
	if hasHan(out) {
		t.Errorf("Telegram 测试文本 测试文本 测试文本 测试文本 测试文本:\n%s", out)
	}
	for _, want := range []string{"测试文本 测试文本", "测试文本", "测试文本", "概述", "查看详情"} {
		if !strings.Contains(out, want) {
			t.Errorf("Telegram 测试文本 %q 测试文本 测试文本:\n%s", want, out)
		}
	}
}

// 说明。
// 说明。
// 说明。
func TestTelegramBatchTitleLocalized(t *testing.T) {
	items := []Item{{Name: "a", Severity: "high"}, {Name: "b", Severity: "low"}}

	// 说明。
	over := telegramBatchTitle(Message{WindowMinutes: 30}, items, 5)
	assertKorean(t, "telegramBatchTitle(测试文本)", over)
	for _, want := range []string{"测试文本 测试文本 · 测试文本 5测试文本", "测试文本 2测试文本 测试文本", "测试文本 3测试文本", "测试文本 测试文本", "测试文本 30测试文本"} {
		if !strings.Contains(over, want) {
			t.Errorf("测试文本 测试文本 %q 测试文本 测试文本: %q", want, over)
		}
	}

	// 说明。
	full := telegramBatchTitle(Message{}, items, 2)
	assertKorean(t, "telegramBatchTitle(测试文本)", full)
	if !strings.Contains(full, "测试文本 测试文本 · 测试文本 2测试文本") {
		t.Errorf("测试文本 测试文本 '测试文本 测试文本 · 测试文本 2测试文本' 测试文本 测试文本 测试文本: %q", full)
	}
	if strings.Contains(full, "测试文本 测试文本") {
		t.Errorf("测试文本 测试文本 '测试文本 测试文本' 测试文本 测试文本: %q", full)
	}
}

// 说明。
// 说明。
func TestTelegramValidateLocalized(t *testing.T) {
	if err := (telegramChannel{}).Validate(map[string]any{}); err == nil {
		t.Fatal("bot_token 测试文本 测试文本 测试文本 测试文本")
	} else {
		assertKorean(t, "telegram Validate(bot_token)", err.Error())
		if !strings.Contains(err.Error(), "Bot Token") {
			t.Errorf("Bot Token 测试文本 测试文本 测试文本: %q", err.Error())
		}
	}
	if err := (telegramChannel{}).Validate(map[string]any{"bot_token": "t"}); err == nil {
		t.Fatal("chat_id 测试文本 测试文本 测试文本 测试文本")
	} else {
		assertKorean(t, "telegram Validate(chat_id)", err.Error())
		if !strings.Contains(err.Error(), "Chat ID") {
			t.Errorf("Chat ID 测试文本 测试文本 测试文本: %q", err.Error())
		}
	}
}

// 说明。
// 说明。
func TestWebhookValidateLocalized(t *testing.T) {
	if err := (webhookChannel{}).Validate(map[string]any{}); err == nil {
		t.Fatal("url 测试文本 测试文本 测试文本 测试文本")
	} else {
		assertKorean(t, "webhook Validate(url)", err.Error())
		if !strings.Contains(err.Error(), "测试文本 URL") {
			t.Errorf("测试文本 URL 测试文本 测试文本 测试文本: %q", err.Error())
		}
	}
	// 说明。
	if err := (webhookChannel{}).Validate(map[string]any{"url": "https://example.com/hook", "method": "DELETE"}); err == nil {
		t.Fatal("DELETE 测试文本 测试文本 测试文本 测试文本")
	} else {
		assertKorean(t, "webhook Validate(method)", err.Error())
		if !strings.Contains(err.Error(), "GET/POST/PUT/PATCH") {
			t.Errorf("允许 测试文本 测试文本 测试文本 测试文本: %q", err.Error())
		}
	}
	// 说明。
	if err := (webhookChannel{}).Validate(map[string]any{"url": "https://example.com/hook", "body_template": "{{"}); err == nil {
		t.Fatal("测试文本 测试文本 测试文本 测试文本 测试文本")
	} else {
		assertKorean(t, "webhook Validate(template)", err.Error())
		if !strings.Contains(err.Error(), "测试文本") {
			t.Errorf("测试文本 测试文本 测试文本 测试文本: %q", err.Error())
		}
	}
}

// 说明。
// 说明。
// 说明。
func TestValidateHTTPURLLocalized(t *testing.T) {
	// 说明。
	if err := validateHTTPURL("ftp://example.com/x"); err == nil {
		t.Fatal("ftp 测试文本 测试文本 测试文本")
	} else {
		assertKorean(t, "validateHTTPURL(scheme)", err.Error())
		if !strings.Contains(err.Error(), "http") {
			t.Errorf("测试文本 测试文本(http/https)测试文本 测试文本 测试文本: %q", err.Error())
		}
	}
	// 说明。
	if err := validateHTTPURL("http://"); err == nil {
		t.Fatal("测试文本 测试文本 测试文本 测试文本 测试文本")
	} else {
		assertKorean(t, "validateHTTPURL(host)", err.Error())
	}
}

// 说明。
// 说明。
func TestHTTPLocalTargetBlockedLocalized(t *testing.T) {
	t.Setenv(AllowLocalTargetsEnv, "")
	err := blockInternalDial("tcp", "127.0.0.1:25", nil)
	if err == nil {
		t.Fatal("测试文本 测试文本 测试文本 拦截测试文本 测试文本")
	}
	assertKorean(t, "blockInternalDial", err.Error())
	if !strings.Contains(err.Error(), AllowLocalTargetsEnv) {
		t.Errorf("拦截 测试文本 测试文本 测试文本(%s)测试文本 测试文本 测试文本: %q", AllowLocalTargetsEnv, err.Error())
	}
}

// 说明。
// 说明。
func TestHTTPRedactHelpersLocalized(t *testing.T) {
	if got := redactRequestTarget(""); !strings.Contains(got, "测试文本 测试文本 测试文本") {
		t.Errorf("测试文本 测试文本 测试文本 测试文本 测试文本 测试文本, 测试文本 测试文本 %q", got)
	}
	// 说明。
	got := redactTransportError(&url.Error{Op: "Get", URL: "http://api.example", Err: nil})
	if hasHan(got) {
		t.Errorf("测试文本 测试文本 测试文本 测试文本 测试文本 测试文本: %q", got)
	}
	if !strings.Contains(got, "测试文本 测试文本 测试文本 测试文本") {
		t.Errorf("Err 测试文本 nil 测试文本 '测试文本 测试文本 测试文本 测试文本' 测试文本 测试文本, 测试文本 测试文本 %q", got)
	}
}

// 说明。
// 说明。
// 说明。
func TestHTTPStatusErrorsLocalized(t *testing.T) {
	t.Setenv(AllowLocalTargetsEnv, "1")
	cases := []struct {
		code int
		want string
	}{
		{http.StatusForbidden, "测试文本 测试文本"},
		{http.StatusInternalServerError, "测试文本 测试文本"},
		{http.StatusTooManyRequests, "测试文本 测试文本"},
	}
	for _, tc := range cases {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(tc.code)
			_, _ = w.Write([]byte("body"))
		}))
		_, err := doJSON(context.Background(), http.MethodPost, srv.URL, nil, map[string]any{"a": 1})
		srv.Close()
		if err == nil {
			t.Fatalf("HTTP %d 测试文本 测试文本 测试文本", tc.code)
		}
		assertKorean(t, "doJSON(HTTP status)", err.Error())
		if !strings.Contains(err.Error(), tc.want) {
			t.Errorf("HTTP %d 测试文本 %q 测试文本 测试文本 测试文本, 测试文本 测试文本 %q", tc.code, tc.want, err.Error())
		}
	}
}
