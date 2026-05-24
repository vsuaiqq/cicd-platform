package runner

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func ArtifactPath(artifactDir, runID, jobID string) string {
	return filepath.Join(artifactDir, runID, jobID+".tar.gz")
}

func SaveArtifactsDocker(ctx context.Context, repoDir, artifactPath, dockerSocket string, paths []string) error {
	if err := os.MkdirAll(filepath.Dir(artifactPath), 0755); err != nil {
		return fmt.Errorf("artifacts: mkdir: %w", err)
	}
	archiveName := filepath.Base(artifactPath)
	tarArgs := "cd /workspace && tar -czf /artifacts/" + archiveName + " " + strings.Join(paths, " ") + " 2>/dev/null"
	args := []string{
		"run", "--rm",
		"-v", repoDir + ":/workspace:ro",
		"-v", filepath.Dir(artifactPath) + ":/artifacts",
		"alpine", "sh", "-c", tarArgs,
	}
	cmd := exec.CommandContext(ctx, "docker", args...)
	applyDockerSocket(cmd, dockerSocket)
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("artifacts: save docker: %w: %s", err, out)
	}
	return nil
}

func ExtractArtifactsDocker(ctx context.Context, repoDir, artifactPath, dockerSocket string) error {
	if _, err := os.Stat(artifactPath); os.IsNotExist(err) {
		return nil
	}
	tarArgs := "cd /workspace && tar -xzf /artifacts/" + filepath.Base(artifactPath)
	args := []string{
		"run", "--rm",
		"-v", repoDir + ":/workspace",
		"-v", filepath.Dir(artifactPath) + ":/artifacts:ro",
		"alpine", "sh", "-c", tarArgs,
	}
	cmd := exec.CommandContext(ctx, "docker", args...)
	applyDockerSocket(cmd, dockerSocket)
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("artifacts: extract docker: %w: %s", err, out)
	}
	return nil
}

func SaveArtifactsHost(repoDir, artifactPath string, paths []string) error {
	if err := os.MkdirAll(filepath.Dir(artifactPath), 0755); err != nil {
		return fmt.Errorf("artifacts: mkdir: %w", err)
	}
	args := append([]string{"-czf", artifactPath, "-C", repoDir}, paths...)
	if out, err := exec.Command("tar", args...).CombinedOutput(); err != nil {
		return fmt.Errorf("artifacts: save host: %w: %s", err, out)
	}
	return nil
}

func ExtractArtifactsHost(repoDir, artifactPath string) error {
	if _, err := os.Stat(artifactPath); os.IsNotExist(err) {
		return nil
	}
	if out, err := exec.Command("tar", "-xzf", artifactPath, "-C", repoDir).CombinedOutput(); err != nil {
		return fmt.Errorf("artifacts: extract host: %w: %s", err, out)
	}
	return nil
}

func applyDockerSocket(cmd *exec.Cmd, dockerSocket string) {
	if dockerSocket != "" {
		cmd.Env = append(os.Environ(), "DOCKER_HOST=unix://"+dockerSocket)
	}
}
