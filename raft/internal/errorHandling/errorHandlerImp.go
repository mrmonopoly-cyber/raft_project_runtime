package errorHandling

import (
	"log"
	"sync"
)

const colorRed = "\033[0;31m"
const colorGreen = "\033[0;32m"

type errorHandlerImp struct {
  errorChannel *chan errno
  wg *sync.WaitGroup
}

func (this *errorHandlerImp) Send(err ...errno) {
  for _,e := range err {
    *this.errorChannel <- e
  }
}

func (this *errorHandlerImp) Start() {
  this.wg.Add(1)
  defer this.wg.Done()

  for {
    select {
      case err := <- *this.errorChannel:
        // take care of the error
        switch err.ty {
          case PANIC:
            log.Printf("%s Panic: %s", colorRed, err.err)
            // graceful restart of the node
          case WARNING:
            log.Printf("%s Warning: %s", colorGreen, err.err)
        }  
    }
  }
}

