package repository

type UserEntity struct {
}

type User interface {
}

type user struct {
}

func NewUser() User {
	return &user{}
}
