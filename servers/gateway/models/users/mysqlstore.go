package users

import (
	"database/sql"
	"fmt"
	"strconv"
)

const (
	insert      = "insert into users (username,passhash,firstname,lastname,photourl) values ( ?,?,?,?,?,? )"
	selectID    = `Select * From users Where id=?`
	getUserName = `Select * From users Where username=?`
	update      = "update users set role=? where id=?"
	del         = "delete from users where id=?"
)

//MySQLStore represents a user.Store backed by MySQL
type MySQLStore struct {
	db *sql.DB
}

//NewMySQLStore constructs a new MySQLStore. It will
//panic if the db pointer is nil.
func NewMySQLStore(db *sql.DB) *MySQLStore {
	if db == nil {
		panic("nil database pointer")
	}
	return &MySQLStore{db}
}

//Insert inserts a user into the database, returning
//the inserted User with its ID field set to the
//new primary key value
func (s *MySQLStore) Insert(user *User) (*User, error) {
	results, err := s.db.Exec(insert,
		user.UserName, user.PassHash, user.FirstName, user.LastName, user.PhotoURL)
	if err != nil {
		return nil, fmt.Errorf("executing insert: %v", err)
	}
	//get the new DBMS-generated primary key value
	id, err := results.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("getting new ID: %v", err)
	}
	//set the ID field of the struct so that callers
	//know what the new ID is
	user.ID = id
	return user, nil
}

//Insert inserts a family into the database, returning
//the inserted User with its ID field set to the
//new primary key value
func (s *MySQLStore) InsertFam(family *FamilyRoom) (*FamilyRoom, error) {
	results, err := s.db.Exec("insert into familyroom (name) values (?)", family.RoomName)
	if err != nil {
		return nil, fmt.Errorf("executing insert: %v", err)
	}
	//get the new DBMS-generated primary key value
	id, err := results.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("getting new ID: %v", err)
	}
	//set the ID field of the struct so that callers
	//know what the new ID is
	family.ID = id
	return family, nil
}

func getByHelper(s *MySQLStore, identifier string, command string) (*User, error) {
	var row *sql.Row
	if command == selectID {
		i, err := strconv.Atoi(identifier)
		if err != nil {
			row = s.db.QueryRow(command, i)
		}
	} else {
		row = s.db.QueryRow(command, identifier)
	}
	user := &User{}
	if err := row.Scan(&user.ID, &user.UserName, &user.PassHash,
		&user.FirstName, &user.LastName, &user.PhotoURL); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("scanning: %v", err)
	}
	return user, nil
}

//GetByID returns a specific user according to given id or ErrUserNotFound
//if the requested user does not exist
func (s *MySQLStore) GetByID(id int64) (*User, error) {
	return getByHelper(s, selectID, string(id))
}

//GetByUserName returns a specific user according to the given username, or ErrNotFound
//if the requested user does not exist
func (s *MySQLStore) GetByUserName(username string) (*User, error) {
	return getByHelper(s, getUserName, username)
}

//Update updates a user to the given user ID
//and returns the newly-inserted User. It returns
//nil and ErrUserNotFound if the task ID does not exist.
func (s *MySQLStore) Update(id int64, updates *Updates) (*User, error) {
	results, err := s.db.Exec(update, updates.Role, id)
	if err != nil {
		return nil, fmt.Errorf("updating: %v", err)
	}
	affected, err := results.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("getting rows affected: %v", err)
	}

	//if no rows were affected, then the requested
	//ID was not in the database
	if affected == 0 {
		return nil, ErrUserNotFound
	}
	return s.GetByID(id)
}

//Delete deletes the user with the given ID
func (s *MySQLStore) Delete(id int64) error {
	_, err := s.db.Exec(del, id)
	if err != nil {
		return fmt.Errorf("deleting: %v", err)
	}
	return nil
}