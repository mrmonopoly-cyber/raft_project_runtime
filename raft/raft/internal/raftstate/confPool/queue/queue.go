package queue

type Queue[T any] interface{
    Pop() T
    Push(v T)
    Size() int
    WaitEl() <- chan int
}

func NewQueue[T any]() Queue[T]{
    return NewQueueImp[T]() 
}
