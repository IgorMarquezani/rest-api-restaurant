package models

type Guest struct {
	Invinting_room   int `json:"invinting_room"`
	User_Id          int `json:"user_id"`
	Permission_level int `json:"permission_level"`
}

