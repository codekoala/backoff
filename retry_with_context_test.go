// +build go1.7

package backoff

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestRetryWithCanceledContext(t *testing.T) {
	f := func() error {
		t.Error("This function shouldn't be called at all")
		return errors.New("error")
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := RetryNotifyWithContext(ctx, f, NewExponentialBackOff(), nil)
	if err != ctx.Err() {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestRetryWithCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	called := false
	f := func() error {
		if called {
			t.Error("This function shouldn't be called more than once")
		} else {
			cancel()
			called = true
		}
		return errors.New("error")
	}

	err := RetryNotifyWithContext(ctx, f, NewExponentialBackOff(), nil)
	if err != ctx.Err() {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestRetryWithSuccess(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	called := 0
	f := func() (err error) {
		switch called {
		case 0:
			called++
			err = errors.New("error")
		case 1:
			err = nil
		case 2:
			t.Error("This function shouldn't be called more than twice")
		}

		return err
	}

	n := func(err error, delay time.Duration) {
	}

	err := RetryNotifyWithContext(ctx, f, NewExponentialBackOff(), n)
	if err != ctx.Err() {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestRetryWithStop(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	called := 0
	f := func() (err error) {
		called++
		return errors.New("error")
	}

	err := RetryNotifyWithContext(ctx, f, &StopBackOff{}, nil)
	if err == nil {
		t.Errorf("expected error but got nil")
	}

	if called != 1 {
		t.Error("Function should be called once")
	}
}
