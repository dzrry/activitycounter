package main

import (
	"flag"
	"fmt"
	"github.com/dzrry/activitycounter/vk/api"
	"log"
)

const groupId = 190873620

func main() {
	login := flag.String("login", "", "login string")
	password := flag.String("password", "", "password string")
	//scope := int64(140488159)
	flag.Parse()

	userClient, err := api.NewVK(2, *login, *password)
	if err != nil {
		log.Fatal(err)
	}

	count, members, err := userClient.GroupGetMembers(groupId)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(count)
	for i := range members {
		fmt.Println(&members[i].UID)
	}
}
