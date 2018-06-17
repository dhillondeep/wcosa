package log

import (
	"fmt"
	"github.com/fatih/color"
)

type Node struct {
	Value stackBuffer
}

// NewQueue returns a new queue with the given initial size.
func NewQueue(size int) *Queue {
	return &Queue{
		nodes: make([]*Node, size),
		size:  size,
	}
}

// Queue is a basic FIFO queue based on a circular list that resizes as needed.
type Queue struct {
	nodes []*Node
	size  int
	head  int
	tail  int
	count int
}

// Push adds a node to the queue.
func (q *Queue) Push(n *Node) {
	if q.head == q.tail && q.count > 0 {
		nodes := make([]*Node, len(q.nodes)+q.size)
		copy(nodes, q.nodes[q.head:])
		copy(nodes[len(q.nodes)-q.head:], q.nodes[:q.head])
		q.head = 0
		q.tail = len(q.nodes)
		q.nodes = nodes
	}
	q.nodes[q.tail] = n
	q.tail = (q.tail + 1) % len(q.nodes)
	q.count++
}

// Pop removes and returns a node from the queue in first to last order.
func (q *Queue) Pop() *Node {
	if q.count == 0 {
		return nil
	}
	node := q.nodes[q.head]
	q.head = (q.head + 1) % len(q.nodes)
	q.count--
	return node
}

type stackBuffer struct {
	text          string
	logType       string
	providedColor *color.Color
}

func pushLog(queue *Queue, logType string, providedColor *color.Color, message string, a ...interface{}) {
	text := message

	if a != nil {
		text = fmt.Sprintf(message, a...)
	}

	buff := stackBuffer{logType: logType, providedColor: providedColor, text: text}
	queue.Push(&Node{buff})
}

func popLog(queue *Queue) stackBuffer {
	return queue.Pop().Value
}
