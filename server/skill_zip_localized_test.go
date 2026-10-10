package server

import (
	"fmt"
	"strings"
	"testing"
	"unicode"
)

// TestSkillZipErrorsLocalized guards the F3b skill_zip.go bundle: every user-facing
// skill-upload error response (fsUploadSkill surfaces these via
// writeErr(400, err.Error()), see server_mgmt.go) must be Chinese, and no zip-method
// display name may carry leftover Korean. Reverting any of these literals back to
// Korean fails here.
func TestSkillZipErrorsLocalized(t *testing.T) {
	// Message constants rendered with sample arguments.
	assertChineseMessage(t, "errSkillZipParse",
		fmt.Errorf(errSkillZipParse, fmt.Errorf("boom")).Error())
	assertChineseMessage(t, "errSkillZipEncrypted",
		fmt.Sprintf(errSkillZipEncrypted, "demo/SKILL.md"))
	assertChineseMessage(t, "errSkillZipUnsupported",
		fmt.Sprintf(errSkillZipUnsupported, "LZMA", 14, "demo/SKILL.md"))

	// A non-zip upload must reach the user as a Chinese hint, not a raw stdlib error.
	if _, err := newSkillZipReader([]byte("this is not a zip archive")); err == nil {
		t.Fatal("비-zip 입력은 오류를 반환해야 합니다")
	} else {
		assertChineseMessage(t, "newSkillZipReader", err.Error())
		if !strings.Contains(err.Error(), "压缩文件") {
			t.Fatalf("parse error = %q, want '压缩文件' 안내", err.Error())
		}
	}

	// No zip-method display name may carry a stray Hangul syllable; the two
	// translated ones (AES, unknown fallback) must be Chinese.
	for m, name := range zipMethodNames {
		for _, r := range name {
			if unicode.Is(unicode.Hangul, r) {
				t.Fatalf("zipMethodNames[%d] = %q 에 한글이 남아 있습니다", m, name)
			}
		}
	}
	assertChineseMessage(t, "zipMethodName(AES)", zipMethodName(zipMethodAES))
	assertChineseMessage(t, "zipMethodName(unknown)", zipMethodName(0xffff))
}
