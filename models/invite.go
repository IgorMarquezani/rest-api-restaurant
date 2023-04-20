package models

import "github.com/api/database"

const (
	ErrInvitedAlready    string = "This user have been already invited to the room"
	ErrNoSuchInvite      string = "This invite does not exist"
	ErrInvalidPermission string = "This permission level does not exist"
)

type Invite struct {
	Id           int    `json:"id"`
	Target       int    `json:"target"`
	InvitingRoom int    `json:"inviting_room"`
	Status       string `json:"status"`
	Permission   uint   `json:"permission"`
}

func InsertInvite(user User, roomId int, permission uint) {
	db := database.GetConnection()

	_, err := db.Query(database.InsertInvite, user.Id, roomId, permission)
	if err != nil {
		panic(err)
	}
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
		err := rows.Scan(&invites[i].Id, &invites[i].Target, &invites[i].InvitingRoom, &invites[i].Status, &invites[i].Permission)
		if err != nil {
			panic(err)
		}
	}

	return invites
}

func SelectInvite(user User, room Room) (Invite, error) {
	var db = database.GetConnection()
	var invite Invite

	rows, err := db.Query(database.SearchInvite, user.Id, room.Id)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	if rows.Next() {
		err := rows.Scan(&invite.Id, &invite.Target, &invite.InvitingRoom, &invite.Status, &invite.Permission)
		if err != nil {
			panic(err)
		}
	}

	return invite, nil
}

func SearchForInvites(v []Invite, id int) (Invite, bool) {
	for _, invite := range v {
		if invite.Id == id {
			return invite, true
		}
	}

	return Invite{}, false
}
