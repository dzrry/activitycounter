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
	for _, member := range members {
		fmt.Println(member.UID)
	}

	db, err := postgres.NewDB("user=postgres password=postgres dbname=activitycounter sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	env := &Env{db}
	if err := env.db.InsertGroupMembers(groupId, members); err != nil {
		log.Fatal(err)
	}
	/*countFromDb, membersFromDb, err := env.db.AllGroupMembers(groupId)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(countFromDb)
	for _, member := range membersFromDb {
		fmt.Println(member.UID)
	}*/
}

type Env struct {
	db postgres.Datastore
}
