// +build postgres
package main

import (
	"os"
	"testing"

	"github.com/piotrkowalczuk/sklog"
)

func TestMain(m *testing.M) {
	config.Parse()

	configPostgres := config.storage.postgres
	configLogger := config.logger

	initLogger(configLogger.adapter, configLogger.format, configLogger.level, sklog.KeySubsystem, "mnemosyne")
	initPostgres(configPostgres.connectionString, configPostgres.retry, logger)
	initStorage(initPostgresStorage(configPostgres.tableName, postgres), logger)

	code := m.Run()

	storage.TearDown()
	postgres.Close()

	os.Exit(code)
}

func TestPostgresStorage_Create(t *testing.T) {
	testStorage_Create(t, storage)
}

func TestPostgresStorage_Get(t *testing.T) {
	testStorage_Get(t, storage)
}

func TestPostgresStorage_List(t *testing.T) {
	testStorage_List(t, storage)
}

func TestPostgresStorage_Exists(t *testing.T) {
	testStorage_Exists(t, storage)
}

func TestPostgresStorage_Abandon(t *testing.T) {
	testStorage_Abandon(t, storage)
}

func TestPostgresStorage_SetData(t *testing.T) {
	testStorage_SetData(t, storage)
}

func TestPostgresStorage_Delete(t *testing.T) {
	testStorage_Delete(t, storage)
}
