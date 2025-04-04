
# Arquivo .env
ENV_FILE = .env
# Extrai as variáveis do .env e transforma em flags --var DB_USER=value
atlas_vars = $(shell grep -E '^(DB_USER|DB_PASSWORD|DB_HOST|DB_PORT|DB_NAME|DB_DEV_NAME)=' $(ENV_FILE) | sed 's/^/--var /' | sed 's/=/=/' )

dev-api:
	air -c .air.toml

dev-consumer:
	air -c .air-consumer-debts.toml

ent-generate:
	go run entgo.io/ent/cmd/ent generate ./pkg/ent/schema

## Status das migrations
atlas-status:
	atlas migrate status --env local $(atlas_vars)

## Aplica migrations
atlas-up:
	atlas migrate apply --env local $(atlas_vars)

## Reverte a última migration
atlas-down:
	atlas migrate down --env local $(atlas_vars)

## Reverte todas as migrations
atlas-reset:
	atlas migrate down --env local --all $(atlas_vars)

## Cria nova migration (uso: make atlas-new NAME=descricao)
atlas-new:
	atlas migrate diff $${NAME} --env local $(atlas_vars) --to ent://pkg/ent/schema

## Gera snapshot do banco atual
atlas-snapshot:
	atlas migrate snapshot "initial" --env local $(atlas_vars)
