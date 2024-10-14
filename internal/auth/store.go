package auth

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/Tutuacs/pkg/db"

	"github.com/Tutuacs/internal/user"
)

type Store struct {
	db      *sql.DB
	extends bool
	Table   string
}

func NewStore(conn ...*sql.DB) (*Store, error) {
	if len(conn) == 0 {
		con, err := db.NewConnection()
		if err != nil {
			return nil, err
		}
		return &Store{db: con, extends: false, Table: "users"}, nil
	}
	return &Store{db: conn[0], extends: true, Table: "users"}, nil
}

func (s *Store) CloseStore() {
	if !s.extends {
		s.db.Close()
	}
}

func (s *Store) GetUserByEmail(email string) (usr *user.User, err error) {
	err = nil
	usr = &user.User{}

	query := "SELECT * FROM " + s.Table + " WHERE email = $1"
	row := s.db.QueryRow(query, email)

	db.ScanRow(row, usr)

	if usr.ID == 0 {
		err = fmt.Errorf("user not found")
		return
	}

	return
}

func (s *Store) GetUserByID(ID int) (*user.User, error) {

	sql := "SELECT * FROM users WHERE id = $1"

	rows, err := s.db.Query(sql, ID)
	if err != nil {
		return nil, err
	}

	usr := new(user.User)

	for rows.Next() {
		err = db.ScanRows(rows, usr)
		if err != nil {
			return nil, err
		}
	}

	return usr, err
}

func (s *Store) CreateUser(user user.User) error {
	query := "INSERT INTO " + s.Table + " (name, email, password) VALUES ($1, $2, $3)"
	_, err := s.db.Exec(query, strings.Split(user.Email, "@")[0], user.Email, user.Password)
	return err
}

func (s *Store) GetLogin(email string) (int64, string, string, error) {
	var userEmail, password string
	var id int64

	query := fmt.Sprintf("SELECT email, password FROM %s WHERE email = $1", s.Table)
	err := s.db.QueryRow(query, email).Scan(&id, &userEmail, &password)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, "", "", fmt.Errorf("user not found")
		}
		return 0, "", "", err
	}

	return id, userEmail, password, nil
}
