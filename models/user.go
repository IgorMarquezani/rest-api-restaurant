package models

type User struct {
  Id    int `json:"id"` 
  Name  string `json:"name"`
  Email string `json:"email"`
  Passwd string `json:"passwd"`
  Img   []byte `json:"img"`
}
