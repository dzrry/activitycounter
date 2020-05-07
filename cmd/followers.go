package main

import (
	"fmt"
	"github.com/dzrry/activitycounter/insta"
	"log"
)

func printFollowersByName(inst *insta.Instagram, username string) {
	users, err := inst.Profiles.ByName(username)
	if err != nil {
		log.Fatal("followers 12 " + err.Error())
	}
	followers := users.Followers()
	usrs := 0
	for followers.Next() {
		usrs++
		for _, user := range followers.Users {
			fmt.Print(user.Username + "  ")
		}
	}
	fmt.Println(usrs)
}
