#!/bin/sh
set -e

: ${MNEMOSYNED_PORT:=8080}
curl -f http://localhost:$((MNEMOSYNED_PORT+1))/health || exit 1