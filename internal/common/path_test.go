package common_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/AtifChy/aiub-notice/internal/common"
)

func withTempLogFile(t *testing.T, testFunc func(logPath string)) {
	t.Helper()

	tmpDir := t.TempDir()
	logPath := filepath.Join(tmpDir, "test.log")

	ogGetLogPath := common.GetLogPath
	common.GetLogPath = func() string { return logPath }
	defer func() { common.GetLogPath = ogGetLogPath }()

	testFunc(logPath)
}

func TestGetLogFile(t *testing.T) {
	tests := []struct {
		name  string
		setup func(path string)
		check func(t *testing.T, path string)
	}{
		{
			name: "create new file if not exist",
			setup: func(path string) {
				// file does not exist, nothing to do
			},
			check: func(t *testing.T, path string) {
				if _, err := os.Stat(path); os.IsNotExist(err) {
					t.Errorf("expected log file to be created, but it does not exist")
				}
			},
		},
		{
			name: "keep small file unchanged",
			setup: func(path string) {
				os.WriteFile(path, []byte("test"), 0664)
			},
			check: func(t *testing.T, path string) {
				data, _ := os.ReadFile(path)
				if string(data) != "test" {
					t.Errorf("expected log file content to be unchanged, got %q", string(data))
				}
			},
		},
		{
			name: "truncate large file",
			setup: func(path string) {
				large := make([]byte, 6*1024*1024) // 6 MB
				os.WriteFile(path, large, 0664)
			},
			check: func(t *testing.T, path string) {
				info, _ := os.Stat(path)
				if info.Size() != 0 {
					t.Errorf("expected log file to be truncated, but size is %d", info.Size())
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			withTempLogFile(t, func(logPath string) {
				tt.setup(logPath)

				file, err := common.GetLogFile()
				if err != nil {
					t.Fatalf("GetLogFile() error: %v", err)
				}
				defer file.Close()

				tt.check(t, logPath)
			})
		})
	}
}
