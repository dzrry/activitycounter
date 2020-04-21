package main

import (
	"flag"
	"fmt"
	"github.com/dzrry/activitycounter/store/postgres"
	"github.com/himidori/golang-vk-api"
	"log"
)

const groupId = 190873620

func main() {
	login := flag.String("login", "", "login string")
	password := flag.String("password", "", "password string")
	//scope := int64(140488159)
	flag.Parse()

	client, err := vkapi.NewVKClient(3, *login, *password)
	if err != nil {
		log.Fatal(err)
	}

	count, members, err := client.GroupGetMembers(groupId, 1000)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(count)
	for i := range members {
		fmt.Println(members[i].UID)
	}

	db, err := postgres.NewDB("user=postgres password=Val040674 dbname=activitycounter sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	env := &Env{db}
	id, err := env.db.InsertGroupMembers(groupId, count, members)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(id)

}

type Env struct {
	db postgres.Datastore
}
