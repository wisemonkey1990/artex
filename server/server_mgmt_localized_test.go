package server

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// TestServerMgmtConstantsLocalized pins that every user-facing error/response
// literal extracted from server_mgmt.go is Korean (Hangul present, no Chinese
// Han). Reverting any one to Chinese makes assertChineseMessage fail on the Han
// ideograph. assertChineseMessage / decodeErrorField are reused from the F3 suite
// (same package server).
func TestServerMgmtConstantsLocalized(t *testing.T) {
	cases := map[string]string{
		"errMgmtPGUnavailable":      errMgmtPGUnavailable,
		"errMgmtTaskIDInvalid":      errMgmtTaskIDInvalid,
		"errMgmtTaskDeleting":       errMgmtTaskDeleting,
		"errMgmtTaskHasAgents":      errMgmtTaskHasAgents,
		"errMgmtAgentKeyFormat":     errMgmtAgentKeyFormat,
		"errMgmtNameEmpty":          errMgmtNameEmpty,
		"errMgmtAgentKeyExists":     errMgmtAgentKeyExists,
		"errMgmtBuiltinNoEditMeta":  errMgmtBuiltinNoEditMeta,
		"errMgmtBuiltinNoDelete":    errMgmtBuiltinNoDelete,
		"errMgmtLLMProfileIDFormat": errMgmtLLMProfileIDFormat,
		"errMgmtLLMProfileInvalid":  errMgmtLLMProfileInvalid,
		"errMgmtNoBuiltinPrompt":    errMgmtNoBuiltinPrompt,
		"errMgmtToolNotFound":       errMgmtToolNotFound,
		"errMgmtNotBuiltinTool":     errMgmtNotBuiltinTool,
		"errMgmtMCPNotFound":        errMgmtMCPNotFound,
		"errMgmtToolDiscoverFail":   errMgmtToolDiscoverFail,
		"errMgmtSkillNameInvalid":   errMgmtSkillNameInvalid,
		"errMgmtSkillNoFile":        errMgmtSkillNoFile,
		"errMgmtSkillNoMDInZip":     errMgmtSkillNoMDInZip,
		"errMgmtSkillReadMDFail":    errMgmtSkillReadMDFail,
		"errMgmtSkillNamePre":       errMgmtSkillNamePre,
		"errMgmtSkillNamePost":      errMgmtSkillNamePost,
		"errMgmtSkillExistsPre":     errMgmtSkillExistsPre,
		"errMgmtSkillExistsPost":    errMgmtSkillExistsPost,
		"errMgmtSkillZipBadPath":    errMgmtSkillZipBadPath,
		"errMgmtSkillZipTooMany":    errMgmtSkillZipTooMany,
		"errMgmtSkillFileTooLarge":  errMgmtSkillFileTooLarge,
		"errMgmtSkillZipTooLarge":   errMgmtSkillZipTooLarge,
		"errMgmtSkillNoMDAfter":     errMgmtSkillNoMDAfter,
		"errMgmtSkillInstallFail":   errMgmtSkillInstallFail,
		"errMgmtLLMActiveDelete":    errMgmtLLMActiveDelete,
		"errMgmtLLMRefChanged":      errMgmtLLMRefChanged,
		"errMgmtLLMRefTimeout":      errMgmtLLMRefTimeout,
		"errMgmtNoAPIKey":           errMgmtNoAPIKey,
		"errMgmtBuildReqFail":       errMgmtBuildReqFail,
		"errMgmtReqFail":            errMgmtReqFail,
		"errMgmtParseRespFail":      errMgmtParseRespFail,
		"errMgmtNoModelList":        errMgmtNoModelList,
		"errMgmtTmplSyntax":         errMgmtTmplSyntax,
		"errMgmtTmplVarPre":         errMgmtTmplVarPre,
		"errMgmtTmplVarPost":        errMgmtTmplVarPost,
		// 说明。
		"errMgmtAPIReturned": fmt.Sprintf(errMgmtAPIReturned, 404, "detail"),
	}
	for label, msg := range cases {
		assertChineseMessage(t, label, msg)
	}
}

// TestServerMgmtAPIReturnedFormat guards that the model-list probe format string
// keeps its status-code and body placeholders (a %d→wrong-verb edit would drop
// them).
func TestServerMgmtAPIReturnedFormat(t *testing.T) {
	msg := fmt.Sprintf(errMgmtAPIReturned, 404, "boom-body")
	assertChineseMessage(t, "errMgmtAPIReturned", msg)
	if !strings.Contains(msg, "404") || !strings.Contains(msg, "boom-body") {
		t.Fatalf("测试文本 测试文本 测试文本 测试文本 测试文本 测试文本 测试文本: %q", msg)
	}
}

