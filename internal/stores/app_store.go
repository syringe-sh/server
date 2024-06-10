package stores

import (
	"database/sql"
	"errors"
)

type User struct {
	Id        int
	Username  string
	Email     string
	CreatedAt string
	Status    string
}

type Key struct {
	Id        int
	PublicKey string
	UserId    int
	CreatedAt string
}

type UserStore interface {
	InsertUser(username, email, status string) (*User, error)
	GetUserByUsername(username string) (*User, error)
	DeleteUserByUsername(username string) error
	UpdateUser(user User) (*User, error)
	GetUserPublicKeys(username string) (*[]Key, error)
}

type KeyStore interface {
	InsertKey(userId int, publicKey string) (*Key, error)
}

type AppStore interface {
	UserStore
	KeyStore
}

type SqliteAppStore struct {
	appDb *sql.DB
}

func NewSqliteAppStore(appDb *sql.DB) SqliteAppStore {
	return SqliteAppStore{appDb}
}

func (s SqliteAppStore) InsertUser(username, email, status string) (*User, error) {
	query := `
		insert into users_ (username_, email_, status_) 
		values ($username, $email, $status) 
		returning id_, username_, email_, status_, created_at_
	`

	row := s.appDb.QueryRow(
		query,
		sql.Named("username", username),
		sql.Named("email", email),
		sql.Named("status", status),
	)

	var insertedUser User

	if err := row.Scan(
		&insertedUser.Id,
		&insertedUser.Username,
		&insertedUser.Email,
		&insertedUser.Status,
		&insertedUser.CreatedAt,
	); err != nil {
		return nil, err
	}

	return &insertedUser, nil
}

func (s SqliteAppStore) GetUserByUsername(username string) (*User, error) {
	query := `
		select id_, username_, email_, status_, created_at_ 
		from users_ 
		where username_ = $1
	`

	row := s.appDb.QueryRow(query, username)

	var user User

	if err := row.Scan(
		&user.Id,
		&user.Username,
		&user.Email,
		&user.Status,
		&user.CreatedAt,
	); err != nil {
		return nil, err
	}

	return &user, nil
}

func (s SqliteAppStore) DeleteUserByUsername(username string) error {
	query := `delete from users_ where username_ = $1`

	res, err := s.appDb.Exec(query, username)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("no user deleted")
	}

	return nil
}

func (s SqliteAppStore) UpdateUser(user User) (*User, error) {
	query := `
		update users_ set email_ = $2, set status_ = $3
		where username_ = $1 
		returning id_, username_, email_, status_, created_at_
	`

	row := s.appDb.QueryRow(query, user.Username, user.Email, user.Status)

	var updatedUser User

	if err := row.Scan(
		&updatedUser.Id,
		&updatedUser.Username,
		&updatedUser.Email,
		&updatedUser.Status,
		&updatedUser.CreatedAt,
	); err != nil {
		return nil, err
	}

	return &user, nil
}

func (s SqliteAppStore) InsertKey(userId int, publicKey string) (*Key, error) {
	query := `
		insert into keys_ (user_id_, ssh_public_key_)
		values ($userId, $publicKey)
		returning id_, user_id_, ssh_public_key_, created_at_
	`

	row := s.appDb.QueryRow(
		query,
		sql.Named("userId", userId),
		sql.Named("publicKey", publicKey),
	)

	var insertedKey Key

	if err := row.Scan(
		&insertedKey.Id,
		&insertedKey.UserId,
		&insertedKey.PublicKey,
		&insertedKey.CreatedAt,
	); err != nil {
		return nil, err
	}

	return &insertedKey, nil
}

func (s SqliteAppStore) GetUserPublicKeys(username string) (*[]Key, error) {
	query := `
		select k.id_, k.user_id_, k.ssh_public_key_, k.created_at_
		from keys_ k 
		inner join
		users_ u
		on k.user_id_ = u.id_
		where u.username_ = $username
	`

	rows, err := s.appDb.Query(
		query,
		sql.Named("username", username),
	)
	if err != nil {
		return nil, err
	}

	var keys []Key

	for rows.Next() {
		var key Key

		if err := rows.Scan(
			&key.Id,
			&key.UserId,
			&key.PublicKey,
			&key.CreatedAt,
		); err != nil {
			return nil, err
		}

		keys = append(keys, key)
	}

	return &keys, nil
}
