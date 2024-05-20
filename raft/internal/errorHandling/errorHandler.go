package errorHandling

import "sync"

type errTypo = int

/* type of error */
const (
  PANIC errTypo = iota
  WARNING
)

/* new custom error type containing the error and its type */
type errno struct {
  ty errTypo
  err error
}

type ErrorHandler interface {
  Send(err ...errno)
  Start()
}

func NewErrorHandler() ErrorHandler {
  var errChan chan errno = make(chan errno)
  
  return &errorHandlerImp{
    &errChan,
    &sync.WaitGroup{},
  }
}

func NewError(ty errTypo, err error) errno {
  return errno{
    ty: ty,
    err: err,
  }
}


