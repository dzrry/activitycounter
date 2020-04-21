package postgres

import (
	"fmt"
	vkapi "github.com/himidori/golang-vk-api"
	"github.com/lib/pq"
)

func (db *DB) AllGroupMembers(groupId int) (int, []*vkapi.User, error) {
	var count int
	userIds := make([]uint8, 0)
	row := db.QueryRow("SELECT count, user_ids FROM group_members WHERE group_id = $1", groupId)
	err := row.Scan(&count, &userIds)
	if err != nil {
		return 0, nil, err
	}
	fmt.Println(len(userIds))
	users := make([]*vkapi.User, 0)
	for _, id := range userIds {
		u := &vkapi.User{UID: int(id)}
		users = append(users, u)
	}
	return count, users, nil
}

func (db *DB) InsertGroupMembers(groupId, count int, users []*vkapi.User) (error) {
	userIds := make([]int, len(users), len(users))
	for i := range users {
		userIds[i] = users[i].UID
	}
	_, err := db.Exec("INSERT INTO group_members (group_id, count, user_ids) VALUES ($1, $2, $3)",
		groupId, count, pq.Array(userIds))
	if err != nil {
		return err
	}
	return err
}

