package backoff

import (
	"errors"
	"log"
	"testing"
	"time"
)

func TestRetry(t *testing.T) {
	const successOn = 3
	var i = 0

	// This function is successful on "successOn" calls.
	f := func() error {
		i++
		log.Printf("f called %d time(s)", i)

		if i == successOn {
			log.Println("OK")
			return nil
		}

		log.Println("error")
		return errors.New("error")
	}

	err := Retry(f, NewExponentialBackOff())
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}
	if i != successOn {
		t.Errorf("invalid number of retries: %d", i)
	}
}

func TestRetryNotifyWithStop(t *testing.T) {
	called := 0
	f := func() (err error) {
		called++
		return errors.New("error")
	}

	err := RetryNotify(f, &StopBackOff{}, nil)
	if err == nil {
		t.Errorf("expected error but got nil")
	}

	if called != 1 {
		t.Error("f should be called once")
	}
}

func TestRetryNotifyWithNotifier(t *testing.T) {
	called := 0
	notified := 0

	f := func() (err error) {
		if called < 1 {
			err = errors.New("error")
		} else {
			err = nil
		}

		called++
		return err
	}

	n := func(err error, delay time.Duration) {
		notified++
	}

	err := RetryNotify(f, &ZeroBackOff{}, n)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	if called != 2 {
		t.Errorf("f should be called twice; called: %d", called)
	}

	if notified != 1 {
		t.Error("should be notified once")
	}
}
