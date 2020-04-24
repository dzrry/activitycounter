package postgres

import (
	"fmt"
	vkapi "github.com/himidori/golang-vk-api"
	"time"
)

func (db *DB) AllGroupMembers(groupId int) (int, []*vkapi.User, error) {
	rows, err := db.Query("SELECT user_id FROM users_activities WHERE group_id = $1", groupId)
	if err != nil {
		return 0, []*vkapi.User{}, err
	}
	defer rows.Close()
	var users []*vkapi.User

	for rows.Next() {
		u := &vkapi.User{}
		err := rows.Scan(&u.UID)
		if err != nil {
			fmt.Println(err)
			continue
		}
		users = append(users, u)
	}

	return len(users), users, nil
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
