package database

import "database/sql"

// Database function that adds a like of a user to a photo
func (db *appdbimpl) InsertLike(l Like) (Like, error) {
	_, err := db.c.Exec("INSERT INTO likes (id, userId, photoid, photoOwner) VALUES (?,?,?,?)", l.LikeId, l.UserId, l.PhotoId, l.PhotoOwner)
	if err != nil {
		return l, err
	}
	return l, nil

}

// Database function that removes a like of a user to a photo
func (db *appdbimpl) RemoveLike(l Like) error {
	res, err := db.c.Exec("DELETE FROM likes WHERE id = ? AND photoId = ? AND userId = ?", l.LikeId, l.PhotoId, l.UserId)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	} else if affected == 0 {
		return ErrLikeDoesNotExist
	}
	return nil
}

// Database function that removes each like of a user due to ban
func (db *appdbimpl) RemoveLikes(user uint64, banned uint64) error {
	_, err := db.c.Exec("DELETE FROM likes WHERE userId = ? AND photoOwner = ?", user, banned)
	if err != nil {
		return err
	}
	return nil
}

// Database function that returns the count of like on a photo
func (db *appdbimpl) GetLikeCount(photoid uint64) (int, error) {
	var count int

	if err := db.c.QueryRow("SELECT COUNT(*) FROM likes WHERE photoId = ?", photoid).Scan(&count); err != nil {
		if err == sql.ErrNoRows {
			return count, ErrLikeDoesNotExist
		}
	}
	return count, nil
}
