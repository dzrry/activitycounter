package postgres

import (
	"fmt"
	vkapi "github.com/himidori/golang-vk-api"
	"time"
)

func (db *DB) AllGroupMembers(groupId int) (int, []*vkapi.User, error) {
	var count int
	userIds := make([]uint8, 0)
	row := db.QueryRow("SELECT count, user_ids FROM users_activities WHERE group_id = $1", groupId)
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

func (db *DB) InsertGroupMembers(groupId int, users []*vkapi.User) error {
	dt := time.Now().Format("01-02-2006 15:04:05")
	stmt, err := db.Prepare(
		"INSERT INTO users_activities (user_id, group_id, event, date) VALUES ($1, $2, 'subscribe', $3)")
	if err != nil {
		return err
	}
	defer stmt.Close()
	for i := range users {
		_, err := stmt.Exec(users[i].UID, groupId, dt)
		if err != nil {
			return err
		}
	}
	return nil
}

