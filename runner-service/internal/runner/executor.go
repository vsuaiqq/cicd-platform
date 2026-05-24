package runner

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/vsuaiqq/cicd/shared/events"
	"github.com/vsuaiqq/cicd/shared/logger"
)

const maxStepOutput = 65_536

type CancelRegistry struct {
	jobs sync.Map
}

func (r *CancelRegistry) Register(jobID string, cancel context.CancelFunc) {
	r.jobs.Store(jobID, cancel)
}

func (r *CancelRegistry) Unregister(jobID string) {
	r.jobs.Delete(jobID)
}

func (r *CancelRegistry) Cancel(jobID string) bool {
	if v, ok := r.jobs.Load(jobID); ok {
		v.(context.CancelFunc)()
		return true
	}
	return false
}

type ExecConfig struct {
	DockerSocket string
	ArtifactDir  string
	CacheDir     string
}

func Execute(ctx context.Context, task *events.PipelineJobTask, baseWorkDir string, cfg *ExecConfig) *events.PipelineJobResult {
	startedAt := time.Now().Unix()

	result := &events.PipelineJobResult{
		RunID:         task.RunID,
		JobID:         task.JobID,
		JobName:       task.JobName,
		ProjectID:     task.ProjectID,
		AttemptNumber: task.AttemptNumber,
		Status:        "success",
		Steps:         make([]events.StepResult, 0, len(task.Steps)),
		StartedAt:     startedAt,
	}

	renumberSteps := func() {
		for i := range result.Steps {
			result.Steps[i].Index = i
		}
	}
	fail := func() *events.PipelineJobResult {
		renumberSteps()
		result.Status = "failed"
		result.FinishedAt = time.Now().Unix()
		return result
	}
	cancelled := func() *events.PipelineJobResult {
		renumberSteps()
		result.Status = "cancelled"
		result.FinishedAt = time.Now().Unix()
		return result
	}

	if task.TimeoutSeconds > 0 {
		var jobCancel context.CancelFunc
		ctx, jobCancel = context.WithTimeout(ctx, time.Duration(task.TimeoutSeconds)*time.Second)
		defer jobCancel()
	}

	if ctx.Err() != nil {
		return cancelled()
	}

	workDir, err := os.MkdirTemp(baseWorkDir, "run-"+task.RunID+"-job-*")
	if err != nil {
		logger.L().Error().Err(err).Msg("mktemp failed")
		return fail()
	}
	defer os.RemoveAll(workDir)

	keyFile := filepath.Join(workDir, "id_rsa")
	if err := os.WriteFile(keyFile, []byte(task.SSHKey), 0600); err != nil {
		logger.L().Error().Err(err).Msg("write ssh key failed")
		return fail()
	}
	sshEnv := fmt.Sprintf(
		"GIT_SSH_COMMAND=ssh -i %s -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null",
		keyFile,
	)

	runIndex := 0
	recordSetupStep := func(name, status, output string, exitCode int, stepStarted int64) {
		result.Steps = append(result.Steps, events.StepResult{
			Index:      runIndex,
			Name:       name,
			Status:     status,
			Output:     output,
			ExitCode:   exitCode,
			StartedAt:  stepStarted,
			FinishedAt: time.Now().Unix(),
		})
		runIndex++
	}

	repoDir := filepath.Join(workDir, "repo")
	cloneStarted := time.Now().Unix()
	cloneCmd := exec.CommandContext(ctx, "git", "clone", "--depth=1", task.RepoURL, repoDir)
	cloneCmd.Env = append(os.Environ(), sshEnv)
	if out, cloneErr := cloneCmd.CombinedOutput(); cloneErr != nil {
		if ctx.Err() != nil {
			return cancelled()
		}
		logger.L().Error().Err(cloneErr).Str("output", string(out)).Msg("git clone failed")
		recordSetupStep("Clone repository", "failed", string(out), 1, cloneStarted)
		return fail()
	}
	recordSetupStep("Clone repository", "success", "", 0, cloneStarted)

	if strings.TrimSpace(task.CommitSHA) != "" {
		checkoutStarted := time.Now().Unix()
		if err := gitCheckout(ctx, repoDir, task.CommitSHA, sshEnv); err != nil {
			if ctx.Err() != nil {
				return cancelled()
			}
			logger.L().Error().Err(err).Str("sha", task.CommitSHA).Msg("git checkout failed")
			recordSetupStep("Checkout commit", "failed", err.Error(), 1, checkoutStarted)
			return fail()
		}
		recordSetupStep("Checkout commit", "success", "", 0, checkoutStarted)
	}

	image := strings.TrimSpace(task.Image)

	if cfg != nil && cfg.ArtifactDir != "" {
		for _, src := range task.ArtifactsDownload {
			artPath := ArtifactPath(cfg.ArtifactDir, task.RunID, src.JobID)
			if _, err := os.Stat(artPath); err != nil {
				errMsg := fmt.Sprintf("artifact from job %q not found at %s", src.JobName, artPath)
				logger.L().Error().Str("job", src.JobName).Str("path", artPath).Msg("artifact not found")
				result.Steps = append(result.Steps, events.StepResult{
					Name:     "download artifact from " + src.JobName,
					Status:   "failed",
					ExitCode: -1,
					Output:   errMsg,
				})
				return fail()
			}
			var extractErr error
			if image != "" {
				extractErr = ExtractArtifactsDocker(ctx, repoDir, artPath, cfg.DockerSocket)
			} else {
				extractErr = ExtractArtifactsHost(repoDir, artPath)
			}
			if extractErr != nil {
				errMsg := fmt.Sprintf("artifact extract from job %q failed: %v", src.JobName, extractErr)
				logger.L().Error().Err(extractErr).Str("job", src.JobName).Msg("artifact extract failed")
				result.Steps = append(result.Steps, events.StepResult{
					Name:     "download artifact from " + src.JobName,
					Status:   "failed",
					ExitCode: -1,
					Output:   errMsg,
				})
				return fail()
			}
			logger.L().Info().Str("job", src.JobName).Msg("artifacts extracted")
		}
	}

	var cacheMounts []string
	if cfg != nil && cfg.CacheDir != "" && image != "" {
		for _, entry := range task.Cache {
			key := EvalCacheKey(entry.Key, repoDir)
			for _, p := range entry.Paths {
				hostDir := CacheHostDir(cfg.CacheDir, task.ProjectID, key, p)
				if err := os.MkdirAll(hostDir, 0755); err == nil {
					cacheMounts = append(cacheMounts, "-v", hostDir+":"+p)
					logger.L().Debug().Str("host", hostDir).Str("container", p).Str("key", key).Msg("cache mount")
				}
			}
		}
	}

	for _, step := range task.Steps {
		var sr events.StepResult
		if image != "" {
			sr = executeStepWithRetry(ctx, step, func(stepCtx context.Context) events.StepResult {
				return executeStepDocker(stepCtx, ctx, repoDir, image, step, task.Env, task.Secrets, cfg, cacheMounts)
			})
		} else {
			sr = executeStepWithRetry(ctx, step, func(stepCtx context.Context) events.StepResult {
				return executeStepHost(stepCtx, ctx, repoDir, step, task.Env, task.Secrets)
			})
		}
		sr.Index = step.Index
		result.Steps = append(result.Steps, sr)
		if sr.Status == "cancelled" {
			return cancelled()
		}
		if sr.Status == "failed" && !step.ContinueOnError {
			return fail()
		}
	}

	if cfg != nil && cfg.ArtifactDir != "" && len(task.ArtifactsUpload) > 0 {
		artPath := ArtifactPath(cfg.ArtifactDir, task.RunID, task.JobID)
		var saveErr error
		if image != "" {
			saveErr = SaveArtifactsDocker(ctx, repoDir, artPath, cfg.DockerSocket, task.ArtifactsUpload)
		} else {
			saveErr = SaveArtifactsHost(repoDir, artPath, task.ArtifactsUpload)
		}
		if saveErr != nil {
			logger.L().Warn().Err(saveErr).Msg("artifact save failed (non-fatal)")
		} else {
			logger.L().Info().Str("path", artPath).Msg("artifacts saved")
		}
	}

	renumberSteps()
	result.FinishedAt = time.Now().Unix()
	return result
}

