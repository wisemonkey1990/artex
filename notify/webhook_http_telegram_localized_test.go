package notify

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

// 이 파일은 F4 ④⑤(telegram·webhook·공용 HTTP 전송 계층) 중국어화를 회귀로부터
// 지킨다. PR#1(60b62f0)이 telegram.go·webhook.go·http.go·dingtalk.go(validateHTTPURL)의
// 사용자 노출 문구를 한국어에서 중국어로 바꿨으므로, 여기서는 그 반대 방향— 한글이
// 남아 있으면 잡아내고 중국어 한자가 있어야 통과하는 쪽으로 단언을 뒤집는다.
// 헬퍼 hasHan·hasHangul 은 notify_localized_test.go 에 있다(같은 패키지, 그대로 재사용).
// assertChineseMessage 는 이 파일 전용이다 — notify_localized_test.go 의 assertKorean 은
// 아직 한국어 단언을 쓰는 다른 파일들(dingtalk_feishu_wecom·html_email·markdown·mask_filter)
// 이 참조하므로 건드리지 않는다.
//
// 데이터(제목·유형·자산·개요)는 전부 ASCII 로 둔다. 그래야 "출력 전체에 한글 0"
// 이라는 단언이 콘텐츠가 아니라 골격 라벨의 회귀만 정확히 포착한다.

// assertChineseMessage 断言字符串含有中文汉字且不含韩文谚字——与 assertKorean 方向相反，
// 用来校验 F4 ④⑤ 完成中文化之后的用户可见文案。
func assertChineseMessage(t *testing.T, where, got string) {
	t.Helper()
	if hasHangul(got) {
		t.Errorf("%s: 韩文谚字仍然残留: %q", where, got)
	}
	if !hasHan(got) {
		t.Errorf("%s: 缺少中文汉字: %q", where, got)
	}
}

// TestTelegramItemLabelsLocalized 는 단건 Telegram 메시지의 모든 라벨 분기를 켠 뒤
// 골격이 중국어이고 한글이 없음을 확인한다.
func TestTelegramItemLabelsLocalized(t *testing.T) {
	m := Message{
		HomeURL: "https://example.com/panel",
		Items: []Item{{
			Name:       "sqli-login", // Title() = Name
			VulnClass:  "injection",  // Title() 과 달라야 유형 줄이 렌더됨
			Severity:   "high",       // SeverityLabel → 높음(한글)
			Summary:    "login form is injectable",
			Assets:     []string{"host-a.example.com"},
			DetailURL:  "https://example.com/f/1",
			FromStatus: "pending", // IsStatusChange()=true → 상태 변경 줄
			ToStatus:   "fixed",
		}},
	}
	out, kept := telegramHTML(m)
	if kept != 1 {
		t.Fatalf("단건 메시지는 1 을 보고해야 합니다, 받은 값 %d", kept)
	}
	if hasHangul(out) {
		t.Errorf("Telegram 본문에 한글이 남아 있습니다:\n%s", out)
	}
	for _, want := range []string{"状态变更", "类型", "资产", "概述", "查看详情"} {
		if !strings.Contains(out, want) {
			t.Errorf("Telegram 본문에 %q 라벨이 없습니다:\n%s", want, out)
		}
	}
}

// TestTelegramBatchTitleLocalized 는 배치(digest) 제목의 넘침·시간창 분기가
// 중국어로 렌더되는지 확인한다. telegramBatchTitle 과 글자까지 맞춰야 채널 간 혼재가
// 생기지 않는다.
func TestTelegramBatchTitleLocalized(t *testing.T) {
	items := []Item{{Name: "a", Severity: "high"}, {Name: "b", Severity: "low"}}

	// 넘침 있음 + 시간창: 요약·총건수·앞 N건·나머지·다음 메시지·최근 N분간.
	over := telegramBatchTitle(Message{WindowMinutes: 30}, items, 5)
	assertChineseMessage(t, "telegramBatchTitle(넘침)", over)
	for _, want := range []string{"最近 30 分钟", "漏洞摘要 · 共 5 项", "仅显示前 2 项", "其余 3 项", "下一条消息"} {
		if !strings.Contains(over, want) {
			t.Errorf("넘침 제목에 %q 가 없습니다: %q", want, over)
		}
	}

	// 넘침 없음 + 시간창 없음: 다음 메시지 안내가 없어야 한다.
	full := telegramBatchTitle(Message{}, items, 2)
	assertChineseMessage(t, "telegramBatchTitle(전량)", full)
	if !strings.Contains(full, "漏洞摘要 · 共 2 项") {
		t.Errorf("전량 제목에 '漏洞摘要 · 共 2 项' 이 있어야 합니다: %q", full)
	}
	if strings.Contains(full, "下一条消息") {
		t.Errorf("넘치지 않았는데 '下一条消息' 안내가 들어갔습니다: %q", full)
	}
}

