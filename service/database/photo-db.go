package database

import "database/sql"

// Database function to upload a photo
func (db *appdbimpl) InsertPhoto(p Photo) (Photo, error) {
	_, err := db.c.Exec("INSERT INTO photos (id, userId, photo, date) VALUES (?,?,?,?)", p.PhotoId, p.UserId, p.File, p.Date)
	if err != nil {
		return p, err
	}
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

	_, err2 := db.c.Exec("DELETE FROM photos WHERE photoId = ?", id)
	if err2 != nil {
		return err2
	}
	return nil
}

// Database function that returns each photo of a specific user
func (db *appdbimpl) GetPhotos(u User, token uint64) ([]Photo, error) {
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
			if err == sql.ErrNoRows {
				return nil, err
			}
		}
		if err := db.c.QueryRow(`SELECT COUNT(*) FROM comments WHERE photoId = ?`, token, b.PhotoId).Scan(&b.CommentNumber); err != nil {
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

// Database function that returns the count of photo uploaded by the user
func (db *appdbimpl) GetPhotosCount(id uint64) (int, error) {
	var count int
	if err := db.c.QueryRow("SELECT COUNT(*) FROM photos WHERE userId = id").Scan(&count); err != nil {
		if err == sql.ErrNoRows {
			return count, ErrUserDoesNotExist
		}
	}
	return count, nil
}
