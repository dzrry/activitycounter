package main

import (
	"fmt"
	"github.com/dzrry/activitycounter/vk/api"
	"log"
)

func checkMembers(client *api.VKClient) {
	tmpl := `
var groupId = %d;
var countByRequest = %d;
var hasNext = %t;
var offset = %d;
var calls = 0;

var items = [];
while (hasNext && calls < 25) {
	calls = calls + 1;
	var members = API.groups.getMembers({"group_id":groupId,"count":countByRequest,"offset":offset});
	items = members.items;

	offset = offset + countByRequest;
	if (offset >= members.count) {
		hasNext = false;
	}
	
}
return {"items":items,"has_next":hasNext,"offset":offset};`

	offset := 0
	countByRequest := 1000
	hasNext := true

	for {
		code := fmt.Sprintf(tmpl, rbkId, countByRequest, hasNext, offset)
		resp, err := client.Execute(code)
		if err != nil {
			log.Fatal(err)
		}

		hasNext = resp["has_next"].(bool)
		offset = int(resp["offset"].(float64))
		fmt.Println(offset)
		respItems := resp["items"].([]interface{})
		if len(respItems) == 0 && !hasNext {
			break
		}

		for _, item := range respItems {
			fmt.Println(item)
		}
	}
}
