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
	var invites = make([]Invite, 0)

	rows, err := db.Query(database.SearchInviteByTarget, u.Id)
	if err != nil {
		panic(err)
	}
  defer rows.Close()

	for i := 0; rows.Next(); i++ {
		invites = append(invites, Invite{})
		err := rows.Scan(&invites[i].Id, &invites[i].Target, &invites[i].InvitingRoom, &invites[i].Status)
		if err != nil {
			panic(err)
		}
	}

	return invites
}
