package models

import (
	"database/sql"
	"errors"
	"strings"
	"time"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        		int
	Name      		string
	Email     		string
	HashedPassword  []byte
	Created   		time.Time
}

type UserModel struct {
	DB *sql.DB
}

func (m *UserModel) Insert(name, email, password string) error {

	// create hash
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	stmt := `INSERT INTO users (name, email, hashed_password, created) 
	VALUES($1, $2, $3, NOW())`

	_, err = m.DB.Exec(stmt, name, email, hashedPassword)
	if err != nil {
		var pqError *pq.Error
		if errors.As(err, &pqError) {
			if pqError.Code == "23505" && strings.Contains(pqError.Message, "users_uc_email") {
				return ErrDuplicateEmail
			}
		}
		return err
	}

	return nil
}

func (m *UserModel) Authenticate(email, password string) (int, error) {
	// stmt := `SELECT id, title, content, created, expires FROM snippets
	// WHERE expires > NOW() and id = $1`

	// row := m.DB.QueryRow(stmt, id) // returns pointer to sql.Row object

	// s := &Snippet{} // Initialise pointer to a new zeroed Snippet struct

	// err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
	// if err != nil {
	// 	// if id not found, row.Scan returns sql.ErrNoRows error
	// 	if errors.Is(err, sql.ErrNoRows) {
	// 		return nil, ErrNoRecord
	// 	} else {
	// 		return nil, err
	// 	}
	// }
	
	return 1, nil
}

func (m *UserModel) Exists(id int) (bool, error) {
	return false, nil
}