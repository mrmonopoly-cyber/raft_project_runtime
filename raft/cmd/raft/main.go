package main

import (
	//"flag"
	"log"
	//"os"
	//"strings"

	"sync"
	ser "raft/internal/server"
)

func main() {

  /*file, _ := os.ReadFile("others_ip")
  file1, _ := os.ReadFile("my_ip")
  myIp := string(file1[:])
  serversId := strings.Split(string(file[:]), "\n")*/

  /*lFlag := flag.Bool("leader", false, "decide if leader or not")
  flag.Parse()

  var isLeader = ser.FOLLOWER

  if (*lFlag == true) {
    isLeader = ser.LEADER
  }*/

  server1 := ser.NewServer(0, "8080", ser.LEADER, []string{"8080", "8090"})
  server2 := ser.NewServer(0, "8090", ser.FOLLOWER, []string{"8080", "8090"})

  var wg sync.WaitGroup

  wg.Add(2)

  go func() {
    defer wg.Done()
    server1.Start()
  } ()

  go func() {
    defer wg.Done()
    server2.Start()
  } ()
  
  wg.Wait()
  
  log.Println("All servers have terminated")
}
