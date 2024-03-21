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
  var fileMyIp, _ = os.ReadFile(workDir + "my_ip")
  var fileOthersIp, _ = os.ReadFile(workDir + "others_ip")
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
