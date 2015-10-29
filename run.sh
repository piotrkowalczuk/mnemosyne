#!/usr/bin/env bash

export MNEMOSYNE_LOGGER_FORMAT="humane"
export MNEMOSYNE_LOGGER_ADAPTER="stdout"
export MNEMOSYNE_LOGGER_LEVEL=6
export MNEMOSYNE_STORAGE_ENGINE="postgres"
export MNEMOSYNE_STORAGE_POSTGRES_CONNECTION_STRING="postgres://localhost/soa?sslmode=disable"
export MNEMOSYNE_STORAGE_POSTGRES_TABLE_NAME="mnemosyne_session"
export MNEMOSYNE_STORAGE_POSTGRES_RETRY=10
export MNEMOSYNE_HOST=localhost
export MNEMOSYNE_PORT=9001
export MNEMOSYNE_SUBSYSTEM="mnemosyne"
make "$@"
