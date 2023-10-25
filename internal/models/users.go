package models

type User struct {
	ID       uint64
	Email    string `json:"email"`
	Password string `json:"password"`
}

type PreUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
