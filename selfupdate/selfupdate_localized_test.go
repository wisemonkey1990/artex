package selfupdate

// 셀프 업데이트 경로의 사용자 노출 진행·오류 문구가 중국어로 나오는지 고정하는
// 회귀 테스트(F20). update.go 의 updateHub 가 Stage 의 진행 콜백 메시지를 SSE
// progress.message 로, Stage/FetchLatest/Rollback 의 반환 오류를 progress.error·
// writeErr 본문·boot_notice 로 그대로 프런트엔드(update-card)에 노출하므로, 이 문구가
// 한국어로 되돌아가면 사용자 화면의 문구가 섞인다. 로그(log.Printf)·주석은 범위 밖이고,
// 여기서는 한글이 없고 중국어 한자가 포함되는지만 단언한다.

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"
)

func hasHangul(s string) bool {
	for _, r := range s {
		if r >= 0xAC00 && r <= 0xD7A3 {
			return true
		}
	}
	return false
}

func hasHanzi(s string) bool {
	for _, r := range s {
		if r >= 0x4E00 && r <= 0x9FFF {
			return true
		}
	}
	return false
}

// assertChineseMessage 会断言错误非空、不含韩文且包含中文汉字。
// %w 包装的底层标准库英文错误可以保留。
func assertChineseMessage(t *testing.T, label string, err error) {
	t.Helper()
	if err == nil {
		t.Fatalf("%s: 오류를 기대했으나 nil 이었습니다", label)
	}
	msg := err.Error()
	if hasHangul(msg) {
		t.Errorf("%s: Korean text remains: %q", label, msg)
	}
	if !hasHanzi(msg) {
		t.Errorf("%s: Chinese text is missing: %q", label, msg)
	}
}

// rtFunc 로 http.Client 에 가짜 응답을 주입한다. GitHub 에 실제로 접속하지 않는다.
type rtFunc func(*http.Request) (*http.Response, error)

func (f rtFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

func stubClient(status int, body string) *http.Client {
	return &http.Client{Transport: rtFunc(func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: status,
			Body:       io.NopCloser(strings.NewReader(body)),
			Header:     make(http.Header),
		}, nil
	})}
}

// TestCheckURLErrorsLocalized: 기존 TestCheckURLRejectsNonGitHub 는 거부 여부만 보고
// 文案语言不会检查。这里确认两种拒绝路径均返回中文错误。
func TestCheckURLErrorsLocalized(t *testing.T) {
	nonHTTPS := checkURL(mustParse(t, "http://api.github.com/x"))
	assertChineseMessage(t, "non-https", nonHTTPS)
	if !strings.Contains(nonHTTPS.Error(), "HTTPS") {
		t.Errorf("non-https 오류에 HTTPS 표기가 없습니다: %q", nonHTTPS.Error())
	}

	nonGitHub := checkURL(mustParse(t, "https://evil.example.com/x"))
	assertChineseMessage(t, "non-github-host", nonGitHub)
	if !strings.Contains(nonGitHub.Error(), "GitHub") {
		t.Errorf("non-github 오류에 GitHub 표기가 없습니다: %q", nonGitHub.Error())
	}
}

// TestFetchLatestErrorsLocalized: 버전 확인 경로의 네 가지 오류(요청 제한·미발행·기타
// 状态、解析失败和缺少 tag）均返回中文错误。使用伪 transport，不访问网络。
func TestFetchLatestErrorsLocalized(t *testing.T) {
	ctx := context.Background()

	_, errRate := FetchLatest(ctx, stubClient(http.StatusForbidden, ""))
	assertChineseMessage(t, "rate-limited", errRate)
	if !strings.Contains(errRate.Error(), "60") {
		t.Errorf("요청 제한 오류에 60 표기가 없습니다: %q", errRate.Error())
	}

	_, errNF := FetchLatest(ctx, stubClient(http.StatusNotFound, ""))
	assertChineseMessage(t, "not-found", errNF)

	_, errStatus := FetchLatest(ctx, stubClient(http.StatusInternalServerError, ""))
	assertChineseMessage(t, "bad-status", errStatus)
	if !strings.Contains(errStatus.Error(), "GitHub") {
		t.Errorf("상태 오류에 GitHub 표기가 없습니다: %q", errStatus.Error())
	}

	_, errJSON := FetchLatest(ctx, stubClient(http.StatusOK, "{not json"))
	assertChineseMessage(t, "bad-json", errJSON)

	_, errTag := FetchLatest(ctx, stubClient(http.StatusOK, "{}"))
	assertChineseMessage(t, "missing-tag", errTag)
	if !strings.Contains(errTag.Error(), "tag") {
		t.Errorf("tag 누락 오류에 tag 표기가 없습니다: %q", errTag.Error())
	}
}

// TestStageMissingAssetLocalized: 릴리스에 현재 플랫폼 패키지가 없으면 Stage 가
// 下载前会因中文错误而停止（FindAsset 失败，不访问网络）。运行环境也可能先被
// checkWritable 拦截；该错误同样是中文，因此这里只检查语言。
func TestStageMissingAssetLocalized(t *testing.T) {
	rel := &Release{TagName: "v9.9.9"} // Assets 비어 있음
	err := Stage(context.Background(), &http.Client{}, rel, "0.0.1", nil)
	assertChineseMessage(t, "stage-missing-asset", err)
}
