package backoff

import (
	"errors"
	"log"
	"testing"
	"time"
)

func TestTicker(t *testing.T) {
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

	b := NewExponentialBackOff()
	ticker := NewTicker(b)

	var err error
	for _ = range ticker.C {
		if err = f(); err != nil {
			t.Log(err)
			continue
		}

		break
	}
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}
	if i != successOn {
		t.Errorf("invalid number of retries: %d", i)
	}
}

func TestTickerStop(t *testing.T) {
	const (
		stopOn    = 2
		successOn = 3
	)

	var (
		i   = 0
		err error
	)

	b := NewConstantBackOff(time.Second)
	ticker := NewTicker(b)

	// This function is successful on "successOn" calls.
	f := func() error {
		i++
		log.Printf("f called %d time(s)", i)

		switch i {
		case stopOn:
			ticker.Stop()
			return nil
		case successOn:
			log.Println("OK")
			return nil
		default:
			log.Println("error")
		}

		return errors.New("error")
	}

	for _ = range ticker.C {
		if err = f(); err != nil {
			t.Log(err)
			continue
		}

		break
	}

	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}
	if i != stopOn {
		t.Errorf("invalid number of retries: %d", i)
	}
}

func TestTickerStopBackoff(t *testing.T) {
	const (
		stopOn    = 2
		successOn = 3
	)

	var (
		i   = 0
		err error
	)

	b := &StopBackOff{}
	ticker := NewTicker(b)

	// This function is successful on "successOn" calls.
	f := func() error {
		i++
		log.Printf("f called %d time(s)", i)

		switch i {
		case stopOn:
			ticker.Stop()
			return nil
		case successOn:
			log.Println("OK")
			return nil
		default:
			log.Println("error")
		}

		return errors.New("error")
	}

	for _ = range ticker.C {
		if err = f(); err != nil {
			t.Log(err)
			continue
		}

		break
	}

	if err == nil {
		t.Errorf("expected error but got none")
	}
	if i != 1 {
		t.Errorf("invalid number of retries: %d", i)
	}
}
