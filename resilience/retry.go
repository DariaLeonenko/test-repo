package main

import (
	"errors"
	"time"
)

func Retry(operation func() error, maxRetries int, baseDelay int) error {
	var err error

	for i := 0; i <= maxRetries; i++ {
		err = operation()
		if err == nil {
			return nil
		}

		delay := time.Duration(baseDelay*(1<<i)) * time.Millisecond
		time.Sleep(delay)
	}

	return errors.New("operation failed after retries: " + err.Error())
}
