import { apiRequest } from './client'

const PREFIX = '/api/v1/ai'

export interface AnalysisResult {
  summary:        string
  root_cause:     string
  fix:            string
  relevant_lines: string[]
}

export function generatePipeline(description: string, token: string): Promise<string> {
  return apiRequest<{ yaml: string }>(`${PREFIX}/generate-pipeline`, {
    method: 'POST',
    body:   JSON.stringify({ description }),
    token,
  }).then((r) => r.yaml ?? '')
}

export function analyzeFailure(
  runId: string,
  jobId: string,
  token: string,
  lang = 'en',
): Promise<AnalysisResult> {
  return apiRequest<AnalysisResult>(`${PREFIX}/analyze-failure`, {
    method: 'POST',
    body:   JSON.stringify({ run_id: runId, job_id: jobId, lang }),
    token,
  })
}
