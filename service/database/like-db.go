package database

import (
	"database/sql"
	"errors"
)

// Database function that adds a like of a user to a photo
func (db *appdbimpl) InsertLike(l Like) (Like, error) {
	res, err := db.c.Exec("INSERT INTO likes (userId, photoid, photoOwner) VALUES (?,?,?)", l.UserId, l.PhotoId, l.PhotoOwner)
	if err != nil {
		return l, err
	}
	// Get the last inserted ID directly from the database
	lastInsertID, err := res.LastInsertId()
	if err != nil {
		return l, err
	}

	// Assign the generated ID to the like
	l.LikeId = uint64(lastInsertID)
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
	}
	if affected == 0 {
		return ErrLikeDoesNotExist
	}
	return nil
}

// Database function that removes each like of a user due to ban
func (db *appdbimpl) RemoveLikes(user uint64, banned uint64) error {
	_, err := db.c.Exec("DELETE FROM likes WHERE userId = ? AND photoOwner = ?", banned, user)
	if err != nil {
		return err
	}
	_, err = db.c.Exec("DELETE FROM likes WHERE userId = ? AND photoOwner = ?", user, banned)
	if err != nil {
		return err
	}

	return nil
}

// Database function that returns the count of like on a photo
func (db *appdbimpl) GetLikeCount(photoid uint64) (int, error) {
	var count int

	err := db.c.QueryRow("SELECT COUNT(*) FROM likes WHERE photoId = ?", photoid).Scan(&count)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, nil
		}
		return 0, err
	}
	return count, nil
}

// Function returns a like
// il token viene usato come id
func (db *appdbimpl) GetLike(photoid uint64, token uint64) (Like, error) {
	var like Like

	err := db.c.QueryRow(`SELECT Id, userId, photoId, photoOwner FROM likes WHERE userId = ? AND photoId = ?`, token, photoid).Scan(&like.LikeId, &like.UserId, &like.PhotoId, &like.PhotoOwner)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return like, ErrLikeDoesNotExist
		}
		return like, err
	}
	return like, nil
}
