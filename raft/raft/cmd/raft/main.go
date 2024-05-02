package main

import (
	"log"
	"os"
	"sync"
	"strings"
	ser "raft/internal/server"
)

func main() {

  var workDir = "/root/mount/raft/"
  var fsRootDir = workDir + "locafs/"
  var fileMyIp, errm = os.ReadFile(workDir + "my_ip")
  if errm != nil {
    panic("could not find my ip")
  }
  var fileOthersIp, erro = os.ReadFile(workDir + "others_ip")
  if erro != nil {
    panic("could not find other ips")
  }
  var stringMyIp []string= strings.Split(string(fileMyIp), "\n")
  var stringOthersIp = string(fileOthersIp)
  var addresses []string = strings.Split(stringOthersIp, "\n")
  for i, v := range addresses {
      log.Println("before trim: ",v)
      addresses[i] = strings.Trim(v,"\n")
      log.Println("after trim: ",addresses[i])
  }

  os.Mkdir(fsRootDir, os.ModePerm)

  log.Println("listening on port: " + "8080" + " with ip: " + stringMyIp[0])

  var server1 ser.Server = ser.NewServer(0, stringMyIp[0], stringMyIp[1], "8080", addresses, fsRootDir)

  var wg sync.WaitGroup


  wg.Add(1)
  go func() {
    defer wg.Done()
    server1.Start()
  } ()
  
  wg.Wait()
  
  log.Println("All servers have terminated")
}
