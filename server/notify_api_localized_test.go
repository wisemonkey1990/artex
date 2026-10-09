package server

import (
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/Autumn-27/artex/db"
)

// 说明。
// 说明。

// 说明。
// 说明。
func TestNotifyErrorConstantsLocalized(t *testing.T) {
	cases := map[string]string{
		"notifyErrBadJSON":         notifyErrBadJSON,
		"notifyErrKindInvalidFmt":  fmt.Sprintf(notifyErrKindInvalidFmt, "email / webhook"),
		"notifyErrNameMissing":     notifyErrNameMissing,
		"notifyErrNameEmpty":       notifyErrNameEmpty,
		"notifyErrModeInvalid":     notifyErrModeInvalid,
		"notifyErrRateNegative":    notifyErrRateNegative,
		"notifyErrChannelID":       notifyErrChannelID,
		"notifyErrKindUnregFmt":    fmt.Sprintf(notifyErrKindUnregFmt, "bogus"),
		"notifyErrDeliveryID":      notifyErrDeliveryID,
		"notifyErrChannelNotFound": notifyErrChannelNotFound,
		"notifyTestName":           notifyTestName,
		"notifyTestClass":          notifyTestClass,
		"notifyTestSummary":        notifyTestSummary,
	}
	for label, msg := range cases {
		assertChineseMessage(t, label, msg)
	}
}

// 说明。
// 说明。
// 说明。
func TestNotifyChannelLookupErrLocalized(t *testing.T) {
	rec := httptest.NewRecorder()
	notifyChannelLookupErr(rec, db.ErrNotificationChannelNotFound)
	if rec.Code != 404 {
		t.Fatalf("测试文本 测试文本 = %d, 测试文本 = 404 (测试文本 %q)", rec.Code, rec.Body.String())
	}
	got := decodeErrorField(t, rec.Body.Bytes())
	if got != notifyErrChannelNotFound {
		t.Fatalf("测试文本 测试文本 = %q, 测试文本 = %q", got, notifyErrChannelNotFound)
	}
	assertChineseMessage(t, "channel_not_found", got)
}

// 说明。
// 说明。
func TestNotifyTestMessageLocalized(t *testing.T) {
	const base = "https://artex.example.test"
	msg := notifyTestMessage(base)
	if len(msg.Items) != 1 {
		t.Fatalf("测试文本 测试文本 测试文本 测试文本 = %d, 测试文本 = 1", len(msg.Items))
	}
	item := msg.Items[0]
	assertChineseMessage(t, "test.name", item.Name)
	assertChineseMessage(t, "test.class", item.VulnClass)
	assertChineseMessage(t, "test.summary", item.Summary)
	if msg.HomeURL != base || item.DetailURL != base {
		t.Fatalf("测试文本 测试文本 测试文本: HomeURL=%q DetailURL=%q (测试文本 %q)", msg.HomeURL, item.DetailURL, base)
	}
}
