package models

type Credentials struct {
	ID       string
	Name     string `json:"name"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Version  int    `json:"version"`
}
