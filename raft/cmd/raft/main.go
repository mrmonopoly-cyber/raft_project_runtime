package main

import (
	//"flag"
	"log"
	"os"
	"strings"

	//"strings"

	"raft/internal/raftstate"
	ser "raft/internal/server"
	"sync"
)

func main() {

  /*file, _ := os.ReadFile("others_ip")
  fileMyIp, _ := os.ReadFile("my_ip")
  stringMyIp := string(fileMyIp[:])
  serversId := strings.Split(string(file[:]), "\n")*/

  /*lFlag := flag.Bool("leader", false, "decide if leader or not")
  flag.Parse()

  var isLeader = ser.FOLLOWER

  if (*lFlag == true) {
    isLeader = ser.LEADER
  }*/

  var fileMyIp, _ = os.ReadFile("my_ip")
  var fileOthersIp, _ = os.ReadFile("others_ip")
  var stringMyIp = string(fileMyIp)
  var stringOthersIp = string(fileOthersIp)
  var addresses []string = strings.Split(stringOthersIp, "\n")

  var server1 *ser.Server = ser.NewServer(0,stringMyIp, "8080", raftstate.FOLLOWER, addresses)
  // server2 := ser.NewServer(0,"localhost", "8090", raftstate.FOLLOWER, []string{"8080", "8090"})

  var wg sync.WaitGroup

  wg.Add(1)

  go func() {
    defer wg.Done()
    server1.Start()
  } ()

  // go func() {
  //   defer wg.Done()
  //   server2.Start()
  // } ()
  
  wg.Wait()
  
  log.Println("All servers have terminated")
}
