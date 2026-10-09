package server

import "testing"

// TestSideQuestionAPIErrorsLocalized guards F21: every user-facing side-question
// 说明。
// exchange status (e.Error) must be Korean — Hangul present, no Chinese Han.
// 说明。
// both layers are localized together to avoid a mixed-language panel. Reuses
// assertChineseMessage (F3a).
func TestSideQuestionAPIErrorsLocalized(t *testing.T) {
	for _, c := range []struct{ label, msg string }{
		{"ctx_not_saved", sideErrCtxNotSaved},
		{"model_config_changed", sideErrModelConfigChanged},
		{"service_unavailable", sideErrServiceUnavailable},
		{"task_archived", sideErrTaskArchived},
		{"worker_deleted", sideErrWorkerDeleted},
		{"bad_question", sideErrBadQuestion},
		{"task_archiving", sideErrTaskArchiving},
		{"request_id_reused", sideErrRequestIDReused},
		{"no_snapshot", sideErrNoSnapshot},
		{"concurrency_limit", sideErrConcurrencyLimit},
		{"answer_stopped", sideErrAnswerStopped},
		{"answer_timeout", sideErrAnswerTimeout},
	} {
		assertChineseMessage(t, c.label, c.msg)
	}
}
