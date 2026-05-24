package handler

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/vsuaiqq/cicd/ai-service/internal/llm"
)

const (
	maxFailedLogLines    = 200
	maxPassedLogLines    = 15
	maxTotalContextBytes = 40_000
)

type StepCtx struct {
	Name      string
	Status    string
	ExitCode  int
	LogOutput string
}

type AnalyzeRequest struct {
	JobName      string
	JobStatus    string
	PipelineYAML string
	Steps        []StepCtx
	Lang         string
}

type AnalysisResult struct {
	Summary       string   `json:"summary"`
	RootCause     string   `json:"root_cause"`
	Fix           string   `json:"fix"`
	RelevantLines []string `json:"relevant_lines"`
}

const systemPromptRU = `Ты — эксперт по CI/CD системам. Анализируй ошибки пайплайнов и объясняй их на русском языке.

По логам упавшего джоба определи что именно пошло не так.

Ответь ТОЛЬКО валидным JSON-объектом в точно таком формате (без markdown, без обёртки в кавычки):
{
  "summary": "одно предложение — что именно упало",
  "root_cause": "2-3 предложения — техническая причина ошибки со ссылками на конкретные сообщения об ошибках и exit-коды",
  "fix": "конкретные пронумерованные шаги для исправления проблемы одной строкой через \n",
  "relevant_lines": ["точная строка из лога 1", "точная строка из лога 2"]
}

Поле fix — строка с нумерованными шагами, НЕ массив.
Массив relevant_lines — до 5 наиболее информативных строк ошибки, скопированных дословно из логов.`

const systemPromptEN = `You are a CI/CD expert. Analyze pipeline failures and explain them in English.

Based on the failed job logs, determine what went wrong.

Reply ONLY with a valid JSON object in exactly this format (no markdown, no code fences):
{
  "summary": "one sentence — what exactly failed",
  "root_cause": "2-3 sentences — technical root cause with references to specific error messages and exit codes",
  "fix": "concrete numbered steps to fix the issue as a single string using \n",
  "relevant_lines": ["exact log line 1", "exact log line 2"]
}

The fix field is a string with numbered steps, NOT an array.
The relevant_lines array — up to 5 most informative error lines, copied verbatim from the logs.`

func BuildMessages(req AnalyzeRequest) []llm.Message {
	systemPrompt := systemPromptRU
	if req.Lang == "en" {
		systemPrompt = systemPromptEN
	}

	return []llm.Message{
		{Role: "system", Content: systemPrompt},
		{Role: "user", Content: buildUserMessage(req)},
	}
}

func ParseAnalysis(content string) (*AnalysisResult, error) {
	content = strings.TrimSpace(content)
	if strings.HasPrefix(content, "```") {
		if i := strings.Index(content[3:], "\n"); i >= 0 {
			content = content[3+i+1:]
		}
		content = strings.TrimSuffix(strings.TrimSpace(content), "```")
	}

	var raw map[string]json.RawMessage
	if err := json.Unmarshal([]byte(content), &raw); err != nil {
		return nil, fmt.Errorf("parse analysis JSON: %w", err)
	}

	result := &AnalysisResult{}

	if v, ok := raw["summary"]; ok {
		_ = json.Unmarshal(v, &result.Summary)
	}
	if v, ok := raw["root_cause"]; ok {
		_ = json.Unmarshal(v, &result.RootCause)
	}
	if v, ok := raw["relevant_lines"]; ok {
		_ = json.Unmarshal(v, &result.RelevantLines)
	}

	if v, ok := raw["fix"]; ok {
		var s string
		if err := json.Unmarshal(v, &s); err == nil {
			result.Fix = s
		} else {
			var arr []string
			if err := json.Unmarshal(v, &arr); err == nil {
				for i, step := range arr {
					arr[i] = fmt.Sprintf("%d. %s", i+1, step)
				}
				result.Fix = strings.Join(arr, "\n")
			}
		}
	}

	return result, nil
}

type GenerateRequest struct {
	Description string
}

func BuildGenerateMessages(req GenerateRequest) []llm.Message {
	const systemPrompt = `Ты — эксперт по CI/CD. Генерируй конфигурацию пайплайна в формате .cicd.yaml на основе описания проекта.

Формат .cicd.yaml:
` + "```yaml" + `
name: Название пайплайна

on:
  push:
    branches:
      - main          # ветки-триггеры

env:                  # глобальные переменные (опционально)
  KEY: value

jobs:
  job-name:           # уникальное имя джоба (только латиница, цифры, дефис)
    image: golang:1.21  # Docker-образ для выполнения
    needs: []           # зависимости от других джобов (имена)
    steps:
      - name: Шаг
        run: |
          команды
        continue_on_error: false
` + "```" + `

Правила:
- Поля name, run обязательны для каждого шага
- image — любой публичный Docker-образ
- needs — массив имён джобов от которых зависит этот
- continue_on_error: true если шаг не должен останавливать джоб при ошибке
- В run используй многострочный YAML (|)
- Возвращай ТОЛЬКО текст YAML без объяснений, без markdown-обёртки, без ` + "```yaml" + ``

	userMsg := "Описание проекта: " + req.Description + "\n\nСгенерируй .cicd.yaml:"

	return []llm.Message{
		{Role: "system", Content: systemPrompt},
		{Role: "user", Content: userMsg},
	}
}

func buildUserMessage(req AnalyzeRequest) string {
	var b strings.Builder

	fmt.Fprintf(&b, "Failed pipeline job: %q (status: %s)\n\n", req.JobName, req.JobStatus)

	b.WriteString("=== Steps ===\n")
	totalBytes := 0
	for _, step := range req.Steps {
		maxLines := maxPassedLogLines
		if step.Status == "failed" || step.ExitCode != 0 {
			maxLines = maxFailedLogLines
		}
		logContent := tailLines(step.LogOutput, maxLines)

		fmt.Fprintf(&b, "--- Step: %s | status: %s | exit_code: %d ---\n",
			step.Name, step.Status, step.ExitCode)
		b.WriteString(logContent)
		b.WriteString("\n\n")

		totalBytes += len(logContent)
		if totalBytes >= maxTotalContextBytes {
			b.WriteString("[remaining logs truncated to stay within context limit]\n")
			break
		}
	}

	if req.PipelineYAML != "" {
		b.WriteString("=== Pipeline configuration (.cicd.yaml) ===\n")
		b.WriteString(truncate(req.PipelineYAML, 3000))
		b.WriteString("\n")
	}

	return b.String()
}

func tailLines(s string, n int) string {
	if s == "" {
		return "(no output)"
	}
	lines := strings.Split(strings.TrimRight(s, "\n"), "\n")
	if len(lines) <= n {
		return s
	}
	skipped := len(lines) - n
	return fmt.Sprintf("[... %d earlier lines omitted ...]\n%s", skipped, strings.Join(lines[skipped:], "\n"))
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "\n[truncated]"
}
