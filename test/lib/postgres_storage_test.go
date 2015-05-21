// +build postgres

package lib

import (
	"log"
	"os"
	"testing"

	"runtime"

	"github.com/go-soa/mnemosyne/lib"
	"github.com/go-soa/mnemosyne/service"
)

var (
	ps *lib.PostgresStorage
)

func TestMain(m *testing.M) {
	runtime.GOMAXPROCS(runtime.NumCPU())

	service.InitConfig("../../conf", "testing")

	l := log.New(os.Stderr, "", log.LstdFlags)
	ssc := service.Config.SessionStorage.Postgres

	service.InitPostgres(ssc, l)
	ps = lib.NewPostgresStorage(service.Postgres, ssc.TableName)

	code := m.Run()

	service.Postgres.Close()

	os.Exit(code)
}

func TestPostgresStorageNew(t *testing.T) {
	testStorageNew(t, ps)
}

func TestPostgresStorageGet(t *testing.T) {
	testStorageGet(t, ps)
}

func TestPostgresStorageExists(t *testing.T) {
	testStorageExists(t, ps)
}

func TestPostgresStorageAbandon(t *testing.T) {
	testStorageAbandon(t, ps)
}

func TestPostgresStorageSetData(t *testing.T) {
	testStorageSetData(t, ps)
}

func TestPostgresStorageCleanup(t *testing.T) {
	testStorageCleanup(t, ps)
}
