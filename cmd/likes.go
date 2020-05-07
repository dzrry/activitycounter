package main

import (
	"fmt"
	"github.com/dzrry/activitycounter/insta"
	"log"
)

func likesCount(inst *insta.Instagram) {
	users, err := inst.Profiles.ByName("nina.shahova")
	if err != nil {
		log.Fatal("likes 12 " + err.Error())
	}
	feed := users.Feed()
	for feed.Next() {
		for _, post := range feed.Items {
			fmt.Printf("Likes: %d on post: %s\n", post.Likes, post.ID)
		}
	}
}
