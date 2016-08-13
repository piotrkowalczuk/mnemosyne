#!/bin/sh
set -e

: ${MNEMOSYNED_PORT:=8080}
: ${MNEMOSYNED_HOST:=0.0.0.0}
: ${MNEMOSYNED_LOG_ADAPTER:=stdout}
: ${MNEMOSYNED_LOG_FORMAT:=json}
: ${MNEMOSYNED_LOG_LEVEL:=6}
: ${MNEMOSYNED_STORAGE:=postgres}
: ${MNEMOSYNED_MONITORING:=false}
: ${MNEMOSYNED_POSTGRES_ADDRESS:=postgres://postgres:postgres@postgres/postgres?sslmode=disable}
: ${MNEMOSYNED_POSTGRES_DEBUG:=false}
: ${MNEMOSYNED_TLS_ENABLED:=false}

if [ "$1" = 'mnemosyned' ]; then
exec mnemosyned -host=$MNEMOSYNED_HOST \
	-port=$MNEMOSYNED_PORT \
	-storage=$MNEMOSYNED_STORAGE \
	-log.adapter=$MNEMOSYNED_LOG_ADAPTER \
	-log.format=$MNEMOSYNED_LOG_FORMAT \
	-log.level=$MNEMOSYNED_LOG_LEVEL \
	-monitoring=$MNEMOSYNED_MONITORING \
	-postgres.address=$MNEMOSYNED_POSTGRES_ADDRESS \
	-tls=$MNEMOSYNED_TLS_ENABLED \
	-tls.cert=$MNEMOSYNED_TLS_CERT \
	-tls.key=$MNEMOSYNED_TLS_KEY
fi

exec "$@"