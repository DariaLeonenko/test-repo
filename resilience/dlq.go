package main

type DeadLetterQueue struct {
	messages []string
}

func NewDeadLetterQueue() *DeadLetterQueue {
	return &DeadLetterQueue{
		messages: []string{},
	}
}

func (d *DeadLetterQueue) Add(msg string) {
	d.messages = append(d.messages, msg)
}

func (d *DeadLetterQueue) GetMessages() []string {
	return d.messages
}

func ProcessWithDLQ(messages []string, handler func(string) error, dlq *DeadLetterQueue) {
	for _, msg := range messages {
		err := handler(msg)
		if err != nil {
			dlq.Add(msg)
		}
	}
}
