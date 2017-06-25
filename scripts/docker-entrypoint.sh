#!/bin/sh
set -e

: ${MNEMOSYNED_PORT:=8080}
: ${MNEMOSYNED_HOST:=0.0.0.0}
: ${MNEMOSYNED_TTL:=24m}
: ${MNEMOSYNED_TTC:=1m}
: ${MNEMOSYNED_CLUSTER_LISTEN:=$(hostname):$MNEMOSYNED_PORT}
: ${MNEMOSYNED_CLUSTER_SEEDS:=}
: ${MNEMOSYNED_LOGGER_ENVIRONMENT:=production}
: ${MNEMOSYNED_LOGGER_LEVEL:=6}
: ${MNEMOSYNED_STORAGE:=postgres}
: ${MNEMOSYNED_MONITORING:=false}
: ${MNEMOSYNED_POSTGRES_ADDRESS:=postgres://postgres:postgres@postgres/postgres?sslmode=disable}
: ${MNEMOSYNED_POSTGRES_TABLE:=session}
: ${MNEMOSYNED_POSTGRES_SCHEMA:=mnemosyne}
: ${MNEMOSYNED_POSTGRES_DEBUG:=false}
: ${MNEMOSYNED_TLS_ENABLED:=false}

if [ "$1" = 'mnemosyned' ]; then
exec mnemosyned -host=${MNEMOSYNED_HOST} \
	-port=${MNEMOSYNED_PORT} \
	-ttl=${MNEMOSYNED_TTL} \
	-ttc=${MNEMOSYNED_TTC} \
	-cluster.listen=${MNEMOSYNED_CLUSTER_LISTEN} \
	-cluster.seeds=${MNEMOSYNED_CLUSTER_SEEDS} \
	-storage=${MNEMOSYNED_STORAGE} \
	-logger.environment=${MNEMOSYNED_LOGGER_ENVIRONMENT} \
	-logger.level=${MNEMOSYNED_LOGGER_LEVEL} \
	-monitoring=${MNEMOSYNED_MONITORING} \
	-postgres.address=${MNEMOSYNED_POSTGRES_ADDRESS} \
	-postgres.table=${MNEMOSYNED_POSTGRES_TABLE} \
	-postgres.schema=${MNEMOSYNED_POSTGRES_SCHEMA} \
	-tls=${MNEMOSYNED_TLS_ENABLED} \
	-tls.cert=${MNEMOSYNED_TLS_CERT} \
	-tls.key=${MNEMOSYNED_TLS_KEY}
fi

exec "$@"