package orchestrator

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/vsuaiqq/cicd/orchestrator-service/internal/db"
	"github.com/vsuaiqq/cicd/orchestrator-service/internal/pipeline"
	orchWs "github.com/vsuaiqq/cicd/orchestrator-service/internal/ws"
	sharedEvents "github.com/vsuaiqq/cicd/shared/events"
	"github.com/vsuaiqq/cicd/shared/logger"
	pb "github.com/vsuaiqq/cicd/shared/proto/gen/projects"
)

type JobPublisher interface {
	PublishJob(ctx context.Context, task *sharedEvents.PipelineJobTask) error
}

type CancelPublisher interface {
	PublishCancel(ctx context.Context, runID, jobID string) error
}

type RunEventPublisher interface {
	PublishRunCompleted(ctx context.Context, evt *sharedEvents.RunCompletedEvent) error
}

type Orchestrator struct {
	repo           *db.Repository
	jobPub         JobPublisher
	cancelPub      CancelPublisher
	runEventPub    RunEventPublisher
	projectsClient pb.ProjectsServiceClient
	hub            *orchWs.Hub
}

func New(
	repo *db.Repository,
	jobPub JobPublisher,
	cancelPub CancelPublisher,
	runEventPub RunEventPublisher,
	projectsClient pb.ProjectsServiceClient,
	hub *orchWs.Hub,
) *Orchestrator {
	return &Orchestrator{
		repo:           repo,
		jobPub:         jobPub,
		cancelPub:      cancelPub,
		runEventPub:    runEventPub,
		projectsClient: projectsClient,
		hub:            hub,
	}
}

func (o *Orchestrator) HandleGitEvent(ctx context.Context, event *sharedEvents.GitEvent) error {
	if event.ProjectID == "" {
		logger.L().Warn().Msg("git event has no project_id, skipping")
		return nil
	}
	logger.L().Info().
		Str("project", event.ProjectID).Str("repo", event.Repository.URL).
		Str("branch", event.Branch).Str("sha", event.CommitSHA).
		Msg("git event")

	info, err := o.getProjectInfo(ctx, event.ProjectID)
	if err != nil {
		return fmt.Errorf("orchestrator: get project info for %s: %w", event.ProjectID, err)
	}
	repoURL := info.repoURL
	if repoURL == "" {
		repoURL = event.Repository.URL
	}

	var pipelineData []byte
	if info.pipelineYAMLOverride != "" {
		logger.L().Info().Str("project", event.ProjectID).Msg("using pipeline YAML override")
		pipelineData = []byte(info.pipelineYAMLOverride)
	} else {
		pipelineData, err = cloneAndReadPipeline(ctx, repoURL, event.CommitSHA, info.sshKey)
		if err != nil {
			logger.L().Warn().Err(err).Str("project", event.ProjectID).Msg("failed to load pipeline")
			return nil
		}
	}

	pl, err := pipeline.LoadBytes(pipelineData)
	if err != nil {
		logger.L().Warn().Err(err).Str("project", event.ProjectID).Msg("invalid pipeline yaml")
		return nil
	}

	if !pl.BranchAllowed(event.Branch) {
		logger.L().Info().Str("branch", event.Branch).Msg("branch not in trigger filter, skipping")
		return nil
	}

	run, err := o.repo.CreateRun(ctx, event.ProjectID, event.CommitSHA, event.Branch, repoURL)
	if err != nil {
		return fmt.Errorf("orchestrator: create run: %w", err)
	}
	logger.L().Info().Str("run_id", run.ID).Str("project", event.ProjectID).Msg("created run")

	if err := o.repo.SetPipelineYAML(ctx, run.ID, string(pipelineData)); err != nil {
		logger.L().Warn().Err(err).Str("run_id", run.ID).Msg("failed to store pipeline yaml")
	}

	jobMap := make(map[string]*db.PipelineJob, len(pl.Jobs))
	for name, plJob := range pl.Jobs {
		displayName := plJob.Name
		if displayName == "" {
			displayName = name
		}
		j, err := o.repo.CreateJob(ctx, run.ID, name, displayName)
		if err != nil {
			return fmt.Errorf("orchestrator: create job %q: %w", name, err)
		}
		jobMap[name] = j
	}

	if err := o.repo.UpdateRunStatus(ctx, run.ID, db.StatusRunning); err != nil {
		return fmt.Errorf("orchestrator: update run status: %w", err)
	}

	o.hub.MarkActive(run.ID)
	o.hub.Broadcast(run.ID, orchWs.Event{
		Type:      orchWs.EventRunUpdated,
		RunID:     run.ID,
		ProjectID: event.ProjectID,
		Status:    db.StatusRunning,
	})
	return o.dispatchReadyJobs(ctx, run.ID, event.ProjectID, pl, jobMap, map[string]bool{}, info.sshKey, repoURL, event.CommitSHA, info.env, info.secrets)
}

