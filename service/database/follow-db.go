package database

import (
	"database/sql"
	"errors"
)

// Database function that adds a follow to another user
func (db *appdbimpl) InsertFollow(f Follow) (Follow, error) {
	res, err := db.c.Exec("INSERT INTO followers (userId, followerId) VALUES (?,?)", f.UserId, f.UserFollowedID)
	if err != nil {
		return f, err
	}
	lastInsertID, err := res.LastInsertId()

	if err != nil {
		return f, err
	}
	f.Id = uint64(lastInsertID)
	return f, nil

}

// Database function that removes a follow to another user
func (db *appdbimpl) RemoveFollow(FollowId uint64, UserId uint64, FollowedId uint64) error {
	res, err := db.c.Exec("DELETE FROM followers WHERE id = ? AND userId = ? AND followerId = ?", FollowId, UserId, FollowedId)

	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return ErrFollowDoesNotExist
	}
	return nil
}

func (db *appdbimpl) RemoveFollows(UserId uint64, FollowedId uint64) error {
	res, err := db.c.Exec("DELETE FROM followers WHERE userId = ? AND followerId = ?", UserId, FollowedId)
	res1, err1 := db.c.Exec("DELETE FROM followers WHERE userId = ? AND followerId = ?", FollowedId, UserId)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return ErrFollowDoesNotExist
	}
	if err1 != nil {
		return err1
	}
	affected1, err := res1.RowsAffected()
	if err1 != nil {
		return err1
	}
	if affected1 == 0 {
		return ErrFollowDoesNotExist
	}
	return nil
}

// Database function that returns a user's count of followers
func (db *appdbimpl) GetFollowersCount(id uint64) (int, error) {
	var count int

	if err := db.c.QueryRow("SELECT COUNT(*) FROM followers WHERE followerId = ?", id).Scan(&count); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, nil
		}
		return count, err
	}
	return count, nil
}

// Database function that returns a user's count of followed
func (db *appdbimpl) GetFollowedCount(id uint64) (int, error) {
	var count int

	if err := db.c.QueryRow("SELECT COUNT(*) FROM followers WHERE userId = ?", id).Scan(&count); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, nil
		}
		return count, err
	}
	return count, nil
}

// Database function that returns the users follow
func (db *appdbimpl) GetFollower(u User, token uint64) (Follow, error) {
	var follow Follow
	if err := db.c.QueryRow("SELECT Id, userId, followerId from followers WHERE followerId = ? AND userId = ?", u.Id, token).Scan(&follow.Id, &follow.UserId, &follow.UserFollowedID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return follow, ErrFollowDoesNotExist
		}
	}

	return follow, nil
}
