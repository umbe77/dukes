// Copyright (c) 2023 Robeto Ughi
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/umbe77/dukes/server"
)

func main() {
	fmt.Println("Hello from dukes!")

	//Initializing tcp Server
	srv := server.NewServer(":3000", "raft", ":3001")

	//Start tcp Server
	if err := srv.Run(); err != nil {
		log.Fatalf("Failed to start server at %s: %s", srv.ListenerAddr, err.Error())
	}

	terminate := make(chan os.Signal, 1)
	signal.Notify(terminate, os.Interrupt, syscall.SIGTERM)
	<-terminate
	log.Println("dukes exited")
}
