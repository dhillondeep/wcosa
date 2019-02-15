package log

import (
	"fmt"

	"github.com/fatih/color"
)

type queueBuffer struct {
	message string
	level   Type
	color   *color.Color
}

// Queue is a basic FIFO queue based on slice
type Queue []queueBuffer

// NewQueue returns a new queue with the given initial size.
func NewQueue(size int) *Queue {
	queue := make(Queue, 0, size)
	return &queue
}

// creates a queue buffer from the logging information
func pushLog(a *Args) {
	if a.queue == nil {
		return
	}
	if a != nil {
		a.message = fmt.Sprintf(a.message, a.args...)
	}
	buff := queueBuffer{level: a.level, color: a.color, message: a.message}
	*a.queue = append(*a.queue, buff)
}

// removes the queue buffer and returns the value
func popLog(queue *Queue) queueBuffer {
	if queue == nil {
		return queueBuffer{}
	}
	el := (*queue)[0]
	*queue = (*queue)[1:]
	return el
}

// This provides a queue that can be used to log at different levels
func GetQueue() *Queue {
	return NewQueue(5)
}

// Copy one queue to another
func CopyQueue(from *Queue, to *Queue, spaces Indentation) {
	for len(*from) > 0 {
		value := popLog(from)
		value.message = string(spaces) + value.message
		*to = append(*to, value)
	}
}

// Print Queue on the console with a set indentation
func PrintQueue(queue *Queue, spaces Indentation) {
	for _, value := range *queue {
		value.message = string(spaces) + value.message
		Write(value.level, value.color, value.message)
	}
}
