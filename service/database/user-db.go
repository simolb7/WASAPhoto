package database

import (
	"database/sql"
)

// Database function that adds a new user in the database upon registration
func (db *appdbimpl) CreateUser(u User) (User, error) {
	res, err := db.c.Exec("INSERT INTO users(username) VALUES (?)", u.Username)

	if err != nil {
		var user User
		if err := db.c.QueryRow("SELECT id, Username FROM users WHERE username = ?", u.Username).Scan(&user.Id, &user.Username); err != nil {
			if err == sql.ErrNoRows {
				return user, ErrUserDoesNotExist
			}
		}
		return user, nil
	}

	LastInsertID, err := res.LastInsertId()
	if err != nil {
		return u, err
	}
	u.Id = uint64(LastInsertID)
	return u, nil

}

// Database function that allow the user to change username
func (db *appdbimpl) SetUsername(u User, username string) (User, error) {
	res, err := db.c.Exec("UPDATE users SET Username=? WHERE id=? AND Username=?", u.Username, u.Id, username)

	if err != nil {
		return u, err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return u, err
	} else if affected == 0 {
		return u, err
	}
	return u, nil
}

// Database function that return the User entity
func (db *appdbimpl) GetUserId(username string) (User, error) {
	var user User
	if err := db.c.QueryRow("SELECT id, username FROM users WHERE username = ?", username).Scan(&user.Id, &user.Username); err != nil {
		if err == sql.ErrNoRows {
			return user, ErrUserDoesNotExist
		}
	}
	return user, nil
}

func (db *appdbimpl) CheckUserByUsername(u User) (User, error) {
	var user User
	if err := db.c.QueryRow(`SELECT id, username FROM users WHERE username = ?`, u.Username).Scan(&user.Id, &user.Username); err != nil {
		if err == sql.ErrNoRows {
			return user, ErrUserDoesNotExist
		}
	}
	return user, nil
}

func (db *appdbimpl) CheckUserById(u User) (User, error) {
	var user User
	if err := db.c.QueryRow(`SELECT id, username FROM users WHERE id = ?`, u.Id).Scan(&user.Id, &user.Username); err != nil {
		if err == sql.ErrNoRows {
			return user, ErrUserDoesNotExist
		}
	}
	return user, nil
}

func (db *appdbimpl) CheckUser(u User) (User, error) {
	var user User
	if err := db.c.QueryRow(`SELECT id, username FROM users WHERE id = ? AND username = ?`, u.Id, u.Username).Scan(&user.Id, &user.Username); err != nil {
		if err == sql.ErrNoRows {
			return user, ErrUserDoesNotExist
		}
	}
	return user, nil
}

// Database function that gets the stream of a user
func (db *appdbimpl) GetStream(u User) ([]PhotoStream, error) {
	var ret []PhotoStream
	rows, err := db.c.Query("SELECT Id, userId, photo, date FROM photos WHERE userId IN (SELECT followerId FROM followers WHERE userId=? AND followerId NOT IN (SELECT userId FROM bans WHERE bannedId=?))")

	if err != nil {
		return ret, ErrUserDoesNotExist
	}

	defer func() {
		_ = rows.Close()
	}()

	for rows.Next() {
		var b PhotoStream
		err = rows.Scan(&b.Id, &b.UserID, &b.File, &b.Date)
		if err != nil {
			return nil, err
		}
		if err := db.c.QueryRow(`SELECT COUNT(*) FROM likes WHERE photoId = ?`, b.Id).Scan(&b.LikeNumber); err != nil {
			if err == sql.ErrNoRows {
				return nil, err
			}
		}
		if err := db.c.QueryRow(`SELECT COUNT(*) FROM comments WHERE photoId = ?`, b.Id).Scan(&b.CommentNumber); err != nil {
			if err == sql.ErrNoRows {
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

// Database function that return true, if one user follows another user, false otherwise
func (db *appdbimpl) GetFollowStatus(user uint64, followed uint64) (bool, error) {
	var ret bool
	if err := db.c.QueryRow(`SELECT EXISTS(SELECT 1 FROM followers WHERE userId= ? AND  followerId= ?)`, user, followed).Scan(&ret); err != nil {
		if err == sql.ErrNoRows {
			return false, err
		}
	}
	return ret, nil
}

// Database function that return true, if one user have banned another user, false otherwise
func (db *appdbimpl) GetBanStatus(user uint64, banned uint64) (bool, error) {
	var ret bool
	if err := db.c.QueryRow(`SELECT EXISTS(SELECT 1 FROM bans WHERE userId=? AND bannedId=?)`, user, banned).Scan(&ret); err != nil {
		if err == sql.ErrNoRows {
			return false, err
		}
	}
	return ret, nil
}