func gitCheckout(ctx context.Context, repoDir, commitSHA, sshEnv string) error {
	checkout := func() error {
		cmd := exec.CommandContext(ctx, "git", "-C", repoDir, "checkout", commitSHA)
		cmd.Env = os.Environ()
		if out, err := cmd.CombinedOutput(); err != nil {
			return fmt.Errorf("%w: %s", err, out)
		}
		return nil
	}

	if err := checkout(); err == nil {
		return nil
	}
	if ctx.Err() != nil {
		return ctx.Err()
	}

	fetchCmd := exec.CommandContext(ctx, "git", "-C", repoDir, "fetch", "--depth=1", "origin", commitSHA)
	fetchCmd.Env = append(os.Environ(), sshEnv)
	if out, err := fetchCmd.CombinedOutput(); err != nil {
		if ctx.Err() != nil {
			return ctx.Err()
		}

		logger.L().Debug().Err(err).Str("output", string(out)).Msg("git fetch for specific SHA failed")
	}

	return checkout()
}

func executeStepWithRetry(
	ctx context.Context,
	step events.PipelineStep,
	fn func(stepCtx context.Context) events.StepResult,
) events.StepResult {
	maxAttempts := step.RetryMax + 1
	var last events.StepResult

	for attempt := 0; attempt < maxAttempts; attempt++ {
		if attempt > 0 {
			logger.L().Info().
				Int("step", step.Index).Str("name", step.Name).
				Int("attempt", attempt+1).Int("max", maxAttempts).
				Msg("retrying step")
		}

		stepCtx := ctx
		var stepCancel context.CancelFunc
		if step.TimeoutSeconds > 0 {
			stepCtx, stepCancel = context.WithTimeout(ctx, time.Duration(step.TimeoutSeconds)*time.Second)
		}

		last = fn(stepCtx)

		if stepCancel != nil {
			stepCancel()
		}
		if last.Status != "failed" {
			return last
		}
		if attempt < maxAttempts-1 {
			logger.L().Info().
				Int("step", step.Index).Str("name", step.Name).
				Int("attempt", attempt+1).
				Msg("step failed, retrying")
		}
	}
	return last
}

