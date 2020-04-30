package main

import (
	"flag"
	"fmt"
	"github.com/dzrry/activitycounter/vk/api"
	"log"
	"strings"
)

const groupId = 190873620
const rbkId = 25232578

func main() {
	client := createUserVk()

	tmpl := `
var groupId = %d;
var items = [%s];
var countByRequest = %d;
var hasNext = true;
var offset = %d;
var calls = 0;

while ((items.length > 0 || hasNext) && calls < 25) {
	if (items.length == 0 && hasNext) {
		var members = API.groups.getMembers({"group_id":groupId,"count":countByRequest,"offset":offset});

		calls = calls + 1;
		offset = offset + 1;
		if (offset >= members.count) {
			hasNext = false;
		}
		items = members.items;
	}
}
return {"items":items,"has_next":hasNext,"offset":offset};`

	offset := 0
	var items []string

	for {
		code := fmt.Sprintf(tmpl, rbkId, strings.Join(items, ","), 1000, offset)
		resp, err := client.Execute(code)
		if err != nil {
			log.Fatal(err)
		}

		offset = int(resp["offset"].(float64))
		respItems := resp["items"].([]interface{})
		if len(respItems) == 0 {
			break
		}

		items = make([]string, 0, len(respItems))
		for _, item := range respItems {
			fmt.Println(item)
			items = append(items, fmt.Sprintf("%f", item.(float64)))
		}
	}
}

func createUserVk() (client *api.VKClient){
	login := flag.String("login", "", "login string")
	password := flag.String("password", "", "password string")
	flag.Parse()

	client, err := api.NewVKClient(3, *login, *password)
	if err != nil {
		log.Fatal(err)
	}
	return
}