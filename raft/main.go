package main

import (
	"flag"
	"log"
	"os"
	"strings"

	//"sync"
	ser "raft/server"
)

func main() {

  file, _ := os.ReadFile("others_ip")
  file1, _ := os.ReadFile("my_ip")
  myIp := string(file1[:])
  serversId := strings.Split(string(file[:]), "\n")

  lFlag := flag.Bool("leader", false, "decide if leader or not")
  flag.Parse()

  var isLeader = ser.FOLLOWER

  if (*lFlag == true) {
    isLeader = ser.LEADER
  }

	server1 := ser.NewServer(0, myIp, isLeader, serversId)

  server1.Start()
  
  log.Println("All servers have terminated")
}
