package logger_test

import (
	"testing"

	"github.com/piotrkowalczuk/mnemosyne/internal/service/logger"
)

func TestInit(t *testing.T) {
	_, err := logger.Init(logger.Opts{
		Environment: "production",
	})
	if err != nil {
		t.Fatalf("unexpected error: %s", err.Error())
	}
	_, err = logger.Init(logger.Opts{
		Environment: "development",
	})
	if err != nil {
		t.Fatalf("unexpected error: %s", err.Error())
	}
	_, err = logger.Init(logger.Opts{})
	if err != nil {
		t.Fatalf("unexpected error: %s", err.Error())
	}
}
