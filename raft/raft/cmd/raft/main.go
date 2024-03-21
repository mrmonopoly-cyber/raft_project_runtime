package main

import (
	//"flag"
	"log"
	"os"
	"strings"

	//"strings"

	ser "raft/internal/server"
	"sync"
)

func main() {

  var workDir = "/root/mount/raft/"
  var fileMyIp, errm = os.ReadFile(workDir + "my_ip")
  if errm != nil {
    panic("could not find my ip")
  }
  var fileOthersIp, erro = os.ReadFile(workDir + "others_ip")
  if erro != nil {
    panic("could not find other ips")
  }
  var stringMyIp = string(fileMyIp)
  var stringOthersIp = string(fileOthersIp)
  var addresses []string = strings.Split(stringOthersIp, "\n")

  log.Println("list of other nodes: ", stringOthersIp)

  var server1 *ser.Server = ser.NewServer(0,stringMyIp, "8080", addresses)

  var wg sync.WaitGroup


  wg.Add(1)
  go func() {
    defer wg.Done()
    server1.Start()
  } ()
  
  wg.Wait()
  
  log.Println("All servers have terminated")
}
