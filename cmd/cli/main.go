package main

import (
	"fmt"
	"log"
	"time"

	"github.com/umbe77/ucd/client"
	"github.com/umbe77/ucd/datatypes"
)

func main() {
	c, err := client.NewClient(":3000")
	if err != nil {
		panic(err)
	}

	defer c.Close()

	//PING
	pingres, err := c.Ping()
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println(pingres)

	//A bunch of Set
	setRes, err := c.Set("Key1", datatypes.String, "Value String")
	if err != nil {
		log.Println(err)
	}
	fmt.Printf("Set %s\n", setRes)
	setRes, err = c.Set("Key2", datatypes.Int, 65999)
	if err != nil {
		log.Println(err)
	}
	fmt.Printf("Set %s\n", setRes)
	setRes, err = c.Set("Key3", datatypes.Bool, false)
	if err != nil {
		log.Println(err)
	}
	fmt.Printf("Set %s\n", setRes)
	setRes, err = c.Set("Key3.1", datatypes.Bool, true)
	if err != nil {
		log.Println(err)
	}
	fmt.Printf("Set %s\n", setRes)
	setRes, err = c.Set("Key4", datatypes.Date, time.Now())
	if err != nil {
		log.Println(err)
	}
	fmt.Printf("Set %s\n", setRes)
}
