#!/bin/sh
set -eu

VAULT_ADDR="${VAULT_ADDR:-http://vault:8200}"
export VAULT_ADDR

echo "waiting for vault at ${VAULT_ADDR}..."
until vault status >/dev/null 2>&1; do
  sleep 1
done

POSTGRES_DSN="postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@postgres:5432/${POSTGRES_DB}?sslmode=disable"

echo "seeding vault KV secrets..."
vault kv put secret/cicd/shared \
  POSTGRES_USER="${POSTGRES_USER}" \
  POSTGRES_PASSWORD="${POSTGRES_PASSWORD}" \
  POSTGRES_DB="${POSTGRES_DB}" \
  POSTGRES_DSN="${POSTGRES_DSN}" \
  ORCHESTRATOR_INTERNAL_API_KEY="${ORCHESTRATOR_INTERNAL_API_KEY}" \
  CLICKHOUSE_PASSWORD="${CLICKHOUSE_PASSWORD}" \
  GRAFANA_ADMIN_USER="${GRAFANA_ADMIN_USER:-admin}" \
  GRAFANA_ADMIN_PASSWORD="${GRAFANA_ADMIN_PASSWORD}"

vault kv put secret/cicd/auth-service \
  JWT_ACCESS_SECRET="${JWT_ACCESS_SECRET}" \
  JWT_REFRESH_SECRET="${JWT_REFRESH_SECRET}" \
  REDIS_PASSWORD="${REDIS_PASSWORD:-}"

vault kv put secret/cicd/projects-service \
  PROJECTS_ENCRYPTION_KEY="${PROJECTS_ENCRYPTION_KEY}" \
  WEBHOOK_BASE_URL="${WEBHOOK_BASE_URL:-}"

vault kv put secret/cicd/ai-service \
  AI_LLM_API_KEY="${AI_LLM_API_KEY:-}" \
  AI_LLM_BASE_URL="${AI_LLM_BASE_URL:-https://gigachat.devices.sberbank.ru/api/v1}" \
  AI_LLM_MODEL="${AI_LLM_MODEL:-GigaChat}"

echo "installing vault policies..."
mkdir -p /tokens
for policy in /policies/*.hcl; do
  name=$(basename "$policy" .hcl)
  vault policy write "cicd-${name}" "$policy"
  token=$(vault token create -policy="cicd-${name}" -format=json | grep '"client_token"' | sed 's/.*"client_token"[[:space:]]*:[[:space:]]*"\([^"]*\)".*/\1/')
  if [ -z "$token" ]; then
    echo "failed to create token for cicd-${name}" >&2
    exit 1
  fi
  printf 'VAULT_TOKEN=%s\n' "$token" > "/tokens/${name}.env"
  echo "  policy cicd-${name} → /tokens/${name}.env"
done

echo "vault init complete"
