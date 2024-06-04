package queue

type Queue[T any] interface{
    Push(v T)
    Pop() T
    Size() int
    WaitEl() <- chan int
}

func NewQueue[T any]() Queue[T]{
    return newQueueImp[T]()
}