// TestTelegramValidateLocalized 는 필수 필드 누락 오류가 중국어이고 어떤 필드가
// 비었는지 알려 주는지 확인한다. 기술 용어 Bot Token·Chat ID 는 그대로 보존한다.
func TestTelegramValidateLocalized(t *testing.T) {
	if err := (telegramChannel{}).Validate(map[string]any{}); err == nil {
		t.Fatal("bot_token 누락은 검증 실패여야 합니다")
	} else {
		assertChineseMessage(t, "telegram Validate(bot_token)", err.Error())
		if !strings.Contains(err.Error(), "Bot Token") {
			t.Errorf("Bot Token 누락을 알려줘야 합니다: %q", err.Error())
		}
	}
	if err := (telegramChannel{}).Validate(map[string]any{"bot_token": "t"}); err == nil {
		t.Fatal("chat_id 누락은 검증 실패여야 합니다")
	} else {
		assertChineseMessage(t, "telegram Validate(chat_id)", err.Error())
		if !strings.Contains(err.Error(), "Chat ID") {
			t.Errorf("Chat ID 누락을 알려줘야 합니다: %q", err.Error())
		}
	}
}

// TestWebhookValidateLocalized 는 범용 Webhook 설정 검증 오류 세 가지가 중국어이고
// 각 오류가 문제의 원인을 알려 주는지 확인한다.
func TestWebhookValidateLocalized(t *testing.T) {
	if err := (webhookChannel{}).Validate(map[string]any{}); err == nil {
		t.Fatal("url 누락은 검증 실패여야 합니다")
	} else {
		assertChineseMessage(t, "webhook Validate(url)", err.Error())
		if !strings.Contains(err.Error(), "目标 URL") {
			t.Errorf("目标 URL 누락을 알려줘야 합니다: %q", err.Error())
		}
	}
	// 지원하지 않는 메서드: 중국어 + 허용 목록 노출.
	if err := (webhookChannel{}).Validate(map[string]any{"url": "https://example.com/hook", "method": "DELETE"}); err == nil {
		t.Fatal("DELETE 는 검증 실패여야 합니다")
	} else {
		assertChineseMessage(t, "webhook Validate(method)", err.Error())
		if !strings.Contains(err.Error(), "GET、POST、PUT 或 PATCH") {
			t.Errorf("허용 메서드 목록을 알려줘야 합니다: %q", err.Error())
		}
	}
	// 템플릿 문법 오류: 중국어(래핑된 원인은 Go 템플릿 오류라 한글이 없다).
	if err := (webhookChannel{}).Validate(map[string]any{"url": "https://example.com/hook", "body_template": "{{"}); err == nil {
		t.Fatal("깨진 템플릿은 검증 실패여야 합니다")
	} else {
		assertChineseMessage(t, "webhook Validate(template)", err.Error())
		if !strings.Contains(err.Error(), "模板") {
			t.Errorf("템플릿 오류임을 알려줘야 합니다: %q", err.Error())
		}
	}
}

