package postgres_test

import (
	"flag"
	"os"
	"testing"

	"go.uber.org/zap"

	"time"

	_ "github.com/lib/pq"
	"github.com/piotrkowalczuk/mnemosyne/internal/service/postgres"
)

var (
	testPostgresAddress string
)

func TestMain(m *testing.M) {
	flag.StringVar(&testPostgresAddress, "postgres.address", getStringEnvOr("MNEMOSYNED_POSTGRES_ADDRESS", "postgres://localhost/test?sslmode=disable"), "")
	flag.Parse()

	os.Exit(m.Run())
}

func getStringEnvOr(env, or string) string {
	if v := os.Getenv(env); v != "" {
		return v
	}
	return or
}

func TestInit_retry(t *testing.T) {
	_, err := postgres.Init(testPostgresAddress, postgres.Opts{
		Logger:  zap.L(),
		Timeout: 2 * time.Second,
		Retry:   600 * time.Nanosecond,
	})
	if err != nil && err != postgres.Timeout {
		t.Fatalf("unexpected error: %s", err.Error())
	}
}

func TestInit_timeout(t *testing.T) {
	_, err := postgres.Init(testPostgresAddress, postgres.Opts{
		Logger:  zap.L(),
		Timeout: 1 * time.Nanosecond,
		Retry:   1 * time.Nanosecond,
	})
	if err != nil && err != postgres.Timeout {
		t.Fatalf("unexpected error: %s", err.Error())
	}
	if err == nil {
		t.Fatal("error expected, got nil")
	}
}
func TestInit_empty(t *testing.T) {
	_, err := postgres.Init(testPostgresAddress, postgres.Opts{
		Logger: zap.L(),
	})
	if err != nil {
		t.Fatalf("unexpected error: %s", err.Error())
	}
}
