package users

import (
	"errors"
	"fmt"
)

func makeNewUser() *User {
	newUser := &NewUser{
		Password:     "test1234",
		PasswordConf: "test1234",
		UserName:     "test",
		FirstName:    "juan",
		LastName:     "oh",
	}
	user, err := newUser.ToUser()
	if err != nil {
		fmt.Printf("unexpected error: %v", err)
		return nil
	}
	return user
}

//MockStore is fake
type MockStore struct {
}

//NewMockStore is fake
func NewMockStore() *MockStore {
	return &MockStore{}
}

//Insert inserts a user information into the database, returning
//the inserted user information with its ID field set to the
//new primary key value
func (s *MockStore) Insert(u *User) (*User, error) {

	return u, nil
}

//GetByID returns a specific user by passing id, or ErrNotFound
//if the requested user does not exist
func (s *MockStore) GetByID(id int64) (*User, error) {
	if id == int64(1) {
		return nil, errors.New("Fake Error")
	}
	return makeNewUser(), nil
}

//GetByUserName returns a specific user by passing username, or ErrNotFound
//if the requested user does not exist
func (s *MockStore) GetByUserName(username string) (*User, error) {
	return makeNewUser(), nil
}

//Update updates a users, setting only the completed state,
//and returns a copy of the updated user. It returns
//nil and ErrNotFound if the user ID does not exist.
func (s *MockStore) Update(id int64, updates *Updates) (*User, error) {
	if id == 2 {
		return nil, errors.New("fake errors")
	}
	return makeNewUser(), nil
}

//Delete selected completed users and returns
//the number of users that were deleted.
func (s *MockStore) Delete(id int64) error {
	return nil
}