func executeStepHost(
	stepCtx, jobCtx context.Context,
	workDir string,
	step events.PipelineStep,
	env, secrets map[string]string,
) events.StepResult {
	stepStart := time.Now().Unix()
	logger.L().Info().Int("step", step.Index).Str("name", step.Name).Msg("running step (host)")

	cmdEnv := buildEnvList(buildEnvList(os.Environ(), env), secrets)

	var buf bytes.Buffer
	cmd := exec.CommandContext(stepCtx, "sh", "-c", step.Run)
	cmd.Dir = workDir
	cmd.Env = cmdEnv
	cmd.Stdout = &buf
	cmd.Stderr = &buf

	runErr := cmd.Run()
	return makeStepResult(step, maskSecrets(buf.String(), secrets), runErr, stepStart, jobCtx)
}

func executeStepDocker(
	stepCtx, jobCtx context.Context,
	workDir, image string,
	step events.PipelineStep,
	env, secrets map[string]string,
	cfg *ExecConfig,
	cacheMounts []string,
) events.StepResult {
	stepStart := time.Now().Unix()
	logger.L().Info().Int("step", step.Index).Str("image", image).Str("name", step.Name).Msg("running step (docker)")

	args := []string{
		"run", "--rm", "--init",
		"-v", workDir + ":/workspace",
		"-w", "/workspace",
	}
	args = append(args, cacheMounts...)
	for k, v := range env {
		args = append(args, "-e", k+"="+v)
	}

	for k, v := range secrets {
		args = append(args, "-e", k+"="+v)
	}
	args = append(args, image, "sh", "-c", step.Run)

	cmd := exec.CommandContext(stepCtx, "docker", args...)
	if cfg != nil && cfg.DockerSocket != "" {
		cmd.Env = append(os.Environ(), "DOCKER_HOST=unix://"+cfg.DockerSocket)
	}

	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf

	runErr := cmd.Run()
	return makeStepResult(step, maskSecrets(buf.String(), secrets), runErr, stepStart, jobCtx)
}

func makeStepResult(step events.PipelineStep, output string, runErr error, startedAt int64, jobCtx context.Context) events.StepResult {
	status := "success"
	exitCode := 0

	if runErr != nil {
		if jobCtx.Err() != nil {
			status = "cancelled"
			exitCode = -1
		} else {
			status = "failed"
			var exitErr *exec.ExitError
			if errors.As(runErr, &exitErr) {
				exitCode = exitErr.ExitCode()
			} else {
				exitCode = 1
			}
		}
	}

	if len(output) > maxStepOutput {
		output = output[:maxStepOutput] + "\n[output truncated]"
	}

	return events.StepResult{
		Index:      step.Index,
		Name:       step.Name,
		Status:     status,
		Output:     output,
		ExitCode:   exitCode,
		StartedAt:  startedAt,
		FinishedAt: time.Now().Unix(),
	}
}

func maskSecrets(output string, secrets map[string]string) string {
	for _, v := range secrets {
		if len(v) < 3 {
			continue
		}
		output = strings.ReplaceAll(output, v, "***")
	}
	return output
}

func buildEnvList(base []string, extra map[string]string) []string {
	out := make([]string, len(base), len(base)+len(extra))
	copy(out, base)
	for k, v := range extra {
		out = append(out, k+"="+v)
	}
	return out
}
