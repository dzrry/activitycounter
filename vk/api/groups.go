package api

import (
	"encoding/json"
	"errors"
	"net/url"
	"strconv"
)

type GroupMembers struct {
	Count   int     `json:"count"`
	Members []*User `json:"items"`
}

type User struct {
	UID int `json:"id"`
}

// TODO: add parameters: sort, offset, count, fields, filter
func (client VK) GroupGetMembers(groupId int) (int, []*User, error) {
	params := url.Values{}
	params.Set("group_id", strconv.Itoa(groupId))

	resp, err := client.MakeRequest("groups.getMembers", params)
	if err != nil {
		return 0, nil, errors.New("groups 26 " + err.Error())
	}

	var gm *GroupMembers
	json.Unmarshal(resp.Response, &gm)
	return gm.Count, gm.Members, nil
}