// TestServerMgmtSkillExistsSentinel guards the cross-stack contract: the server's
// duplicate-skill message must carry the "已存在" marker that
// web/src/app/(main)/system/skills/page.tsx greps (msg.includes("已存在")) to
// switch into the overwrite-confirm flow. Drift here silently disables overwrite.
func TestServerMgmtSkillExistsSentinel(t *testing.T) {
	msg := errMgmtSkillExistsPre + "my-skill" + errMgmtSkillExistsPost
	assertChineseMessage(t, "skillExists", msg)
	if !strings.Contains(msg, "已存在") {
		t.Fatalf("测试文本 测试文本 测试文本 '已存在' 测试文本 测试文本: %q", msg)
	}
	if !strings.Contains(msg, "my-skill") {
		t.Fatalf("测试文本 测试文本 测试文本 测试文本 测试文本: %q", msg)
	}
}

// TestServerMgmtPGGate503Localized exercises the real shared DB gate every
// management handler passes through: with no PostgreSQL handle, pg() writes 503
// and returns nil. The 503 body must be Korean.
func TestServerMgmtPGGate503Localized(t *testing.T) {
	s := &Server{m: &Manager{}} // pg == nil
	rec := httptest.NewRecorder()
	if got := s.pg(rec); got != nil {
		t.Fatal("pg() 测试文本 DB 测试文本 测试文本 nil 测试文本 测试文本 测试文本")
	}
	if rec.Code != 503 {
		t.Fatalf("测试文本 测试文本 503 测试文本 测试文本 %d", rec.Code)
	}
	assertChineseMessage(t, "pg.503", decodeErrorField(t, rec.Body.Bytes()))
}

// TestServerMgmtDeleteTaskBadIDLocalized drives pgDeleteTask through its first
// guard (non-numeric id → 400) without any DB or engine.
func TestServerMgmtDeleteTaskBadIDLocalized(t *testing.T) {
	s := &Server{}
	req := httptest.NewRequest(http.MethodDelete, "/api/tasks/abc", nil)
	req.SetPathValue("id", "abc")
	rec := httptest.NewRecorder()
	s.pgDeleteTask(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("测试文本 测试文本 400 测试文本 测试文本 %d (测试文本 %s)", rec.Code, rec.Body.Bytes())
	}
	got := decodeErrorField(t, rec.Body.Bytes())
	assertChineseMessage(t, "deleteTask.badID", got)
	if got != errMgmtTaskIDInvalid {
		t.Fatalf("errMgmtTaskIDInvalid 测试文本 测试文本 测试文本: %q", got)
	}
}

// TestServerMgmtListModelsNoKeyLocalized drives pgListModels with no api key and
// no profile id, so it returns {"ok":false,"error":...} before touching the DB.
func TestServerMgmtListModelsNoKeyLocalized(t *testing.T) {
	s := &Server{}
	req := httptest.NewRequest(http.MethodPost, "/api/llm/models", strings.NewReader("{}"))
	rec := httptest.NewRecorder()
	s.pgListModels(rec, req)
	if rec.Code != 200 {
		t.Fatalf("测试文本 测试文本 200 测试文本 测试文本 %d", rec.Code)
	}
	got := decodeErrorField(t, rec.Body.Bytes())
	assertChineseMessage(t, "listModels.noKey", got)
	if got != errMgmtNoAPIKey {
		t.Fatalf("errMgmtNoAPIKey 测试文本 测试文本 测试文本: %q", got)
	}
}

// TestServerMgmtValidateTemplateLocalized exercises the prompt-template editor
// validator (pure function, no DB): both the parse-error and the
// disallowed-variable branches must return Korean.
func TestServerMgmtValidateTemplateLocalized(t *testing.T) {
	if msg := validateTemplate("{{", nil); msg == "" {
		t.Fatal("测试文本 测试文本 测试文本 测试文本 测试文本 测试文本")
	} else {
		assertChineseMessage(t, "validateTemplate.syntax", msg)
	}
	msg := validateTemplate("{{.Bogus}}", nil)
	if msg == "" {
		t.Fatal("允许 测试文本 测试文本 测试文本 测试文本 测试文本 测试文本 测试文本")
	}
	assertChineseMessage(t, "validateTemplate.var", msg)
	if !strings.Contains(msg, "Bogus") {
		t.Fatalf("测试文本 测试文本 测试文本 测试文本 测试文本: %q", msg)
	}
}

// TestServerMgmtGlobalPromptVarsLocalized pins the template-variable catalog help
// text (rendered in the prompt editor's variable palette) as Korean.
func TestServerMgmtGlobalPromptVarsLocalized(t *testing.T) {
	if len(globalPromptVars) == 0 {
		t.Fatal("globalPromptVars 测试文本 测试文本 测试文本")
	}
	for _, v := range globalPromptVars {
		assertChineseMessage(t, "globalPromptVar."+v.Name, v.Description)
	}
}