// TestValidateHTTPURLLocalized 는 다섯 채널(telegram·webhook·dingtalk·feishu·wecom)이
// 공유하는 URL 검증 헬퍼의 오류가 중국어인지 확인한다. 이 헬퍼가 한국어로 남아 있으면
// 중국어 접두("目标 URL 无效: ...")와 섞여 반한반중이 된다.
func TestValidateHTTPURLLocalized(t *testing.T) {
	// 지원하지 않는 스킴.
	if err := validateHTTPURL("ftp://example.com/x"); err == nil {
		t.Fatal("ftp 스킴은 거부되어야 합니다")
	} else {
		assertChineseMessage(t, "validateHTTPURL(scheme)", err.Error())
		if !strings.Contains(err.Error(), "http") {
			t.Errorf("지원 스킴(http/https)을 알려줘야 합니다: %q", err.Error())
		}
	}
	// 호스트 이름 없음.
	if err := validateHTTPURL("http://"); err == nil {
		t.Fatal("호스트 없는 주소는 거부되어야 합니다")
	} else {
		assertChineseMessage(t, "validateHTTPURL(host)", err.Error())
	}
}

// TestHTTPLocalTargetBlockedLocalized 는 로컬/링크 로컬 주소 차단 오류가 중국어이고
// 우회 방법(환경 변수)을 알려 주는지 확인한다.
func TestHTTPLocalTargetBlockedLocalized(t *testing.T) {
	t.Setenv(AllowLocalTargetsEnv, "") // 명시적으로 꺼서 차단 경로를 탄다
	err := blockInternalDial("tcp", "127.0.0.1:25", nil)
	if err == nil {
		t.Fatal("로컬 주소는 기본적으로 차단되어야 합니다")
	}
	assertChineseMessage(t, "blockInternalDial", err.Error())
	if !strings.Contains(err.Error(), AllowLocalTargetsEnv) {
		t.Errorf("차단 오류는 우회 방법(%s)을 알려줘야 합니다: %q", AllowLocalTargetsEnv, err.Error())
	}
}

// TestHTTPRedactHelpersLocalized 는 주소 탈감(redact) 자리표시자와 전송 오류의
// "알 수 없는 오류" 분기가 중국어인지 확인한다.
func TestHTTPRedactHelpersLocalized(t *testing.T) {
	if got := redactRequestTarget(""); !strings.Contains(got, "无法解析") {
		t.Errorf("해석 불가 주소는 중국어 자리표시자여야 합니다, 받은 값 %q", got)
	}
	// *url.Error 의 Err 가 nil 인 분기.
	got := redactTransportError(&url.Error{Op: "Get", URL: "http://api.example", Err: nil})
	if hasHangul(got) {
		t.Errorf("전송 오류 탈감 결과에 한글이 남았습니다: %q", got)
	}
	if !strings.Contains(got, "未知错误") {
		t.Errorf("Err 가 nil 이면 '未知错误' 여야 합니다, 받은 값 %q", got)
	}
}

// TestHTTPStatusErrorsLocalized 는 실제 doJSON 왕복으로 상태코드 분류(거부/서버
// 오류/제한)의 사용자 노출 오류가 중국어인지 확인한다. 127.0.0.1 httptest 로의
// 전송은 기본 차단이라 이 테스트에서만 명시적으로 허용한다.
func TestHTTPStatusErrorsLocalized(t *testing.T) {
	t.Setenv(AllowLocalTargetsEnv, "1")
	cases := []struct {
		code int
		want string
	}{
		{http.StatusForbidden, "拒绝了请求"},
		{http.StatusInternalServerError, "发生错误"},
		{http.StatusTooManyRequests, "限流或响应超时"},
	}
	for _, tc := range cases {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(tc.code)
			_, _ = w.Write([]byte("body"))
		}))
		_, err := doJSON(context.Background(), http.MethodPost, srv.URL, nil, map[string]any{"a": 1})
		srv.Close()
		if err == nil {
			t.Fatalf("HTTP %d 는 오류여야 합니다", tc.code)
		}
		assertChineseMessage(t, "doJSON(HTTP status)", err.Error())
		if !strings.Contains(err.Error(), tc.want) {
			t.Errorf("HTTP %d 오류에 %q 가 있어야 합니다, 받은 값 %q", tc.code, tc.want, err.Error())
		}
	}
}
