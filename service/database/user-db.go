package database

import (
	"database/sql"
	"errors"
)

// Database function that adds a new user in the database upon registration, or logs the user

func (db *appdbimpl) CreateUser(u User) (User, error) {
	var user User
	err := db.c.QueryRow("SELECT id, Username FROM users WHERE username = ?", u.Username).Scan(&user.Id, &user.Username)
	if err == nil {
		// L'utente esiste già, fai login
		return user, err
	}
	if !errors.Is(err, sql.ErrNoRows) {
		// Un errore diverso da ErrNoRows si è verificato durante la query
		return user, ErrUserDoesNotExist
	}

	// creo utente
	res, err := db.c.Exec("INSERT INTO users(username) VALUES (?)", u.Username)
	if err != nil {
		return u, err
	}
	LastInsertID, err := res.LastInsertId()
	if err != nil {
		return u, err
	}
	u.Id = uint64(LastInsertID)
	return u, nil

}

// Database function that allow the user to change username
func (db *appdbimpl) SetUsername(u User, newusername string) (User, error) {
	var conflictingUser User
	err := db.c.QueryRow("SELECT Id, Username FROM users WHERE Username = ? AND Id <> ?", newusername, u.Id).Scan(&conflictingUser.Id, &conflictingUser.Username)
	if err == nil {
		// new username is already used
		return u, ErrUsernameAlreadyExists
	}
	if !errors.Is(err, sql.ErrNoRows) {
		// Un errore diverso si è verificato durante la query
		return u, err
	}
	// if it is all correct, change username
	res, err := db.c.Exec("UPDATE users SET Username=? WHERE Id=?", newusername, u.Id)

	if err != nil {
		return u, err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return u, err
	}
	if affected == 0 {
		return u, ErrUserDoesNotExist
	}

	u.Username = newusername
	return u, nil
}

// Database function that return the User entity
func (db *appdbimpl) GetUserId(username string) (User, error) {
	var user User
	err := db.c.QueryRow("SELECT id, username FROM users WHERE username = ?", username).Scan(&user.Id, &user.Username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// specific error
			return user, ErrUserDoesNotExist
		}

		// for manage general errors
		return user, err
	}
	return user, nil
}

// Check if an user exist by username
func (db *appdbimpl) CheckUserByUsername(u User) (User, error) {
	var user User
	if err := db.c.QueryRow(`SELECT id, username FROM users WHERE username = ?`, u.Username).Scan(&user.Id, &user.Username); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return user, ErrUserDoesNotExist
		}
	}
	return user, nil
}

// Return an user form an id
func (db *appdbimpl) GetUsername(id uint64) (User, error) {
	var user User
	if err := db.c.QueryRow(`SELECT id, username FROM users WHERE id = ?`, id).Scan(&user.Id, &user.Username); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return user, ErrUserDoesNotExist
		}
	}
	return user, nil
}

// Check user by id
func (db *appdbimpl) CheckUserById(u User) (User, error) {
	var user User
	err := db.c.QueryRow(`SELECT id, username FROM users WHERE id = ?`, u.Id).Scan(&user.Id, &user.Username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return user, ErrUserDoesNotExist
		}
	}
	return user, nil
}

// Check user
func (db *appdbimpl) CheckUser(u User) (User, error) {
	var user User
	if err := db.c.QueryRow(`SELECT id, username FROM users WHERE id = ? AND username = ?`, u.Id, u.Username).Scan(&user.Id, &user.Username); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return user, ErrUserDoesNotExist
		}
	}
	return user, nil
}

// Database function that gets the stream of a user
func (db *appdbimpl) GetStream(u User) ([]Photo, error) {

	rows, err := db.c.Query("SELECT Id, userId, photo, date FROM photos WHERE userId IN (SELECT followerId FROM followers WHERE userId=? AND followerId NOT IN (SELECT userId FROM bans WHERE bannedId=?))", u.Id, u.Id)

	if err != nil {
		return nil, ErrUserDoesNotExist
	}

	defer func() {
		_ = rows.Close()
	}()

	var ret []Photo
	for rows.Next() {
		var b Photo
		err = rows.Scan(&b.PhotoId, &b.UserId, &b.File, &b.Date)
		if err != nil {
			return nil, err
		}
		if err := db.c.QueryRow(`SELECT COUNT(*) FROM likes WHERE photoId = ?`, b.PhotoId).Scan(&b.LikeNumber); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, err
			}
		}
		if err := db.c.QueryRow(`SELECT COUNT(*) FROM comments WHERE photoId = ?`, b.PhotoId).Scan(&b.CommentNumber); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, err
			}
		}

		ret = append(ret, b)
	}

	if rows.Err() != nil {
		return nil, err
	}
	return ret, nil
}
