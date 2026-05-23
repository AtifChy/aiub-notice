package notice

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

// setupCacheTestEnv redirects the cache directory to a temporary folder
// to ensure tests run in isolation and don't overwrite user data.
func setupCacheTestEnv(t *testing.T) string {
	t.Helper()

	tmpDir := t.TempDir()

	// Redirect app data/cache directories for both Windows and Unix
	t.Setenv("LOCALAPPDATA", tmpDir)
	t.Setenv("XDG_CACHE_HOME", tmpDir)

	return tmpDir
}

func Test_CachedNotices(t *testing.T) {
	setupCacheTestEnv(t)

	// 1. Verify reading from non-existent cache fails gracefully
	_, err := GetCachedNotices()
	if err == nil {
		t.Fatal("expected error reading non-existent cache, got nil")
	}

	// Mock notice data
	mockNotices := []Notice{
		{
			Date:  time.Date(2026, 5, 23, 0, 0, 0, 0, time.UTC),
			Title: "Test Notice 1",
			Desc:  "This is notice 1",
			Link:  "https://example.com/1",
		},
		{
			Date:  time.Date(2026, 5, 24, 0, 0, 0, 0, time.UTC),
			Title: "Test Notice 2",
			Desc:  "This is notice 2",
			Link:  "https://example.com/2",
		},
	}

	// 2. Write notices to the cache
	err = storeCachedNotices(mockNotices)
	if err != nil {
		t.Fatalf("failed to store cached notices: %v", err)
	}

	// 3. Read back notices from the cache
	cached, err := GetCachedNotices()
	if err != nil {
		t.Fatalf("failed to read cached notices: %v", err)
	}

	if len(cached) != len(mockNotices) {
		t.Fatalf("expected %d notices, got %d", len(mockNotices), len(cached))
	}

	for i, n := range cached {
		if !n.Date.Equal(mockNotices[i].Date) {
			t.Errorf("notice %d: expected Date %v, got %v", i, mockNotices[i].Date, n.Date)
		}
		if n.Title != mockNotices[i].Title {
			t.Errorf("notice %d: expected Title %q, got %q", i, mockNotices[i].Title, n.Title)
		}
		if n.Desc != mockNotices[i].Desc {
			t.Errorf("notice %d: expected Desc %q, got %q", i, mockNotices[i].Desc, n.Desc)
		}
		if n.Link != mockNotices[i].Link {
			t.Errorf("notice %d: expected Link %q, got %q", i, mockNotices[i].Link, n.Link)
		}
	}
}

func Test_SeenNotices(t *testing.T) {
	setupCacheTestEnv(t)

	// 1. Reading seen notices when file does not exist should return empty map, no error
	seen, err := LoadSeenNotices()
	if err != nil {
		t.Fatalf("unexpected error when seen notices file is missing: %v", err)
	}
	if len(seen) != 0 {
		t.Fatalf("expected empty seen notices map, got size %d", len(seen))
	}

	// Mock seen notices
	mockSeen := map[string]struct{}{
		"https://example.com/notice-a": {},
		"https://example.com/notice-b": {},
	}

	// 2. Save seen notices
	err = SaveSeenNotices(mockSeen)
	if err != nil {
		t.Fatalf("failed to save seen notices: %v", err)
	}

	// 3. Load seen notices back
	loaded, err := LoadSeenNotices()
	if err != nil {
		t.Fatalf("failed to load seen notices: %v", err)
	}

	if len(loaded) != len(mockSeen) {
		t.Fatalf("expected %d seen notices, got %d", len(mockSeen), len(loaded))
	}

	for k := range mockSeen {
		if _, ok := loaded[k]; !ok {
			t.Errorf("expected key %q to exist in loaded seen notices", k)
		}
	}

	// 4. Corrupt seen notices file manually
	seenPath, err := GetSeenNoticesPath()
	if err != nil {
		t.Fatalf("failed to get seen notices path: %v", err)
	}

	err = os.WriteFile(seenPath, []byte("{invalid-json}"), 0o644)
	if err != nil {
		t.Fatalf("failed to write corrupted file: %v", err)
	}

	// Loading corrupted seen notices should fall back to empty map, no error
	fallback, err := LoadSeenNotices()
	if err != nil {
		t.Fatalf("expected no error when loading corrupted seen notices, got: %v", err)
	}
	if len(fallback) != 0 {
		t.Fatalf("expected empty map fallback for corrupted file, got size %d", len(fallback))
	}
}

func Test_GetSeenNoticesPath(t *testing.T) {
	setupCacheTestEnv(t)

	path, err := GetSeenNoticesPath()
	if err != nil {
		t.Fatalf("unexpected error getting seen notices path: %v", err)
	}

	expectedBase := "seen_notices.json"
	if filepath.Base(path) != expectedBase {
		t.Errorf("expected base name %q, got %q", expectedBase, filepath.Base(path))
	}
}
