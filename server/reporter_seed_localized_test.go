package server

import "testing"

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
func TestReporterSeedLabelsLocalized(t *testing.T) {
	assertChineseMessage(t, "reporter_name", reporterAgentName)
	assertChineseMessage(t, "reporter_description", reporterAgentDescription)
}
