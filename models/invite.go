package models

import "github.com/api/database"

type Invite struct {
	Id           int    `json:"id"`
	Target       int    `json:"target"`
	InvitingRoom int    `json:"inviting_room"`
	Status       string `json:"status"`
}

func SelectInviteByUser(u User) []Invite {
	var db = database.GetConnection()
	var invites = make([]Invite, 8)

	search, err := db.Query(database.SearchInviteByTarget, u.Id)
	if err != nil {
		panic(err)
	}

	for i := 0; search.Next(); i++ {
		invites = append(invites, Invite{})
		err := search.Scan(&invites[i].Id, &invites[i].Target, &invites[i].InvitingRoom, &invites[i].Status)
		if err != nil {
			panic(err)
		}
	}

	return invites
}
