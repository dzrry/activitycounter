package postgres

import (
	vkapi "github.com/himidori/golang-vk-api"
)

func (db *DB) AllGroupMembers(groupId int) (int, []*vkapi.User, error) {
	var count int
	users := make([]*vkapi.User, 0)
	row := db.QueryRow("SELECT * FROM group_members WHERE id = $1", groupId)
	err := row.Scan(&count, &users)
	if err != nil {
		return 0, nil, err
	}
	return count, users, nil
}

func (db *DB) InsertGroupMembers(groupId, count int, users []*vkapi.User) (int64, error) {
	userIds := make([]int, len(users), cap(users))
	for i := range users {
		userIds[i] = users[i].UID
	}
	res, err := db.Exec("INSERT INTO group_members (group_id, count, user_ids) VALUES ($1, $2, $3)",
		groupId, count, userIds)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, err
}