func (o *Orchestrator) HandleJobResult(ctx context.Context, result *sharedEvents.PipelineJobResult) error {
	logger.L().Info().
		Str("run_id", result.RunID).Str("job_id", result.JobID).Str("status", result.Status).
		Msg("job result")

	if err := o.repo.UpdateJobStatus(ctx, result.JobID, result.Status); err != nil {
		return fmt.Errorf("orchestrator: update job status: %w", err)
	}
	if result.AttemptNumber > 1 {
		_ = o.repo.UpdateJobAttempt(ctx, result.JobID, result.AttemptNumber)
	}

	for _, sr := range result.Steps {
		startedAt := time.Unix(sr.StartedAt, 0)
		finishedAt := time.Unix(sr.FinishedAt, 0)
		if err := o.repo.UpsertStep(ctx, &db.PipelineStep{
			JobID:      result.JobID,
			StepIndex:  sr.Index,
			Name:       sr.Name,
			Status:     sr.Status,
			LogOutput:  sr.Output,
			ExitCode:   sr.ExitCode,
			StartedAt:  &startedAt,
			FinishedAt: &finishedAt,
		}); err != nil {
			logger.L().Error().Err(err).Msg("upsert step error")
		}
	}

	jobs, err := o.repo.GetJobsByRun(ctx, result.RunID)
	if err != nil {
		return fmt.Errorf("orchestrator: get jobs: %w", err)
	}

	jobName := jobNameByID(jobs, result.JobID)
	steps := make([]orchWs.StepEvent, len(result.Steps))
	for i, sr := range result.Steps {
		steps[i] = orchWs.StepEvent{
			Index:      sr.Index,
			Name:       sr.Name,
			Status:     sr.Status,
			ExitCode:   sr.ExitCode,
			LogOutput:  sr.Output,
			StartedAt:  sr.StartedAt,
			FinishedAt: sr.FinishedAt,
		}
	}
	now := time.Now()
	o.hub.Broadcast(result.RunID, orchWs.Event{
		Type:            orchWs.EventJobUpdated,
		RunID:           result.RunID,
		Status:          result.Status,
		JobName:         jobName,
		JobID:           result.JobID,
		Steps:           steps,
		JobFinishedAtMs: now.UnixMilli(),
	})

	completed, allDone, anyFailed := computeJobState(jobs)
	if allDone {
		finalStatus := db.StatusSuccess
		if anyFailed {
			finalStatus = db.StatusFailed
		}
		_ = o.repo.UpdateRunStatus(ctx, result.RunID, finalStatus)

		finishedRun, _ := o.repo.GetRun(ctx, result.RunID)
		projectID := ""
		if finishedRun != nil {
			projectID = finishedRun.ProjectID
		}
		o.hub.Broadcast(result.RunID, orchWs.Event{
			Type:            orchWs.EventRunUpdated,
			RunID:           result.RunID,
			ProjectID:       projectID,
			Status:          finalStatus,
			RunFinishedAtMs: time.Now().UnixMilli(),
		})
		logger.L().Info().Str("run_id", result.RunID).Str("status", finalStatus).Msg("run finished")
		if finishedRun != nil {
			o.publishRunCompleted(ctx, finishedRun, finalStatus)
		}
		return nil
	}

	run, err := o.repo.GetRun(ctx, result.RunID)
	if err != nil || run == nil {
		logger.L().Warn().Err(err).Str("run_id", result.RunID).Msg("run not found after job result")
		return nil
	}
	pl, err := pipeline.LoadBytes([]byte(run.PipelineYAML))
	if err != nil {
		logger.L().Warn().Err(err).Str("run_id", run.ID).Msg("invalid stored pipeline yaml")
		return nil
	}
	jobMap := jobMapByName(jobs)

	if result.Status == db.StatusFailed || result.Status == db.StatusCancelled {

		if plJob, ok := pl.Jobs[jobName]; ok && result.Status == db.StatusFailed &&
			plJob.Retry.Max > 0 && result.AttemptNumber < plJob.Retry.Max+1 {

			nextAttempt := result.AttemptNumber + 1
			logger.L().Info().
				Str("job", jobName).Int("attempt", nextAttempt).Int("max", plJob.Retry.Max+1).Str("run_id", result.RunID).
				Msg("retrying job")

			dbJob, err := o.repo.GetJobByID(ctx, result.JobID)
			if err == nil && dbJob != nil {
				_ = o.repo.UpdateJobAttempt(ctx, dbJob.ID, nextAttempt)

				_ = o.repo.DeleteStepsByJob(ctx, dbJob.ID)
				info, _ := o.getProjectInfo(ctx, run.ProjectID)

				delay := retryDelay(nextAttempt)
				logger.L().Info().Str("delay", delay.String()).Str("job", jobName).Msg("retry back-off")
				select {
				case <-ctx.Done():
					return ctx.Err()
				case <-time.After(delay):
				}

				if dispatchErr := o.dispatchJob(ctx, result.RunID, run.ProjectID, dbJob, plJob, pl,
					info.sshKey, run.RepoURL, run.CommitSHA, info.env, info.secrets,
					jobMap, nextAttempt); dispatchErr != nil {
					logger.L().Error().Err(dispatchErr).Str("job", jobName).Msg("retry dispatch failed")
				} else {
					return nil
				}
			}
		}

		o.skipDependents(ctx, result.RunID, jobName, pl, jobMap)

		updatedJobs, err := o.repo.GetJobsByRun(ctx, result.RunID)
		if err == nil {
			if _, allDone, _ := computeJobState(updatedJobs); allDone {
				_ = o.repo.UpdateRunStatus(ctx, result.RunID, db.StatusFailed)
				projectID := ""
				if run != nil {
					projectID = run.ProjectID
				}
				o.hub.Broadcast(result.RunID, orchWs.Event{
					Type:            orchWs.EventRunUpdated,
					RunID:           result.RunID,
					ProjectID:       projectID,
					Status:          db.StatusFailed,
					RunFinishedAtMs: time.Now().UnixMilli(),
				})
				logger.L().Info().Str("run_id", result.RunID).Msg("run finished: failed (all dependents skipped)")
				if run != nil {
					o.publishRunCompleted(ctx, run, db.StatusFailed)
				}
			}
		}
		return nil
	}

	info, err := o.getProjectInfo(ctx, run.ProjectID)
	if err != nil {
		logger.L().Warn().Err(err).Str("project", run.ProjectID).Msg("could not get project info")
		return nil
	}
	return o.dispatchReadyJobs(ctx, result.RunID, run.ProjectID, pl, jobMap, completed, info.sshKey, run.RepoURL, run.CommitSHA, info.env, info.secrets)
}

