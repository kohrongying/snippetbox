package models

import (
	"database/sql"
	"time"
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
	
	// stmt := `INSERT INTO snippets (title, content, created, expires) 
	// VALUES($1, $2, NOW(), NOW() + INTERVAL '1 DAY' * $3)
	// RETURNING id;`

	// var id int
	// err := m.DB.QueryRow(stmt, title, content, expires).Scan(&id)
	// if err != nil {
	// 	return 0, err
	// }

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