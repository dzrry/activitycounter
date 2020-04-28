package postgres

import (
	"fmt"
	"github.com/dzrry/activitycounter/vk/api"
	"time"
)

func (db *DB) AllGroupMembers(groupId int) (int, []*api.User, error) {
	rows, err := db.Query("SELECT user_id FROM users_activities WHERE group_id = $1", groupId)
	if err != nil {
		return 0, []*api.User{}, err
	}
	defer rows.Close()
	var users []*api.User

	for rows.Next() {
		u := &api.User{}
		err := rows.Scan(&u.UID)
		if err != nil {
			fmt.Println(err)
			continue
		}
		users = append(users, u)
	}

	return len(users), users, nil
}

func (db *DB) InsertGroupMembers(groupId int, users []*api.User) error {
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
