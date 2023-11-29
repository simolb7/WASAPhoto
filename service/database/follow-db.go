package database

import "database/sql"

// Database function that adds a follow to another user
func (db *appdbimpl) InsertFollow(f Follow) (Follow, error) {
	_, err := db.c.Exec("INSERT INTO followers (id, userId, userFollowed) VALUES (?,?,?)", f.Id, f.UserId, f.UserFollowedID)
	if err != nil {
		return f, err
	}
	return f, nil

}

// Database function that removes a follow to another user
func (db *appdbimpl) RemoveFollow(f Follow) error {
	res, err := db.c.Exec("DELETE FROM followers WHERE id = ? AND userId = ? AND userFollowed = ?", f.Id, f.UserId, f.UserFollowedID)

	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	} else if affected == 0 {
		return ErrFollowDoesNotExist
	}
	return nil
}

// Database function that returns a user's count of followers
func (db *appdbimpl) GetFollowersCount(id uint64) (int, error) {
	var count int

	if err := db.c.QueryRow("SELECT COUNT(*) FROM followers WHERE userFollowed = ?", id).Scan(&count); err != nil {
		if err == sql.ErrNoRows {
			return count, ErrFollowDoesNotExist
		}
	}
	return count, nil
}

// Database function that returns a user's count of followed
func (db *appdbimpl) GetFollowedCount(id uint64) (int, error) {
	var count int

	if err := db.c.QueryRow("SELECT COUNT(*) FROM followers WHERE userId = ?", id).Scan(&count); err != nil {
		if err == sql.ErrNoRows {
			return count, ErrFollowDoesNotExist
		}
	}
	return count, nil
}

// Database function that returns the users follow
func (db *appdbimpl) GetFollowers(u User, token uint64) (Follow, error) {
	var follow Follow
	if err := db.c.QueryRow("SELECT Id, userId, userFollowed from followers WHERE userFollowed = ?", u.Id, token).Scan(&follow.Id, &follow.UserFollowedID, &follow.UserId); err != nil {
		if err == sql.ErrNoRows {
			return follow, ErrFollowDoesNotExist
		}
	}

	return follow, nil
}