func (o *Orchestrator) dispatchReadyJobs(
	ctx context.Context,
	runID string,
	projectID string,
	pl *pipeline.Pipeline,
	jobMap map[string]*db.PipelineJob,
	completed map[string]bool,
	sshKey, repoURL, commitSHA string,
	projectEnv map[string]string,
	secrets map[string]string,
) error {
	for name, plJob := range pl.Jobs {
		dbJob, ok := jobMap[name]
		if !ok || dbJob.Status != db.StatusPending {
			continue
		}
		if !allCompleted(plJob.Needs, completed) {
			continue
		}

		if bool(plJob.Approval) {
			logger.L().Info().Str("job", name).Str("run_id", runID).Msg("job requires approval")
			_ = o.repo.UpdateJobStatus(ctx, dbJob.ID, db.StatusAwaitingApproval)
			o.hub.Broadcast(runID, orchWs.Event{
				Type:    orchWs.EventJobAwaitingApproval,
				RunID:   runID,
				Status:  db.StatusAwaitingApproval,
				JobName: dbJob.Name,
				JobID:   dbJob.ID,
			})
			continue
		}
		if err := o.dispatchJob(ctx, runID, projectID, dbJob, plJob, pl, sshKey, repoURL, commitSHA, projectEnv, secrets, jobMap, 1); err != nil {
			logger.L().Error().Err(err).Str("job", name).Msg("failed to dispatch job")
		}
	}
	return nil
}

