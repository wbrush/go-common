package helpers

import (
	"context"
	"time"
)

type (
	Module interface {
		Run()
		Title() string
		GracefulStop(context.Context) error
	}

	Daemon interface {
		Module
		SetInterval(interval time.Duration)
		GetInterval() time.Duration
	}
)
