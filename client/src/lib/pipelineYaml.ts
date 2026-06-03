import yaml from 'js-yaml'

export interface PerformanceGateMetricRule {
  name: string
  direction?: 'lower_is_better' | 'higher_is_better' | string
  max?: number
  min?: number
}

export interface PerformanceGateJobConfig {
  sourceJob: string
  metrics?: PerformanceGateMetricRule[]
  baseline?: {
    window_days?: number
    min_samples?: number
    branch?: string
  }
  adaptive?: {
    enabled?: boolean
    sigma_factor?: number
    max_regression_pct?: number
  }
}

export function parsePerformanceGateJobs(
  pipelineYaml: string | undefined,
): Record<string, PerformanceGateJobConfig> {
  if (!pipelineYaml?.trim()) return {}
  try {
    const doc = yaml.load(pipelineYaml) as { jobs?: Record<string, Record<string, unknown>> } | null
    const jobs = doc?.jobs
    if (!jobs || typeof jobs !== 'object') return {}

    const out: Record<string, PerformanceGateJobConfig> = {}
    for (const [name, job] of Object.entries(jobs)) {
      const pg = job?.performance_gate as Record<string, unknown> | undefined
      if (!pg || typeof pg !== 'object') continue

      const metricsRaw = pg.metrics as Array<Record<string, unknown>> | undefined
      const metrics = Array.isArray(metricsRaw)
        ? metricsRaw
            .filter((m) => m && typeof m.name === 'string')
            .map((m) => ({
              name: String(m.name),
              direction: m.direction != null ? String(m.direction) : undefined,
              max: m.max != null ? Number(m.max) : undefined,
              min: m.min != null ? Number(m.min) : undefined,
            }))
        : undefined

      const baseline = pg.baseline as Record<string, unknown> | undefined
      const adaptive = pg.adaptive as Record<string, unknown> | undefined

      out[name] = {
        sourceJob: String(pg.source_job ?? ''),
        metrics,
        baseline: baseline
          ? {
              window_days: baseline.window_days != null ? Number(baseline.window_days) : undefined,
              min_samples: baseline.min_samples != null ? Number(baseline.min_samples) : undefined,
              branch: baseline.branch != null ? String(baseline.branch) : undefined,
            }
          : undefined,
        adaptive: adaptive
          ? {
              enabled: adaptive.enabled != null ? Boolean(adaptive.enabled) : undefined,
              sigma_factor: adaptive.sigma_factor != null ? Number(adaptive.sigma_factor) : undefined,
              max_regression_pct:
                adaptive.max_regression_pct != null ? Number(adaptive.max_regression_pct) : undefined,
            }
          : undefined,
      }
    }
    return out
  } catch {
    return {}
  }
}

export function isPerformanceGateJob(
  jobName: string,
  gateJobs: Record<string, PerformanceGateJobConfig>,
): boolean {
  return Object.prototype.hasOwnProperty.call(gateJobs, jobName)
}
