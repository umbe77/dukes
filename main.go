package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"path"
	"syscall"

	"github.com/umbe77/dukes/client"
	"github.com/umbe77/dukes/server"
)

func main() {
	fmt.Println("Hello from dukes!")

	serverAddress := ":3000"
	nodeId := "Node0"
	raftAddr := "localhost:3001"
	joinAddr := ""

	if val, ok := os.LookupEnv("DUKES_SRV_ADDR"); ok {
		serverAddress = val
	}

	if val, ok := os.LookupEnv("DUKES_NODEID"); ok {
		nodeId = val
	}

	raftDir := path.Join("raft", nodeId)

	if val, ok := os.LookupEnv("DUKES_RAFT_BIND"); ok {
		raftAddr = val
	}

	if val, ok := os.LookupEnv("DUKES_JOIN_ADDR"); ok {
		joinAddr = val
	}

	//Initializing Server
	srv := server.NewServer(serverAddress, raftDir, raftAddr)

	//Start tcp Server and raft Cluster
	if err := srv.Run(nodeId, joinAddr == ""); err != nil {
		log.Fatalf("Failed to start server at %s: %s", srv.ListenerAddr, err.Error())
	}

	if joinAddr != "" {
		joinClient, err := client.NewClient(joinAddr)
		if err != nil {
			log.Fatalf("Cannot connect to %s: %s", joinAddr, err)
		}
		defer joinClient.Close()
		if err := joinClient.Join(nodeId, raftAddr); err != nil {
			log.Fatalf("Cannot join to %s: %s", joinAddr, err)
		}
	}

	terminate := make(chan os.Signal, 1)
	signal.Notify(terminate, os.Interrupt, syscall.SIGTERM)
	<-terminate
	log.Println("dukes exited")
}
