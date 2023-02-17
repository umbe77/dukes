// Copyright (c) 2023 Robeto Ughi
// 
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package main

import (
	"fmt"
	"log"
	"time"

	"github.com/umbe77/dukes/client"
	"github.com/umbe77/dukes/datatypes"
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
	setRes, err = c.Set("Key2", datatypes.Int, int32(65999))
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

	getRes, err := c.Get("Key2")
	if err != nil {
		log.Println(err)
	}
	fmt.Printf("Get %v\n", getRes.ToAny())

	hasRes, err := c.Has("key3")
	if err != nil {
		log.Println(err)
	}
	fmt.Printf("Has key3 %v\n", hasRes)

	hasRes, err = c.Has("Key4")
	if err != nil {
		log.Println(err)
	}
	fmt.Printf("Has Key4 %v\n", hasRes)

	hasRes, err = c.Has("Key3.1")
	if err != nil {
		log.Println(err)
	}

	fmt.Printf("Has Key3.1 %v\n", hasRes)
	err = c.Del("Key3.1")
	if err != nil {
		log.Println(err)
	}
	fmt.Println("Key3.1 Deleted")

	hasRes, err = c.Has("Key3.1")
	if err != nil {
		log.Println(err)
	}
	fmt.Printf("Has Key3.1 %v\n", hasRes)

	fmt.Println("-------------------")
	for k := range c.Dump() {
		fmt.Printf("%s\n", k)
	}
}
