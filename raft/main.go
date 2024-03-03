package main

import (
	"log"
	//"os"
	ser "raft/server"
)

func main() {

	log.Println("Ciaoo")

	serversId := []int{8080, 8090, 9000, 5060, 3029}

	server1 := ser.NewServer(0, 8080, ser.LEADER, serversId)
	server2 := ser.NewServer(0, 8090, ser.FOLLOWER, serversId)
	server3 := ser.NewServer(0, 9000, ser.FOLLOWER, serversId)
	server4 := ser.NewServer(0, 5060, ser.FOLLOWER, serversId)
	server5 := ser.NewServer(0, 3029, ser.FOLLOWER, serversId)

	go server1.Start()
	go server2.Start()
	go server3.Start()
	go server4.Start()
	server5.Start()

	log.Println("nn")
}
