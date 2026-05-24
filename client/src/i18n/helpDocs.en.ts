export const helpDocs = {
  overview: {
    title: '.cicd.yaml reference',
    desc:
      'Place a `.cicd.yaml` file at the root of your repository to define your pipeline. Flow reads this file when a matching trigger fires (push, pull request, manual run, or cron schedule) and executes jobs in dependency order.',
    tip:
      'You can override the repository file from Settings → Pipeline config. The UI-configured YAML takes precedence. From the Runs tab, editors can start a manual run, rerun a completed run, or rerun only failed jobs.',
    minimalExample: 'Minimal working example:',
  },
  global: {
    title: 'Global fields',
    desc: 'Top-level keys that apply to the entire pipeline file.',
    fields: {
      name: 'Human-readable name for the pipeline. Shown in run detail and logs.',
      on: 'Trigger configuration. See the on (triggers) section.',
      env:
        'Key/value pairs injected as environment variables into every job and step. Lowest precedence inside `.cicd.yaml` — overridden by job `env`, step `env`, then secrets.',
      environments:
        'Named deployment targets referenced by jobs via `environment:`. Each entry can define a URL and protection rules (branch filter, wait timer, required reviewers). See the environments section.',
      jobs: 'Dictionary of jobs to run. Keys are internal job IDs (alphanumeric + hyphens).',
    },
  },
  on: {
    title: 'on — triggers',
    desc:
      'Defines when the pipeline should run. You can combine `push`, `pull_request`, `manual`, and `schedule` in the same file — each matching event starts a new run.',
    fields: {
      push:
        'Runs on push webhooks. Filter by `branches`, `paths`, and `paths_ignore`. Path filters inspect files changed in the pushed commits.',
      pull_request:
        'Runs on pull-request webhooks. Filter by target `branches`, event `types`, `drafts`, and path filters on changed files.',
      manual:
        'Declares inputs for manual runs started from the Runs tab. Each input is exposed to steps as `MANUAL_INPUT_<NAME>` (uppercase).',
      schedule:
        'Cron-driven runs. Each list entry needs a 5-field `cron` expression and optional `branch` / `timezone` (default UTC).',
      branches:
        'Branch name globs (`main`, `release/*`). For `push` — source branch; for `pull_request` — target (base) branch. Omitted = any branch.',
      paths:
        'Glob list (`**`, `*`, leading `!` for negation). Run only when at least one changed file matches after excludes.',
      paths_ignore:
        'Inverse filter: skip the run when every changed file matches one of these globs (e.g. docs-only commits).',
      types:
        'PR event types: `opened`, `synchronize`, `reopened` (default if omitted).',
      drafts:
        'When `false` (default), draft pull requests are skipped. Note: GitHub draft detection may be unreliable depending on the git host.',
      inputs:
        'Under `manual.inputs.<name>`: `description`, `default`, and `required` (UI hint only).',
      cron: 'Standard 5-field cron: `minute hour day-of-month month day-of-week`.',
      timezone: 'IANA time zone for the cron entry (e.g. `Europe/Amsterdam`). Default: `UTC`.',
    },
    pathTip:
      'Path filters on `push` use file lists from the webhook. For `pull_request`, path filtering works only when the git host provides changed-file metadata.',
  },
  jobs: {
    title: 'jobs',
    desc:
      'A map of job definitions. Each key is a unique job ID used in `needs` lists. Jobs with no dependencies run in parallel; jobs with `needs` wait until all listed jobs finish (success or skipped). A matrix job fans out into parallel variants — see the matrix section.',
    tip:
      'Jobs that share no `needs` relationships run in parallel. Use `needs` to fan-in multiple parallel jobs before a final stage.',
  },
  jobFields: {
    title: 'Job fields',
    desc: 'Fields available inside each job definition.',
    fields: {
      name: 'Display name shown in the UI. Falls back to the job ID if omitted.',
      image:
        'Docker image for this job (`golang:1.22`, `node:20`, …). Optional — when omitted, steps run directly on the runner host (no Docker, no cache bind-mounts).',
      needs:
        'List of job IDs that must finish before this job starts. Creates the dependency graph. Circular dependencies are rejected at parse time.',
      env:
        'Environment variables scoped to this job. Merges with (and overrides) pipeline `env`. Supports `${{ secrets.NAME }}`, `${{ matrix.KEY }}`, `${{ branch }}`, `${{ sha }}`.',
      timeout:
        'Maximum wall-clock time for the entire job. Go duration strings: `5m`, `1h30m`. Defaults to no limit.',
      retry:
        'Extra attempts on failure. Shorthand `retry: 2` or long form `retry: { max: 2 }` — up to 3 total attempts.',
      allow_failure:
        'If `true`, a failing job does not fail the overall pipeline. Downstream jobs that `needs` this one still run.',
      manual:
        'If `true`, the job pauses for manual approval before executing. Equivalent to `approval: required`.',
      approval:
        'Alias for `manual`. Accepts `required`, `true`, `yes`, or `1`.',
      strategy:
        'Matrix expansion — see the matrix section. Produces parallel job variants with `MATRIX_<KEY>` env vars.',
      if:
        'Boolean expression — when false the job is skipped (dependents treat it as succeeded). Evaluated at dispatch time after dependencies finish. See conditions.',
      environment:
        'Name of an entry from top-level `environments`. Applies protection rules and sets `DEPLOY_ENV` / `DEPLOY_URL` env vars.',
      cache: 'Filesystem cache between runs (Docker jobs only). See the cache section.',
      artifacts: 'Upload paths and/or download artifacts from dependency jobs. See the artifacts section.',
      steps: 'Ordered list of steps to execute inside the job.',
    },
  },
  steps: {
    title: 'steps',
    desc:
      'Each job contains an ordered list of steps. Steps run sequentially in the same container (or on the host when `image` is omitted). A failing step aborts subsequent steps unless `continue_on_error` is set.',
    fields: {
      name: 'Step label displayed in the UI and logs.',
      run:
        'Shell script to execute. Multi-line scripts use the YAML block literal (`|`). Runs under `/bin/sh -e` by default.',
      env: 'Environment variables scoped to this step only. Overrides job and pipeline env.',
      timeout: 'Maximum time allowed for this step. Overrides job-level timeout for this step.',
      retry: 'Extra attempts for this step on failure. Combined with job-level retry.',
      continue_on_error:
        'If `true`, the job continues even if this step exits non-zero. The step is marked failed but later steps still run.',
      if:
        'Boolean expression — when false the step is recorded as skipped and never executed. Same syntax as job `if:`, evaluated at dispatch time (not after prior steps). See conditions.',
    },
  },
  matrix: {
    title: 'strategy — matrix',
    desc:
      'Expands one job definition into multiple parallel variants. Each combination gets its own run with `MATRIX_<KEY>=<value>` env vars and `${{ matrix.<key> }}` resolved in strings (image, env, cache key, step commands, …).',
    fields: {
      matrix:
        'Map of variable names to value lists. Cartesian product across keys defines variants.',
      include:
        'Extra explicit combinations added on top of the matrix product.',
      exclude:
        'Combinations removed from the matrix product.',
      fail_fast:
        'When `true`, the first failing variant cancels remaining siblings in the same matrix group.',
    },
    tip:
      'Matrix job IDs in the UI appear as `<jobId> (<value1>, <value2>, …)`. Downstream `needs` should reference the parent job ID — the orchestrator waits for all variants.',
  },
  conditions: {
    title: 'if — conditions',
    desc:
      'Job and step `if:` fields accept boolean expressions. When the expression is false, the job/step is marked `skipped`. Expressions are evaluated at dispatch time — they cannot inspect stdout or exit codes from earlier steps in the same job.',
    fields: {
      literals: 'String and boolean literals in single quotes: `branch == \'main\'`, `event != \'pull_request\'`.',
      context:
        'Context variables: `branch`, `sha`, `event` (`push`, `pull_request`, `manual`, `schedule`), `env.NAME`, `matrix.NAME`, `needs.<jobId>.status`, `needs.<jobId>.result`.',
      functions:
        'Status helpers: `success()`, `failure()`, `cancelled()`, `always()`. Combine with `==`, `!=`, `&&`, `||`, `!`.',
    },
    warn:
      'Because `if:` is evaluated at dispatch time, expressions like `success()` on a step inside the same job always reflect the planned run, not live step outcomes. Use job-level `if:` with `needs.<job>.result` to gate on upstream jobs.',
  },
  environments: {
    title: 'environments',
    desc:
      'Top-level `environments` defines named deployment targets. Jobs opt in with `environment: <name>`. Protection rules are additive — all configured rules must pass before the job dispatches.',
    fields: {
      url:
        'Deployment URL. Auto-injected as `DEPLOY_URL`; `DEPLOY_ENV` is set to the environment name.',
      branches:
        'Source branch globs allowed to deploy. Empty = any branch.',
      wait_timer:
        'Minimum delay after the job becomes runnable before dispatch (e.g. `5m`, `30s`).',
      required_reviewers:
        'User IDs allowed to approve. When non-empty, the job pauses in awaiting_approval like `manual: true`.',
      reviewer_count:
        'Parsed from YAML; currently only a single approver from the list is enforced (N-of-M is not yet implemented).',
    },
    tip:
      'Environment protection stacks with job `manual:` / `approval:`. A job can use both — e.g. `environment: production` with `manual: true` for a double gate.',
  },
  cache: {
    title: 'cache',
    desc:
      'Caches directories between runs using a content-addressed key. On a cache hit paths are restored before steps run; on a miss they are saved after completion. Cache requires a Docker `image` — host jobs ignore `cache`.',
    fields: {
      key:
        'Cache identifier. Supports `${{ checksum "file" }}`, `${{ branch }}`, `${{ matrix.KEY }}`.',
      paths:
        'Filesystem paths to cache inside the container. Relative paths resolve from `/workspace`.',
    },
    tip:
      'Cache keys are shared across branches on the same runner. Include a branch or matrix qualifier in the key when isolation is needed.',
  },
  artifacts: {
    title: 'artifacts',
    desc:
      'Files produced by a job (`paths`) are archived after steps complete. Other jobs in the same run can fetch them via `download` before their steps execute. Artifacts are scoped to a single run — there is no cross-run artifact store or download button in the UI yet.',
    fields: {
      paths: 'Glob patterns or paths to archive after the job succeeds. Relative to the workspace root.',
      download:
        'List of `{ job: <jobId> }` entries. Artifacts from those jobs are extracted into the workspace before steps run.',
    },
  },
  expressions: {
    title: 'Template expressions',
    desc:
      'String interpolation inside `${{ ... }}` in env values, cache keys, image names, and step commands.',
    fields: {
      checksum:
        'SHA-256 checksum of the given file (hex). Resolved on the runner from the checked-out repo.',
      branch: 'Name of the branch that triggered this run.',
      sha: 'Full commit SHA of the triggering event.',
      secrets:
        'Value of a project secret named `NAME`. Use inside `env` values. Secret values are masked in logs.',
      matrix:
        'Matrix axis value for the current variant: `${{ matrix.go }}`, `${{ matrix.os }}`, etc.',
    },
  },
  secrets: {
    title: 'Secrets & environment',
    desc: 'Flow resolves environment variables in the following priority order (highest wins):',
    priorityList: [
      'Secrets — set in Settings → Secrets. Always highest priority; values masked in logs.',
      'Step env — `env` inside a specific step.',
      'Job env — `env` inside a job definition.',
      'Pipeline env — top-level `env` in `.cicd.yaml`.',
      'Project env vars — Settings → Environment variables. Lowest priority.',
      'Manual inputs — `on.manual.inputs` appear as `MANUAL_INPUT_<NAME>` (uppercase) at trigger time.',
      'Environment vars — `DEPLOY_ENV` and `DEPLOY_URL` when `environment:` is set.',
    ],
    warn:
      'Secrets are write-only after saving. Reference them with `${{ secrets.NAME }}` in `env` or as `$NAME` in shell commands.',
  },
  limitations: {
    title: 'Known limitations',
    desc: 'Current behaviour you should be aware of when authoring pipelines:',
    items: [
      'Job/step `if:` is evaluated at dispatch time, not after each step finishes inside the same job.',
      'Cache works only when the job has a Docker `image`. Host jobs skip caching.',
      'Artifacts are available only within the same pipeline run (via `artifacts.download`).',
      'Live log streaming exists on the backend; the UI currently shows step output after the job completes.',
      'PR path filters depend on changed-file metadata from the git host — not all hosts populate it reliably.',
      'GitHub draft PR detection may report `draft: false` due to API limitations.',
      '`reviewer_count` is parsed but only one approver from `required_reviewers` is enforced today.',
      'Manual run inputs marked `required` are enforced in the UI only, not on the server.',
    ],
  },
  example: {
    title: 'Full example',
    desc:
      'A Go project with push/PR/manual/schedule triggers, matrix tests, caching, artifacts, environment protection, and conditional deploy. See also `test-repo/.cicd.yaml` in the platform repository.',
  },
} as const
