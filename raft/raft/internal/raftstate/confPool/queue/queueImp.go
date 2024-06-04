package queue

import (
    "log"
)

type QueueImp[T any] struct {
    buffer []T
    notifyAdd chan int
    C   <- chan int
}

// Pop implements Queue.
func (q *QueueImp[T]) Pop() T {
    var res = q.buffer[0]
    q.buffer = q.buffer[1:]
    return res
}

// Push implements Queue.
func (q *QueueImp[T]) Push(v T) {
    q.buffer = append(q.buffer, v)
    log.Println("pushing notification on: ", q.notifyAdd)
    q.notifyAdd <- 1
}

// Size implements Queue.
func (q *QueueImp[T]) Size() int {
    return len(q.buffer)
}

func NewQueueImp[T any]() *QueueImp[T] {
	var res = &QueueImp[T]{
        buffer: nil,
        notifyAdd: make(chan int),
    }
    res.C = res.notifyAdd
    return res

}
