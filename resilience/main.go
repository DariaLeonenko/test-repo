package main

import (
	"errors"
	"fmt"
	"time"
)

func main() {

	fmt.Println("=== Retry ===")
	counter := 0

	err := Retry(func() error {
		counter++
		if counter < 3 {
			return errors.New("fail")
		}
		return nil
	}, 5, 100)

	fmt.Println("Retry result:", err)

	fmt.Println("\n=== Timeout ===")

	err = Timeout(func() error {
		time.Sleep(2 * time.Second)
		return nil
	}, 1000)

	fmt.Println("Timeout result:", err)

	fmt.Println("\n=== DLQ ===")

	messages := []string{"msg1", "msg2", "msg3"}
	dlq := NewDeadLetterQueue()

	ProcessWithDLQ(messages, func(msg string) error {
		if msg == "msg2" {
			return errors.New("error")
		}
		fmt.Println("Processed:", msg)
		return nil
	}, dlq)

	fmt.Println("DLQ:", dlq.GetMessages())
}
