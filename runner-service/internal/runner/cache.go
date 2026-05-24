package runner

import (
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var checksumRe = regexp.MustCompile(`\$\{\{\s*checksum\s+"([^"]+)"\s*\}\}`)

func EvalCacheKey(keyTemplate, repoDir string) string {
	return checksumRe.ReplaceAllStringFunc(keyTemplate, func(match string) string {
		sub := checksumRe.FindStringSubmatch(match)
		if len(sub) < 2 {
			return match
		}
		path := filepath.Join(repoDir, sub[1])
		data, err := os.ReadFile(path)
		if err != nil {
			return "missing"
		}
		sum := sha256.Sum256(data)
		return fmt.Sprintf("%x", sum)[:16]
	})
}

func CacheHostDir(cacheDir, projectID, key, containerPath string) string {
	keyHash := fmt.Sprintf("%x", sha256.Sum256([]byte(key)))[:16]
	sanitized := strings.NewReplacer(
		"/", "_",
		"\\", "_",
		":", "_",
		" ", "_",
	).Replace(strings.TrimPrefix(containerPath, "/"))
	if sanitized == "" {
		sanitized = "root"
	}
	return filepath.Join(cacheDir, projectID, keyHash, sanitized)
}
