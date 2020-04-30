package main

import (
	"fmt"
	"github.com/dzrry/activitycounter/store/postgres"
	"github.com/dzrry/activitycounter/vk/api"
	"log"
	"time"
)

func check(client *api.VKClient) {
	count, members, err := client.GroupGetMembers(groupId, 1000, 0)
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
	/*if err := env.db.InsertGroupMembers(groupId, members); err != nil {
		log.Fatal(err)
	}*/
	for {
		go func() {
			countFromDb, membersFromDb, err := env.db.AllGroupMembers(groupId)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(countFromDb)
			for _, member := range membersFromDb {
				fmt.Println(member.UID)
			}
		}()
		time.Sleep(5 * time.Minute)
	}
}

type Env struct {
	db postgres.Datastore
}

