package models

type User struct {
	UID      int64  `json:"uid" db:"user_id" `
	UserName string `json:"username" db:"username"`
	Pwd      string `json:"pwd" db:"password"`
	Email    string `json:"email" db:"email"`
	Token    string
}
