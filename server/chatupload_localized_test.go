package server

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
)

// 说明。
// 说明。
// 说明。

// 说明。
// 说明。
// 说明。
func TestChatUploadErrorConstantsLocalized(t *testing.T) {
	cases := map[string]string{
		"scope_invalid": errChatUploadScopeInvalid,
		"bad_id":        errChatUploadBadID,
		"task_deleting": errChatUploadTaskDeleting,
		"mkdir":         errChatUploadMkdir,
		"parse":         errChatUploadParse,
		"no_file":       errChatUploadNoFile,
		"save_failed":   errChatUploadSaveFailed,
	}
	for label, msg := range cases {
		assertChineseMessage(t, label, msg)
	}
}

// 说明。
// 说明。
// 说明。
// 说明。
func TestChatUploadHandlerResponsesLocalized(t *testing.T) {
	// 说明。
	noFileBody := func() (*bytes.Buffer, string) {
		var buf bytes.Buffer
		mw := multipart.NewWriter(&buf)
		_ = mw.WriteField("other", "x")
		_ = mw.Close()
		return &buf, mw.FormDataContentType()
	}

	t.Run("scope-invalid", func(t *testing.T) {
		s := &Server{}
		req := httptest.NewRequest(http.MethodPost, "/api/chat/upload?scope=bogus&id=x", nil)
		rec := httptest.NewRecorder()
		s.chatUpload(rec, req)
		if rec.Code != http.StatusBadRequest {
			t.Fatalf("测试文本 测试文本 = %d, 测试文本 = 400 (测试文本 %q)", rec.Code, rec.Body.String())
		}
		if got := decodeErrorField(t, rec.Body.Bytes()); got != errChatUploadScopeInvalid {
			t.Fatalf("测试文本 测试文本 = %q, 测试文本 = %q", got, errChatUploadScopeInvalid)
		}
		assertChineseMessage(t, "scope-invalid", errChatUploadScopeInvalid)
	})

	t.Run("bad-id", func(t *testing.T) {
		s := &Server{}
		// 说明。
		req := httptest.NewRequest(http.MethodPost, "/api/chat/upload?scope=session&id=../evil", nil)
		rec := httptest.NewRecorder()
		s.chatUpload(rec, req)
		if rec.Code != http.StatusBadRequest {
			t.Fatalf("测试文本 测试文本 = %d, 测试文本 = 400 (测试文本 %q)", rec.Code, rec.Body.String())
		}
		if got := decodeErrorField(t, rec.Body.Bytes()); got != errChatUploadBadID {
			t.Fatalf("测试文本 测试文本 = %q, 测试文本 = %q", got, errChatUploadBadID)
		}
		assertChineseMessage(t, "bad-id", errChatUploadBadID)
	})

	t.Run("no-file", func(t *testing.T) {
		s := &Server{m: &Manager{dir: t.TempDir()}}
		body, ctype := noFileBody()
		req := httptest.NewRequest(http.MethodPost, "/api/chat/upload?scope=session&id=sess1", body)
		req.Header.Set("Content-Type", ctype)
		rec := httptest.NewRecorder()
		s.chatUpload(rec, req)
		if rec.Code != http.StatusBadRequest {
			t.Fatalf("测试文本 测试文本 = %d, 测试文本 = 400 (测试文本 %q)", rec.Code, rec.Body.String())
		}
		if got := decodeErrorField(t, rec.Body.Bytes()); got != errChatUploadNoFile {
			t.Fatalf("测试文本 测试文本 = %q, 测试文本 = %q", got, errChatUploadNoFile)
		}
		assertChineseMessage(t, "no-file", errChatUploadNoFile)
	})

	t.Run("task-deleting", func(t *testing.T) {
		s := &Server{m: &Manager{tasks: map[string]*Task{"t1": {ID: "t1"}}}, engine: &Engine{}}
		s.engine.deleting.Store("t1", true)
		req := httptest.NewRequest(http.MethodPost, "/api/chat/upload?scope=task&id=t1", nil)
		rec := httptest.NewRecorder()
		s.chatUpload(rec, req)
		if rec.Code != http.StatusConflict {
			t.Fatalf("测试文本 测试文本 = %d, 测试文本 = 409 (测试文本 %q)", rec.Code, rec.Body.String())
		}
		if got := decodeErrorField(t, rec.Body.Bytes()); got != errChatUploadTaskDeleting {
			t.Fatalf("测试文本 测试文本 = %q, 测试文本 = %q", got, errChatUploadTaskDeleting)
		}
		assertChineseMessage(t, "task-deleting", errChatUploadTaskDeleting)
	})
}
