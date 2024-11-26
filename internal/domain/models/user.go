package models

type User struct {
	ID       int64  `json:"id"`
	Email    string `json:"email"`
	PassHash []byte `json:"passHash"`
	Name     string `json:"name"`
	Image    string `json:"image"`
}
