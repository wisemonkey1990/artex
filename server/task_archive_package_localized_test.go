package server

import (
	"os"
	"path/filepath"
	"testing"
)

// 说明。
// 说明。
// 说明。
//
// 说明。
// 说明。
// 说明。
// 说明。
// 说明。
//
// 说明。
// 说明。
func TestTaskArchiveStagedAndOriginalCoexistErrorLocalized(t *testing.T) {
	dataDir := t.TempDir()
	archivePath := taskArchivePath(dataDir, 21, "42")
	if err := os.MkdirAll(filepath.Dir(archivePath), archiveDirMode); err != nil {
		t.Fatal(err)
	}
	// 说明。
	if err := os.WriteFile(archivePath, []byte("original"), archiveFileMode); err != nil {
		t.Fatal(err)
	}
	staged := archivePath + ".deleting-21"
	if err := os.WriteFile(staged, []byte("archive"), archiveFileMode); err != nil {
		t.Fatal(err)
	}
	_, moved, err := stageTaskArchivePackageDelete(archivePath, 21)
	if err == nil {
		t.Fatalf("测试文本 测试文本 测试文本 测试文本 测试文本 测试文本 测试文本 测试文本 (moved=%v)", moved)
	}
	if moved {
		t.Fatal("测试文本 测试文本 测试文本 moved 测试文本 true 测试文本 测试文本 测试文本 测试文本")
	}
	assertChineseMessage(t, "staged_and_original_coexist", err.Error())
}
