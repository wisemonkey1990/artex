package selfupdate

// 说明。
// 说明。
// 说明。
// 说明。
// 说明。
// 说明。

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
		t.Fatalf("%s: 测试文本 测试文本 nil 测试文本", label)
	}
	msg := err.Error()
	if hasHangul(msg) {
		t.Errorf("%s: Korean text remains: %q", label, msg)
	}
	if !hasHanzi(msg) {
		t.Errorf("%s: Chinese text is missing: %q", label, msg)
	}
}

// 说明。
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

// 说明。
// 文案语言不会检查。这里确认两种拒绝路径均返回中文错误。
func TestCheckURLErrorsLocalized(t *testing.T) {
	nonHTTPS := checkURL(mustParse(t, "http://api.github.com/x"))
	assertChineseMessage(t, "non-https", nonHTTPS)
	if !strings.Contains(nonHTTPS.Error(), "HTTPS") {
		t.Errorf("non-https 测试文本 HTTPS 测试文本 测试文本: %q", nonHTTPS.Error())
	}

	nonGitHub := checkURL(mustParse(t, "https://evil.example.com/x"))
	assertChineseMessage(t, "non-github-host", nonGitHub)
	if !strings.Contains(nonGitHub.Error(), "GitHub") {
		t.Errorf("non-github 测试文本 GitHub 测试文本 测试文本: %q", nonGitHub.Error())
	}
}

// 说明。
// 状态、解析失败和缺少 tag）均返回中文错误。使用伪 transport，不访问网络。
func TestFetchLatestErrorsLocalized(t *testing.T) {
	ctx := context.Background()

	_, errRate := FetchLatest(ctx, stubClient(http.StatusForbidden, ""))
	assertChineseMessage(t, "rate-limited", errRate)
	if !strings.Contains(errRate.Error(), "60") {
		t.Errorf("测试文本 测试文本 测试文本 60 测试文本 测试文本: %q", errRate.Error())
	}

	_, errNF := FetchLatest(ctx, stubClient(http.StatusNotFound, ""))
	assertChineseMessage(t, "not-found", errNF)

	_, errStatus := FetchLatest(ctx, stubClient(http.StatusInternalServerError, ""))
	assertChineseMessage(t, "bad-status", errStatus)
	if !strings.Contains(errStatus.Error(), "GitHub") {
		t.Errorf("测试文本 测试文本 GitHub 测试文本 测试文本: %q", errStatus.Error())
	}

	_, errJSON := FetchLatest(ctx, stubClient(http.StatusOK, "{not json"))
	assertChineseMessage(t, "bad-json", errJSON)

	_, errTag := FetchLatest(ctx, stubClient(http.StatusOK, "{}"))
	assertChineseMessage(t, "missing-tag", errTag)
	if !strings.Contains(errTag.Error(), "tag") {
		t.Errorf("tag 测试文本 测试文本 tag 测试文本 测试文本: %q", errTag.Error())
	}
}

// 说明。
// 下载前会因中文错误而停止（FindAsset 失败，不访问网络）。运行环境也可能先被
// checkWritable 拦截；该错误同样是中文，因此这里只检查语言。
func TestStageMissingAssetLocalized(t *testing.T) {
	rel := &Release{TagName: "v9.9.9"}
	err := Stage(context.Background(), &http.Client{}, rel, "0.0.1", nil)
	assertChineseMessage(t, "stage-missing-asset", err)
}
