package main

import (
	"fmt"
	"log"

	"github.com/umbe77/ucd/cache"
	"github.com/umbe77/ucd/server"
)

func main() {
	fmt.Println("Hello from cache!")
	cache := cache.NewCache()
	srv := server.NewServer(":3000", cache)
	log.Fatal(srv.Run())
}
