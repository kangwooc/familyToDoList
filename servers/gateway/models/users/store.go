package users

import (
	"errors"
)

//ErrUserNotFound is returned when the user can't be found
var ErrUserNotFound = errors.New("user not found")

//Store represents a store for Users
type Store interface {
	//GetByID returns the User with the given ID
	GetByID(id int64) (*User, error)

	//GetByUserName returns the User with the given Username
	GetByUserName(username string) (*User, error)

	//Insert inserts the user into the database, and returns
	//the newly-inserted User, complete with the DBMS-assigned ID
	Insert(user *User) (*User, error)

	//Insert inserts the user into the database, and returns
	//the newly-inserted User, complete with the DBMS-assigned ID
	InsertFam(family *FamilyRoom) (*FamilyRoom, error)

	UpdateToMember(id int64, updates *Updates) (*User, error)

	GetRoomName(id int64) (*FamilyRoom, error)

	GetByRoomName(roomname string) ([]*User, error)
	//Update applies UserUpdates to the given user ID
	//and returns the newly-updated user
	Update(id int64, updates *Updates) (*User, error)

	//Delete deletes the user with the given ID
	Delete(id int64) error

	GetAdmin(roomname string, role string) (*User, error)
	// update the score after receiving points from task
	UpdateScore(id int64, point int) (*User, error)
}
