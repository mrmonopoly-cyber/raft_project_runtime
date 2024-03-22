package main

import (
	// 	//"flag"
	"bufio"
	"log"
	"net"
	"os"
	"sync"

	// 	"strings"
	//
	"strings"
	// ser "raft/internal/server"
	// "sync"
)

func main(){
    var wg sync.WaitGroup
    wg.Add(1)
    var workDir = "/root/mount/raft/"
    var fileOthersIp, erro = os.ReadFile(workDir + "others_ip")
    if erro != nil {
        log.Printf("could not find other ips")
    }
    var stringOthersIp = string(fileOthersIp)
    var addresses []string = strings.Split(stringOthersIp, "\n")
    log.Printf("creating listener")
    var list net.Listener
    var err error
    list,err = net.Listen("tcp",":8080")
    if err !=nil {
        panic("failed to create listener")
    }

    var numAddrs = len(addresses)-1
    log.Printf("addresses found: %v", numAddrs)
    if numAddrs == 0 {
        log.Printf("accepting connection")
        var con net.Conn
        con, err = list.Accept()
        if err != nil {
            panic("accept failed")
        }
        var mex string

        go func () {
            for {
                mex = "to you"
                _,err = con.Write([]byte(mex + "\n"))
            }
        }()

        log.Printf("recr mex")
        go func (){
            for{
                mex, err = bufio.NewReader(con).ReadString('\n')
                if err != nil{
                    log.Panicf("error in receiving the mex: %v", err)
                }
                log.Printf("message received %v", mex)
            }
        }()
        wg.Wait()

        con.Close()

    }else{
        var con net.Conn
        var err error
        log.Printf("enstablish connection")
        con,err = net.Dial("tcp",addresses[0] + ":8080")
        if err != nil {
            panic("error enstablish connection")
        }

        log.Printf("sending message")
        go func ()  {
            for{
                var mex = "hello"
                _,err = con.Write([]byte(mex + "\n"))
                if err != nil {
                    panic("error sending data")
                }
            }
        }()

        var mex string
        log.Println("data sent")


        go func (){
            for {
                mex, err = bufio.NewReader(con).ReadString('\n')
                if err != nil{
                    log.Panicf("error in receiving the mex: %v", err)
                }
                log.Printf("data received: %v", mex)
            }
        }()

        wg.Wait()
        con.Close()
    }


    log.Println("listening on port: " + "8080")
}

// func main() {
//
//   var workDir = "/root/mount/raft/"
//   var fileMyIp, errm = os.ReadFile(workDir + "my_ip")
//   if errm != nil {
//     panic("could not find my ip")
//   }
//   var fileOthersIp, erro = os.ReadFile(workDir + "others_ip")
//   if erro != nil {
//     panic("could not find other ips")
//   }
//   var stringMyIp = strings.Split(string(fileMyIp), "\n")
//   var stringOthersIp = string(fileOthersIp)
//   var addresses []string = strings.Split(stringOthersIp, "\n")
//
//
//   log.Println("listening on port: " + "8080")
//
//   var server1 *ser.Server = ser.NewServer(0, stringMyIp[0], "8080", addresses)
//
//   var wg sync.WaitGroup
//
//
//   wg.Add(1)
//   go func() {
//     defer wg.Done()
//     server1.Start()
//   } ()
//   
//   wg.Wait()
//   
//   log.Println("All servers have terminated")
// }
