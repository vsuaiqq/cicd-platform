export const helpDocs = {
  overview: {
    title: '.cicd.yaml reference',
    desc:
      'Place a `.cicd.yaml` file at the root of your repository to define your pipeline. Flow reads this file when a push webhook matches the configured branch filter and executes jobs in dependency order.',
    tip:
      'You can override the repository file from Settings → Pipeline config. The UI-configured YAML takes precedence over `.cicd.yaml` in the repo.',
    minimalExample: 'Minimal working example:',
  },
  global: {
    title: 'Global fields',
    desc: 'Top-level keys supported in `.cicd.yaml` today.',
    fields: {
      name: 'Human-readable name for the pipeline. Shown in run detail and logs.',
      on: 'Trigger configuration. Currently only `on.push` is supported — see the triggers section.',
      env:
        'Key/value pairs injected as environment variables into every job and step. Overridden by job-level `env`. Project secrets and Settings → Environment variables are merged separately (see Secrets & env).',
      jobs: 'Dictionary of jobs to run. Keys are internal job IDs (alphanumeric + hyphens).',
    },
  },
  on: {
    title: 'on — triggers',
    desc:
      'Defines when the pipeline runs. Only `on.push` is implemented today. Other trigger types (`pull_request`, `manual`, `schedule`) and path filters are not parsed — they are ignored if present in YAML.',
    fields: {
      push:
        'Runs on push webhooks from the connected git host. Use `branches` to restrict which source branches start a run.',
      branches:
        'Branch name list (`main`, `develop`, `*`). If omitted or empty, any branch triggers a run. Exact match only — glob patterns like `release/*` are not expanded yet.',
    },
    exampleNote:
      'A push to a branch not listed under `branches` does not start a pipeline run.',
  },
  jobs: {
    title: 'jobs',
    desc:
      'A map of job definitions. Each key is a unique job ID used in `needs` lists. Jobs with no dependencies run in parallel; jobs with `needs` wait until all listed jobs finish successfully (failed or skipped upstream jobs block dependents).',
    tip:
      'A job is either a runner job (has `steps`, executed in a container or on the host) or a performance gate job (has `performance_gate`, evaluated by the orchestrator — no steps, no container).',
  },
  jobFields: {
    title: 'Job fields',
    desc: 'Fields available inside each job definition. Use `steps` for ordinary jobs or `performance_gate` for adaptive quality gates — not both.',
    fields: {
      name: 'Display name shown in the UI. Falls back to the job ID if omitted.',
      image:
        'Docker image for this job (`golang:1.25`, `node:20`, …). Optional — when omitted, steps run directly on the runner host (no Docker, no cache bind-mounts). Required for `cache`.',
      needs:
        'List of job IDs that must finish before this job starts. Creates the dependency graph. Circular dependencies are rejected at parse time. For `performance_gate`, `source_job` must appear here.',
      env:
        'Environment variables scoped to this job. Merges with (and overrides) pipeline `env`. Values may reference project secrets via `${{ secrets.NAME }}`.',
      timeout:
        'Maximum wall-clock time for the entire job. Go duration strings: `5m`, `1h30m`. Defaults to no limit.',
      retry:
        'Extra attempts on job failure. Shorthand `retry: 2` or long form `retry: { max: 2 }` — up to `max + 1` total attempts (initial run plus retries).',
      approval:
        'When `required`, `true`, `yes`, or `1`, the job pauses in awaiting approval until a project editor approves it from the run detail page.',
      cache: 'Filesystem cache between runs (Docker jobs only). See the cache section.',
      artifacts: 'Upload paths and/or download artifacts from dependency jobs. See the artifacts section.',
      steps: 'Ordered list of steps for runner jobs. Required unless the job defines `performance_gate`.',
      performance_gate:
        'Adaptive performance quality gate — see the Performance gate section. The orchestrator evaluates metrics from `source_job` internally; no runner steps are executed.',
    },
  },
  steps: {
    title: 'steps',
    desc:
      'Each runner job contains an ordered list of steps. Steps run sequentially in the same container (or on the host when `image` is omitted). A failing step aborts subsequent steps unless `continue_on_error` is set.',
    fields: {
      name: 'Step label displayed in the UI and logs.',
      run:
        'Shell script to execute. Multi-line scripts use the YAML block literal (`|`). Runs under `/bin/sh -e` by default.',
      timeout: 'Maximum time allowed for this step (`5m`, `30s`, …).',
      retry: 'Extra attempts for this step on failure. Up to `retry + 1` total attempts.',
      continue_on_error:
        'If `true`, the job continues even if this step exits non-zero. The step is marked failed but later steps still run.',
    },
  },
  cache: {
    title: 'cache',
    desc:
      'Caches directories between runs using a content-addressed key. On a cache hit paths are restored before steps run; on a miss they are saved after completion. Requires a Docker `image` — host jobs ignore `cache`.',
    fields: {
      key:
        'Cache identifier. Supports `${{ checksum "relative/path" }}` — SHA-256 of the file (first 16 hex chars), resolved on the runner from the checked-out repo.',
      paths:
        'Filesystem paths to cache inside the container. Relative paths resolve from `/workspace`.',
    },
    tip:
      'Cache keys are shared across branches on the same runner. Include a branch name or commit qualifier in the key when isolation is needed (plain text in the key string — `${{ branch }}` is not supported).',
  },
  artifacts: {
    title: 'artifacts',
    desc:
      'Files produced by a job (`paths`) are archived after steps complete successfully. Other jobs in the same run can fetch them via `download` before their steps execute. Artifacts are scoped to a single run — there is no cross-run artifact store or download button in the UI yet.',
    fields: {
      paths: 'Paths to archive after the job succeeds. Relative to the workspace root.',
      download:
        'List of `{ job: <jobId> }` entries. Artifacts from those jobs are extracted into the workspace before steps run.',
    },
  },
  loadTesting: {
    title: 'Load testing',
    desc:
      'Load testing is an ordinary runner job — any job with `steps` can run a load test. There is no special job type. After the job finishes, the runner looks for `.flow/perf-metrics.json` in the workspace and attaches the metrics to the job result.',
    fields: {
      metricsFile:
        'Path: `.flow/perf-metrics.json` (fixed). The runner reads this file automatically; you do not configure the path in YAML.',
      format:
        'JSON object with `version` (optional, default 1), optional `tool` string, and a `metrics` map of numeric values.',
      exampleMetrics:
        'Common keys: `http_req_duration_p95`, `http_req_duration_avg`, `http_reqs`, `http_req_failed_rate`. Names are free-form — the performance gate references them by name.',
    },
    tip:
      'A load test job works without a performance gate. Add a `performance_gate` job downstream when you want automatic pass/fail decisions based on metrics.',
    exampleNote: 'See `test-repo/scripts/load-test.sh` in the platform repository for a working example that writes the metrics file.',
  },
  performanceGate: {
    title: 'Performance gate',
    desc:
      'A performance gate is a platform job type that compares load-test metrics from a source job against constant and/or adaptive thresholds. It runs inside the orchestrator (via analytics-service) — no container, no steps. On failure, downstream jobs are skipped.',
    fields: {
      source_job:
        'Job ID of the runner job that produced `.flow/perf-metrics.json`. Must also appear in this job\'s `needs` list. Cannot be another performance gate job.',
      metrics:
        'List of metric rules. Each entry has `name`, `direction` (`lower_is_better` or `higher_is_better`), and optional constant bounds: `max` for lower-is-better, `min` for higher-is-better.',
      baseline_window_days:
        'How many days of historical runs to sample for adaptive thresholds. Default: `30`.',
      baseline_min_samples:
        'Minimum historical samples per metric before adaptive checks apply. Default: `3`. Below this threshold the gate is in cold start — adaptive checks are skipped, constant checks still apply.',
      baseline_branch:
        'Optional branch filter for baseline history. Defaults to the current run branch.',
      adaptive_enabled:
        'When `true` (default), adaptive thresholds are computed from baseline mean and standard deviation. Set `false` for constant-only gates.',
      adaptive_sigma_factor:
        'Multiplier for standard deviation in adaptive bounds. Default: `2.0`.',
      adaptive_max_regression_pct:
        'Maximum allowed regression vs baseline mean (percent). Default: `15`.',
    },
    adaptiveFormula:
      'lower_is_better: threshold = min(μ + k·σ, μ × (1 + max_regression% / 100)). higher_is_better: threshold = max(μ − k·σ, μ × (1 − max_regression% / 100)). A metric passes only if it satisfies both constant bounds (when configured) and the effective threshold.',
    coldStart:
      'During cold start (fewer than `min_samples` historical values), adaptive checks are skipped and only constant `max`/`min` limits apply. The run detail page shows a cold-start badge.',
    uiNote:
      'Gate results appear on the run detail page: per-metric verdict, baseline stats, constant vs adaptive thresholds, and overall pass/fail.',
    defaultMetrics:
      'If `metrics` is omitted, defaults are: `http_req_duration_p95` (lower), `http_req_failed_rate` (lower), `http_reqs` (higher) — with no constant bounds.',
  },
  secrets: {
    title: 'Secrets & environment',
    desc: 'Flow resolves environment variables in the following priority order (highest wins):',
    priorityList: [
      'Secrets — set in Settings → Secrets. Always highest priority; values masked in logs.',
      'Job env — `env` inside a job definition.',
      'Pipeline env — top-level `env` in `.cicd.yaml`.',
      'Project env vars — Settings → Environment variables. Lowest priority.',
    ],
    interpolationTitle: 'Supported interpolations',
    interpolationList: [
      '`${{ secrets.NAME }}` — in job or pipeline `env` values only. Resolved by the orchestrator before dispatch.',
      '`${{ checksum "path" }}` — in `cache.key` only. Resolved by the runner from the checked-out repo.',
      '`$NAME` — in shell `run` scripts for any env var or secret already injected into the job environment.',
    ],
    warn:
      'Secrets are write-only after saving. `${{ branch }}`, `${{ sha }}`, `${{ matrix.* }}`, and job/step `if:` expressions are not supported yet.',
  },
  limitations: {
    title: 'Limitations & roadmap',
    desc: 'Current platform behaviour and features not yet implemented:',
    items: [
      'Triggers: only `on.push` with optional `branches`. No pull_request, manual trigger block, schedule, or path filters.',
      'Conditions: no job or step `if:` expressions.',
      'Matrix: no `strategy.matrix` job fan-out.',
      'Environments: no top-level `environments` block or job `environment:` field.',
      'allow_failure: not supported — a failed job fails the run (except `continue_on_error` on individual steps).',
      'Step env: per-step `env` is not parsed from YAML.',
      'Cache works only when the job has a Docker `image`. Host jobs skip caching.',
      'Artifacts are available only within the same pipeline run (via `artifacts.download`). No UI download button yet.',
      'Logs: the UI shows step output after the job completes; live streaming is backend-only today.',
      'Performance gate cold start: adaptive thresholds require enough historical runs on the same branch; until then only constant `max`/`min` apply.',
      'Metrics file: fixed path `.flow/perf-metrics.json`; custom paths are not configurable.',
    ],
  },
  example: {
    title: 'Full example',
    desc:
      'A Go project with lint, test, build, staging deploy, load test, adaptive performance gate, and production deploy. This matches `test-repo/.cicd.yaml` in the platform repository.',
  },
} as const