func (o *Orchestrator) CancelRun(ctx context.Context, runID string) error {
	runningJobIDs, err := o.repo.CancelRunNonTerminalJobs(ctx, runID)
	if err != nil {
		return fmt.Errorf("orchestrator: cancel run %s: %w", runID, err)
	}
	if err := o.repo.UpdateRunStatus(ctx, runID, db.StatusCancelled); err != nil {
		return fmt.Errorf("orchestrator: cancel run status %s: %w", runID, err)
	}

	for _, jobID := range runningJobIDs {
		if pubErr := o.cancelPub.PublishCancel(ctx, runID, jobID); pubErr != nil {
			logger.L().Error().Err(pubErr).Str("job_id", jobID).Msg("failed to publish cancel")
		}
	}

	cancelledRun, _ := o.repo.GetRun(ctx, runID)
	projectID := ""
	if cancelledRun != nil {
		projectID = cancelledRun.ProjectID
	}
	o.hub.Broadcast(runID, orchWs.Event{
		Type:            orchWs.EventRunUpdated,
		RunID:           runID,
		ProjectID:       projectID,
		Status:          db.StatusCancelled,
		RunFinishedAtMs: time.Now().UnixMilli(),
	})
	logger.L().Info().Str("run_id", runID).Int("signals", len(runningJobIDs)).Msg("run cancelled")
	if cancelledRun != nil {
		o.publishRunCompleted(ctx, cancelledRun, db.StatusCancelled)
	}
	return nil
}

func (o *Orchestrator) ApproveJob(ctx context.Context, runID, jobID string) error {
	dbJob, err := o.repo.GetJobByID(ctx, jobID)
	if err != nil || dbJob == nil {
		return fmt.Errorf("orchestrator: job %s not found", jobID)
	}
	if dbJob.RunID != runID {
		return fmt.Errorf("orchestrator: job %s does not belong to run %s", jobID, runID)
	}
	if dbJob.Status != db.StatusAwaitingApproval {
		return fmt.Errorf("orchestrator: job %s is not awaiting approval (status=%s)", jobID, dbJob.Status)
	}

	run, err := o.repo.GetRun(ctx, runID)
	if err != nil || run == nil {
		return fmt.Errorf("orchestrator: run %s not found", runID)
	}
	pl, err := pipeline.LoadBytes([]byte(run.PipelineYAML))
	if err != nil {
		return fmt.Errorf("orchestrator: load pipeline yaml for run %s: %w", runID, err)
	}
	plJob, ok := pl.Jobs[dbJob.Name]
	if !ok {
		return fmt.Errorf("orchestrator: job %q not found in pipeline", dbJob.Name)
	}

	info, err := o.getProjectInfo(ctx, run.ProjectID)
	if err != nil {
		return fmt.Errorf("orchestrator: get project info: %w", err)
	}

	logger.L().Info().Str("job", dbJob.Name).Str("run_id", runID).Msg("job approved")

	o.hub.MarkActive(runID)

	_ = o.repo.UpdateJobStatus(ctx, dbJob.ID, db.StatusPending)
	jobs, _ := o.repo.GetJobsByRun(ctx, runID)
	jm := jobMapByName(jobs)
	return o.dispatchJob(ctx, runID, run.ProjectID, dbJob, plJob, pl, info.sshKey, run.RepoURL, run.CommitSHA, info.env, info.secrets, jm, dbJob.Attempt)
}

