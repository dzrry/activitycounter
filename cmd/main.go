package main

import (
	"flag"
	"github.com/dzrry/activitycounter/insta"
	"github.com/dzrry/activitycounter/vk/api"
	"log"
)

const groupId = 190873620
const rbkId = 25232578

func main() {
	inst := createInsta()
	if err := inst.Login(); err != nil {
		log.Fatal("main 17 " + err.Error())
	}
	defer inst.Logout()

	likersNames(inst)
}

func createUserVk() (client *api.VKClient) {
	login := flag.String("login", "", "login string")
	password := flag.String("password", "", "password string")
	flag.Parse()

	client, err := api.NewVKClient(3, *login, *password)
	if err != nil {
		log.Fatal(err)
	}
	return
}

func createInsta() (inst *insta.Instagram) {
	login := flag.String("login", "", "login string")
	password := flag.String("password", "", "password string")
	flag.Parse()

	inst = insta.New(*login, *password)
	return
}
