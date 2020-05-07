package main

import (
	"fmt"
	"github.com/dzrry/activitycounter/insta"
	"log"
)

func likesCount(inst *insta.Instagram) {
	user, err := inst.Profiles.ByName("nina.shahova")
	if err != nil {
		log.Fatal("likes 12 " + err.Error())
	}
	feed := user.Feed()
	for feed.Next() {
		for _, post := range feed.Items {
			fmt.Printf("Likes: %d on post: %s\n", post.Likes, post.ID)
		}
	}
}

func likersNames(inst *insta.Instagram) {
	user, err := inst.Profiles.ByName("nina.shahova")
	if err != nil {
		log.Fatal("likes 25" + err.Error())
	}
	user, err = inst.Profiles.ByID(user.ID)
	if err != nil {
		log.Fatal("likes 29" + err.Error())
	}

	feed := user.Feed()
	for feed.Next() {
		for _, post := range feed.Items {
			post.SyncLikers()
			fmt.Printf("\n\nLikers for post %s", post.ID)
			for _, liker := range post.Likers {
				fmt.Print(liker.Username + " ")
			}
		}
	}
}