func (o *Orchestrator) RejectJob(ctx context.Context, runID, jobID string) error {
	dbJob, err := o.repo.GetJobByID(ctx, jobID)
	if err != nil || dbJob == nil {
		return fmt.Errorf("orchestrator: job %s not found", jobID)
	}
	if dbJob.RunID != runID {
		return fmt.Errorf("orchestrator: job %s does not belong to run %s", jobID, runID)
	}
	if dbJob.Status != db.StatusAwaitingApproval {
		return fmt.Errorf("orchestrator: job %s is not awaiting approval (status=%s)", jobID, dbJob.Status)
	}

	_ = o.repo.UpdateJobStatus(ctx, dbJob.ID, db.StatusCancelled)
	o.hub.Broadcast(runID, orchWs.Event{
		Type:    orchWs.EventJobUpdated,
		RunID:   runID,
		Status:  db.StatusCancelled,
		JobName: dbJob.Name,
		JobID:   dbJob.ID,
	})

	run, err := o.repo.GetRun(ctx, runID)
	if err == nil && run != nil {
		if pl, err := pipeline.LoadBytes([]byte(run.PipelineYAML)); err == nil {
			jobs, _ := o.repo.GetJobsByRun(ctx, runID)
			jobMap := jobMapByName(jobs)
			o.skipDependents(ctx, runID, dbJob.Name, pl, jobMap)

			updatedJobs, _ := o.repo.GetJobsByRun(ctx, runID)
			if _, allDone, _ := computeJobState(updatedJobs); allDone {
				_ = o.repo.UpdateRunStatus(ctx, runID, db.StatusFailed)
				rejProjectID := ""
				if run != nil {
					rejProjectID = run.ProjectID
				}
				o.hub.Broadcast(runID, orchWs.Event{
					Type:            orchWs.EventRunUpdated,
					RunID:           runID,
					ProjectID:       rejProjectID,
					Status:          db.StatusFailed,
					RunFinishedAtMs: time.Now().UnixMilli(),
				})
				if run != nil {
					o.publishRunCompleted(ctx, run, db.StatusFailed)
				}
			}
		}
	}
	logger.L().Info().Str("job", dbJob.Name).Str("run_id", runID).Msg("job rejected")
	return nil
}

func (o *Orchestrator) dispatchJob(
	ctx context.Context,
	runID string,
	projectID string,
	dbJob *db.PipelineJob,
	plJob *pipeline.Job,
	pl *pipeline.Pipeline,
	sshKey, repoURL, commitSHA string,
	projectEnv map[string]string,
	secrets map[string]string,
	jobMap map[string]*db.PipelineJob,
	attemptNumber int,
) error {
	if err := o.repo.UpdateJobStatus(ctx, dbJob.ID, db.StatusRunning); err != nil {
		return err
	}
	yamlStepOffset := 1
	if strings.TrimSpace(commitSHA) != "" {
		yamlStepOffset = 2
	}
	o.hub.Broadcast(runID, orchWs.Event{
		Type:           orchWs.EventJobUpdated,
		RunID:          runID,
		Status:         db.StatusRunning,
		JobName:        dbJob.Name,
		JobID:          dbJob.ID,
		JobStartedAtMs: time.Now().UnixMilli(),
	})

	steps := make([]sharedEvents.PipelineStep, len(plJob.Steps))
	for i, s := range plJob.Steps {
		steps[i] = sharedEvents.PipelineStep{
			Index:           yamlStepOffset + i,
			Name:            s.Name,
			Run:             s.Run,
			ContinueOnError: s.ContinueOnError,
			TimeoutSeconds:  s.Timeout.Seconds(),
			RetryMax:        s.Retry,
		}
	}

	mergedEnv := resolveEnvSecrets(pl.MergedEnv(plJob, projectEnv), secrets)

	var artifactsDownload []sharedEvents.ArtifactSource
	if plJob.Artifacts != nil {
		for _, dl := range plJob.Artifacts.Download {
			if dep, ok := jobMap[dl.Job]; ok {
				artifactsDownload = append(artifactsDownload, sharedEvents.ArtifactSource{
					JobName: dl.Job,
					JobID:   dep.ID,
				})
			}
		}
	}

	var cacheEntries []sharedEvents.CacheEntry
	if plJob.Cache != nil {
		cacheEntries = []sharedEvents.CacheEntry{{
			Key:   plJob.Cache.Key,
			Paths: plJob.Cache.Paths,
		}}
	}

	var artifactsUpload []string
	if plJob.Artifacts != nil {
		artifactsUpload = plJob.Artifacts.Paths
	}

	task := &sharedEvents.PipelineJobTask{
		RunID:             runID,
		ProjectID:         projectID,
		JobID:             dbJob.ID,
		JobName:           dbJob.Name,
		Image:             plJob.Image,
		RepoURL:           repoURL,
		CommitSHA:         commitSHA,
		SSHKey:            sshKey,
		Env:               mergedEnv,
		Secrets:           secrets,
		TimeoutSeconds:    plJob.Timeout.Seconds(),
		RetryMax:          plJob.Retry.Max,
		AttemptNumber:     attemptNumber,
		Cache:             cacheEntries,
		ArtifactsUpload:   artifactsUpload,
		ArtifactsDownload: artifactsDownload,
		Steps:             steps,
	}
	logger.L().Info().
		Str("job", dbJob.Name).Str("job_id", dbJob.ID).Int("attempt", attemptNumber).Str("run_id", runID).
		Msg("dispatching job")
	return o.jobPub.PublishJob(ctx, task)
}

