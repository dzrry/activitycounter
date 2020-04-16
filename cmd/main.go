package main

import (
	"flag"
	"fmt"
	"github.com/dzrry/activitycounter/vk/api"
	"log"
)

func main() {
	login := flag.String("login", "", "login string")
	password := flag.String("password", "", "password string")
	scope := int64(140488159)
	flag.Parse()

	userClient, err := api.NewClientFromLogin(*login, *password, scope)
	if err != nil {
		log.Fatal("no user client")
	}
	fmt.Println(userClient)

}
