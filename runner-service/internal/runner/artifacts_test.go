package runner

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNormalizeArtifactPaths(t *testing.T) {
	got := normalizeArtifactPaths([]string{"./bin/", " ./dist ", "/out"})
	want := []string{"bin", "dist", "out"}
	if len(got) != len(want) {
		t.Fatalf("got %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("got %v, want %v", got, want)
		}
	}
}

func TestSaveAndExtractArtifactsHost(t *testing.T) {
	dir := t.TempDir()
	repo := filepath.Join(dir, "repo")
	if err := os.MkdirAll(filepath.Join(repo, "bin"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(repo, "bin", "app"), []byte("binary"), 0755); err != nil {
		t.Fatal(err)
	}

	archive := filepath.Join(dir, "artifacts", "run", "job.tar.gz")
	if err := SaveArtifactsHost(repo, archive, []string{"bin"}); err != nil {
		t.Fatalf("save: %v", err)
	}
	if _, err := os.Stat(archive); err != nil {
		t.Fatalf("archive missing: %v", err)
	}

	extractRepo := filepath.Join(dir, "extract")
	if err := os.MkdirAll(extractRepo, 0755); err != nil {
		t.Fatal(err)
	}
	if err := ExtractArtifactsHost(extractRepo, archive); err != nil {
		t.Fatalf("extract: %v", err)
	}
	if _, err := os.Stat(filepath.Join(extractRepo, "bin", "app")); err != nil {
		t.Fatalf("extracted binary missing: %v", err)
	}
}