func (o *Orchestrator) skipDependents(
	ctx context.Context,
	runID, failedName string,
	pl *pipeline.Pipeline,
	jobMap map[string]*db.PipelineJob,
) {
	for name, plJob := range pl.Jobs {
		dbJob, ok := jobMap[name]
		if !ok || dbJob.Status != db.StatusPending {
			continue
		}
		for _, need := range plJob.Needs {
			if need == failedName {
				logger.L().Info().Str("job", name).Str("failed", failedName).Msg("skipping job (dependency failed)")
				_ = o.repo.UpdateJobStatus(ctx, dbJob.ID, db.StatusSkipped)
				o.hub.Broadcast(runID, orchWs.Event{
					Type:    orchWs.EventJobUpdated,
					RunID:   runID,
					Status:  db.StatusSkipped,
					JobName: name,
					JobID:   dbJob.ID,
				})
				o.skipDependents(ctx, runID, name, pl, jobMap)
				break
			}
		}
	}
}

func retryDelay(attempt int) time.Duration {
	d := 5 * time.Second
	for i := 2; i < attempt; i++ {
		d *= 2
		if d > 60*time.Second {
			return 60 * time.Second
		}
	}
	return d
}

func allCompleted(needs []string, completed map[string]bool) bool {
	for _, n := range needs {
		if !completed[n] {
			return false
		}
	}
	return true
}

func jobNameByID(jobs []*db.PipelineJob, id string) string {
	for _, j := range jobs {
		if j.ID == id {
			return j.Name
		}
	}
	return ""
}

func jobMapByName(jobs []*db.PipelineJob) map[string]*db.PipelineJob {
	m := make(map[string]*db.PipelineJob, len(jobs))
	for _, j := range jobs {
		m[j.Name] = j
	}
	return m
}

func computeJobState(jobs []*db.PipelineJob) (completed map[string]bool, allDone, anyFailed bool) {
	completed = make(map[string]bool, len(jobs))
	allDone = true
	for _, j := range jobs {
		switch j.Status {
		case db.StatusSuccess:
			completed[j.Name] = true
		case db.StatusFailed:
			anyFailed = true
			completed[j.Name] = true
		case db.StatusSkipped:
			completed[j.Name] = true
		case db.StatusCancelled:
			anyFailed = true
			completed[j.Name] = true
		case db.StatusPending, db.StatusRunning, db.StatusAwaitingApproval:
			allDone = false
		}
	}
	return
}

type projectInfo struct {
	sshKey               string
	repoURL              string
	pipelineYAMLOverride string
	env                  map[string]string
	secrets              map[string]string
}

func (o *Orchestrator) getProjectInfo(ctx context.Context, projectID string) (*projectInfo, error) {
	resp, err := o.projectsClient.GetProjectInternal(ctx, &pb.GetProjectInternalRequest{ProjectId: projectID})
	if err != nil {
		if st, ok := status.FromError(err); ok && st.Code() == codes.NotFound {
			return nil, fmt.Errorf("project %s not found", projectID)
		}
		return nil, fmt.Errorf("GetProjectInternal grpc: %w", err)
	}
	env := make(map[string]string, len(resp.EnvVars))
	for _, v := range resp.EnvVars {
		env[v.Key] = v.Value
	}
	secrets := make(map[string]string, len(resp.Secrets))
	for _, s := range resp.Secrets {
		secrets[s.Key] = s.Value
	}
	return &projectInfo{
		sshKey:               resp.SshKey,
		repoURL:              resp.RepoUrl,
		pipelineYAMLOverride: resp.PipelineYamlOverride,
		env:                  env,
		secrets:              secrets,
	}, nil
}

