package models

//type User struct {
//	ID       uint64
//	Email    string `json:"email"`
//	Password string `json:"password"`
//}

type User struct {
	ID       uint64
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Name     string `json:"name"`
	Pass     string `json:"pass"`
	Birthday uint64 `json:"birthday"`
}

type UserWithoutPass struct {
	ID       uint64
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Name     string `json:"name"`
	Birthday uint64 `json:"birthday"`
}

type UserWithoutID struct {
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Name     string `json:"name"`
	Pass     string `json:"pass"`
	Birthday uint64 `json:"birthday"`
}

type PreUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
