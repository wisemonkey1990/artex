package server

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Autumn-27/artex/db"
)

// 说明。
// 说明。
// 说明。
// 说明。
// 说明。

// 说明。
// 说明。
// 说明。
// 说明。
// 说明。
func TestCustomToolErrorConstantsLocalized(t *testing.T) {
	for _, c := range []struct{ name, msg string }{
		{"keyFormat", errCustomToolKeyFormat},
		{"kindInvalid", errCustomToolKindInvalid},
		{"httpSchemaRequired", errCustomToolHTTPSchemaRequired},
		{"keyExists", errCustomToolKeyExists},
		{"editCustomOnly", errCustomToolEditCustomOnly},
		{"badBody", errCustomToolBadBody},
		{"shellNoExec", errCustomToolShellNoExec},
		{"unknownKindPrefix", errCustomToolUnknownKindPrefix},
	} {
		assertChineseMessage(t, "customtool."+c.name, c.msg)
	}
}

// 说明。
// 说明。
// 说明。
// 说明。
func newCustomToolServer() *Server {
	return &Server{m: &Manager{pg: &db.DB{}}}
}

// 说明。
// 说明。
// 说明。
func TestCustomToolCreateValidationLocalized(t *testing.T) {
	cases := []struct {
		name, body, want string
	}{
		{"bad-key", `{"key":"BadKey","kind":"command"}`, errCustomToolKeyFormat},
		{"bad-kind", `{"key":"goodkey","kind":"bogus"}`, errCustomToolKindInvalid},
		{"http-no-schema", `{"key":"goodkey","kind":"http"}`, errCustomToolHTTPSchemaRequired},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			s := newCustomToolServer()
			req := httptest.NewRequest(http.MethodPost, "/api/tools/custom", strings.NewReader(c.body))
			rec := httptest.NewRecorder()
			s.pgCreateCustomTool(rec, req)
			if rec.Code != 400 && rec.Code != 409 {
				t.Fatalf("测试文本 测试文本 测试文本 测试文本 测试文本 %d (测试文本 %q)", rec.Code, rec.Body.String())
			}
			got := decodeErrorField(t, rec.Body.Bytes())
			if got != c.want {
				t.Fatalf("测试文本 测试文本 测试文本: got %q want %q", got, c.want)
			}
			assertChineseMessage(t, "customtool.create."+c.name, got)
		})
	}
}

// 说明。
// 说明。
func TestCustomToolTestRunValidationLocalized(t *testing.T) {
	cases := []struct {
		name, body, want string
	}{
		{"bad-body", `not-json`, errCustomToolBadBody},
		{"shell", `{"kind":"shell"}`, errCustomToolShellNoExec},
		{"unknown-kind", `{"kind":"bogus"}`, errCustomToolUnknownKindPrefix + "bogus"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			s := newCustomToolServer()
			req := httptest.NewRequest(http.MethodPost, "/api/tools/custom/test", strings.NewReader(c.body))
			rec := httptest.NewRecorder()
			s.pgTestCustomTool(rec, req)
			if rec.Code != 400 {
				t.Fatalf("测试文本 测试文本 400 测试文本 测试文本 %d (测试文本 %q)", rec.Code, rec.Body.String())
			}
			got := decodeErrorField(t, rec.Body.Bytes())
			if got != c.want {
				t.Fatalf("测试文本 测试文本 测试文本: got %q want %q", got, c.want)
			}
			// 说明。
			assertChineseMessage(t, "customtool.test."+c.name, errCustomToolUnknownKindPrefix)
		})
	}
}