var secretRefRe = regexp.MustCompile(`\$\{\{\s*secrets\.([A-Z_][A-Z0-9_]*)\s*\}\}`)

func resolveSecretRefs(s string, secrets map[string]string) string {
	return secretRefRe.ReplaceAllStringFunc(s, func(match string) string {
		sub := secretRefRe.FindStringSubmatch(match)
		if len(sub) < 2 {
			return match
		}
		if val, ok := secrets[sub[1]]; ok {
			return val
		}
		return match
	})
}

func resolveEnvSecrets(env map[string]string, secrets map[string]string) map[string]string {
	if len(secrets) == 0 {
		return env
	}
	out := make(map[string]string, len(env))
	for k, v := range env {
		out[k] = resolveSecretRefs(v, secrets)
	}
	return out
}

func maskSecretsFromLog(s string, secrets map[string]string) string {
	for _, v := range secrets {
		if len(v) >= 3 {
			s = strings.ReplaceAll(s, v, "***")
		}
	}
	return s
}

func (o *Orchestrator) publishRunCompleted(ctx context.Context, run *db.PipelineRun, finalStatus string) {
	if o.runEventPub == nil {
		return
	}
	dur := 0
	if run.StartedAt != nil && run.FinishedAt != nil {
		dur = int(run.FinishedAt.Sub(*run.StartedAt).Seconds())
		if dur < 0 {
			dur = 0
		}
	}
	finishedAt := time.Now().Unix()
	if run.FinishedAt != nil {
		finishedAt = run.FinishedAt.Unix()
	}
	startedAt := finishedAt
	if run.StartedAt != nil {
		startedAt = run.StartedAt.Unix()
	}
	evt := &sharedEvents.RunCompletedEvent{
		RunID:       run.ID,
		ProjectID:   run.ProjectID,
		Status:      finalStatus,
		Branch:      run.Branch,
		CommitSHA:   run.CommitSHA,
		DurationSec: dur,
		StartedAt:   startedAt,
		FinishedAt:  finishedAt,
	}
	if err := o.runEventPub.PublishRunCompleted(ctx, evt); err != nil {
		logger.L().Error().Err(err).Str("run_id", run.ID).Msg("analytics publish error")
	}
}

func cloneAndReadPipeline(ctx context.Context, repoURL, commitSHA, sshKey string) ([]byte, error) {
	dir, err := os.MkdirTemp("", "cicd-pipeline-*")
	if err != nil {
		return nil, fmt.Errorf("mktemp: %w", err)
	}
	defer os.RemoveAll(dir)

	keyFile := filepath.Join(dir, "key")
	if err := os.WriteFile(keyFile, []byte(sshKey), 0600); err != nil {
		return nil, fmt.Errorf("write ssh key: %w", err)
	}

	sshCmd := "ssh -i " + keyFile +
		" -o StrictHostKeyChecking=no" +
		" -o UserKnownHostsFile=/dev/null" +
		" -o BatchMode=yes"

	gitEnv := []string{
		"HOME=/tmp",
		"PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin",
		"GIT_SSH_COMMAND=" + sshCmd,
		"GIT_TERMINAL_PROMPT=0",
	}

	cloneDir := filepath.Join(dir, "repo")
	clone := exec.CommandContext(ctx, "git", "clone", "--depth=1", repoURL, cloneDir)
	clone.Env = gitEnv
	if out, err := clone.CombinedOutput(); err != nil {
		return nil, fmt.Errorf("git clone: %w (output: %s)", err, string(out))
	}

	if commitSHA != "" {
		fetch := exec.CommandContext(ctx, "git", "-C", cloneDir, "fetch", "--depth=1", "origin", commitSHA)
		fetch.Env = gitEnv
		_ = fetch.Run()

		checkout := exec.CommandContext(ctx, "git", "-C", cloneDir, "checkout", commitSHA)
		checkout.Env = gitEnv[:2]
		_ = checkout.Run()
	}

	return os.ReadFile(filepath.Join(cloneDir, ".cicd.yaml"))
}
