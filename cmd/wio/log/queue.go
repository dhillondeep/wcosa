package log

import (
    "fmt"
    "github.com/fatih/color"
)

type queueBuffer struct {
    text          string
    logType       Type
    providedColor *color.Color
}

// Queue is a basic FIFO queue based on slice
type Queue []queueBuffer

// NewQueue returns a new queue with the given initial size.
func NewQueue(size int) *Queue {
    queue := make(Queue, 0, size)
    return &queue
}

// creates a queue buffer from the logging information
func pushLog(queue *Queue, logType Type, providedColor *color.Color, message string, a ...interface{}) {
    if queue == nil {
        return
    }
    text := message

    if a != nil {
        text = fmt.Sprintf(message, a...)
    }

    buff := queueBuffer{logType: logType, providedColor: providedColor, text: text}
    *queue = append(*queue, buff)
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
