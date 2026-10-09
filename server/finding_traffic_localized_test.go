package server

import "testing"

// 说明。
// 说明。
// 说明。
// 说明。
// 说明。
// 说明。
//
// 说明。
// 说明。
// 说明。
func TestFindingTrafficErrorsLocalized(t *testing.T) {
	for _, c := range []struct {
		label, msg string
	}{
		{"inherited_readonly", errFindingTrafficInheritedReadonly},
		{"select_required", errFindingTrafficSelectRequired},
		{"version_required", errFindingTrafficVersionRequired},
		{"binding_ids_required", errFindingTrafficBindingIDsRequired},
	} {
		assertChineseMessage(t, c.label, c.msg)
	}
}
