package database

import (
	"database/sql"
	"errors"
)

// Database function to upload a photo
func (db *appdbimpl) InsertPhoto(p Photo) (Photo, error) {
	result, err := db.c.Exec("INSERT INTO photos (userId, photo, date) VALUES (?,?,?)", p.UserId, p.File, p.Date)
	if err != nil {
		return p, err
	}

	// generate id
	lastInsertID, err := result.LastInsertId()
	if err != nil {
		return p, err
	}

	p.PhotoId = uint64(lastInsertID)
	return p, nil

}

// Database function that remove a photo
func (db *appdbimpl) RemovePhoto(id uint64) error {
	_, err := db.c.Exec("DELETE FROM likes WHERE photoId = ?", id)
	if err != nil {
		return err
	}

	_, err1 := db.c.Exec("DELETE FROM comments WHERE photoId = ?", id)
	if err1 != nil {
		return err1
	}

	_, err2 := db.c.Exec("DELETE FROM photos WHERE Id = ?", id)
	if err2 != nil {
		return err2
	}
	return nil
}

// Database function that returns each photo of a specific user
func (db *appdbimpl) GetPhotos(u User) ([]Photo, error) {
	var ret []Photo
	rows, err := db.c.Query("SELECT Id, userId, photo, date FROM photos WHERE userId = ?", u.Id)

	if err != nil {
		return ret, ErrUserDoesNotExist
	}

	defer func() {
		_ = rows.Close()
	}()

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

	return ret, rows.Err()
}

// The function returns photos of users that a specific user follows, sorted in a descending order
func (db *appdbimpl) GetPhotosFollower(token uint64) ([]Photo, error) {
	var ret []Photo
	rows, err := db.c.Query("SELECT Id, userId, photo, date FROM photos WHERE userId IN (SELECT followerId FROM followers WHERE userId=?) ORDER BY date DESC LIMIT 10", token)

	if err != nil {
		return ret, ErrUserDoesNotExist
	}

	defer func() {
		_ = rows.Close()
	}()

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

	return ret, rows.Err()
}

// Database function that returns the count of photo uploaded by the user
func (db *appdbimpl) GetPhotosCount(id uint64) (int, error) {
	var count int
	if err := db.c.QueryRow("SELECT COUNT(*) FROM photos WHERE userId = ?", id).Scan(&count); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return count, ErrPhotoDoesNotExist
		}
	}
	return count, nil
}

// Check if a photo exists
func (db *appdbimpl) CheckPhoto(p Photo) (Photo, error) {
	var photo Photo
	if err := db.c.QueryRow(`SELECT Id, userId, photo, date FROM photos WHERE Id=?`, p.PhotoId).Scan(&photo.PhotoId, &photo.UserId, &photo.File, &photo.Date); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return photo, ErrPhotoDoesNotExist
		}
	}
	return photo, nil
}
