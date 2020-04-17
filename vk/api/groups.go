package api

type GroupMembers struct {
	Count   int     `json:"count"`
	Members []*User `json:"items"`
}

type User struct {
	UID int `json:"id"`
}
