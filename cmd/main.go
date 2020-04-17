package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/dzrry/activitycounter/vk/api"
	"log"
	"net/url"
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

	values := make(url.Values)
	values.Set("group_id", "190873620")
	res, _ := userClient.Do(api.NewRequest("groups.getMembers", userClient.GetToken(), values))
	fmt.Println(res.Response)
	if err != nil {
		log.Fatal(err)
	}

	var members *api.GroupMembers
	json.Unmarshal(res.Response, &members)
	if err != nil {
		log.Fatal("cannot unmarshal json")
	}
	for i := range members.Members {
		fmt.Println(members.Members[i].UID)
	}
}
