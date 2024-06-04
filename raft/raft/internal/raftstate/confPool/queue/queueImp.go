package queue

import (
    "log"
)

type queueImp[T any] struct {
    buffer []T
    notifyAdd chan int
    res <-chan int
}

// Pop implements Queue.
func (q *queueImp[T]) Pop() T {
    var res = q.buffer[0]
    q.buffer = q.buffer[1:]
    return res
}

// Push implements Queue.
func (q *queueImp[T]) Push(v T) {
    q.buffer = append(q.buffer, v)
    go func(){

        log.Println("pushing notification on: ", q.notifyAdd)
        q.notifyAdd <- 1
    }()
}

// Size implements Queue.
func (q *queueImp[T]) Size() int {
    return len(q.buffer)
}

// WaitEl implements Queue.
func (q *queueImp[T]) WaitEl() <-chan int {
    return q.res
}

func NewQueueImp[T any]() Queue[T] {
	var res = &queueImp[T]{
        buffer: nil,
        notifyAdd: make(chan int),
        res: nil,
    }
    res.res = res.notifyAdd
    return res
}
