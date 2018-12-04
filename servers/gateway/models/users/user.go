package users

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

//gravatarBasePhotoURL is the base URL for Gravatar image requests.
const gravatarBasePhotoURL = "https://www.gravatar.com/avatar/"

//bcryptCost is the default bcrypt cost to use when hashing passwords
var bcryptCost = 13

//User represents a user account in the database
type User struct {
	ID        int64  `json:"id"`
	PassHash  []byte `json:"-"` //never JSON encoded/decoded
	UserName  string `json:"username"`
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	PhotoURL  string `json:"photourl"`
	Role      string `json:"personrole"`
	RoomName  string `json:"roomname"`
}

//Credentials represents user sign-in credentials
type Credentials struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}

//NewUser represents a new user signing up for an account
type NewUser struct {
	Password     string `json:"password"`
	PasswordConf string `json:"passwordConf"`
	UserName     string `json:"userName"`
	FirstName    string `json:"firstName"`
	LastName     string `json:"lastName"`
}

//Updates represents allowed updates to a user profile
type Updates struct {
	Role     string `json:"personrole"`
	RoomName string `json:"roomname"`
}

//FamilyRoom represents family room table
type FamilyRoom struct {
	ID       int64  `json:"id"`
	RoomName string `json:"roomname"`
}

//Validate validates the new user and returns an error if
//any of the validation rules fail, or nil if its valid
func (nu *NewUser) Validate() error {
	if len(nu.Password) < 6 {
		return fmt.Errorf("Password must be at least 6 characters")
	}
	if nu.Password != nu.PasswordConf {
		return fmt.Errorf("Password and PasswordConf must match")
	}
	if len(nu.UserName) == 0 || strings.Contains(nu.UserName, " ") {
		return fmt.Errorf("UserName must be non-zero length and may not contain spaces")
	}
	if nu.FirstName == "" {
		return fmt.Errorf("FirstName must be non-zero length")
	}
	if nu.LastName == "" {
		return fmt.Errorf("LastName must be non-zero length")
	}
	return nil
}

//ToUser converts the NewUser to a USer, setting the
//PhotoURL and PassHash fields appropriately
func (nu *NewUser) ToUser() (*User, error) {
	err := nu.Validate()
	if err != nil {
		return nil, err
	}
	user := &User{}
	user.UserName = nu.UserName
	user.FirstName = nu.FirstName
	user.LastName = nu.LastName
	user.Role = "default"
	er := user.SetPassword(nu.Password)
	if er != nil {
		return nil, er
	}
	hasher := md5.New()
	user.PhotoURL = gravatarBasePhotoURL + hex.EncodeToString(hasher.Sum(nil))
	return user, nil
}

//FullName returns the user's full name, in the form:
// "<FirstName> <LastName>"
//If either first or last name is an empty string, no
//space is put between the names. If both are missing,
//this returns an empty string
func (u *User) FullName() string {
	return u.FirstName + " " + u.LastName
}

//SetPassword hashes the password and stores it in the PassHash field
func (u *User) SetPassword(password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
	if err != nil {
		return err
	}
	u.PassHash = hash
	return nil
}

//Authenticate compares the plaintext password against the stored hash
//and returns an error if they don't match, or nil if they do
func (u *User) Authenticate(password string) error {
	return bcrypt.CompareHashAndPassword(u.PassHash, []byte(password))
}

//ApplyUpdates applies the updates to the user. An error
//is returned if the updates are invalid
func (u *User) ApplyUpdates(updates *Updates) error {
	u.Role = updates.Role
	u.RoomName = updates.RoomName
	// log.Println("this is userrrrr %v", u)
	// log.Println("this is userrrrrole %v", u.Role)
	// log.Println("this is userrrrroom %v", u.RoomName)

	return nil
}
