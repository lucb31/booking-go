package booking

import (
	"errors"
	"slices"
)

type User struct {
	Id         int
	Name       string
	Department string
}

var Users = []User{{1, "User A", "IT Department"}}

func GetUserById(id int) (*User, error) {
	idx := slices.IndexFunc(Users, func(user User) bool { return user.Id == id })
	if idx == -1 {
		return nil, errors.New("Unknown user id")
	}
	return &Users[idx], nil
}
