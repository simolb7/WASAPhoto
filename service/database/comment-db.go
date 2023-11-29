package database

import "database/sql"

// Database function that adds a comment of a user to a photo
func (db *appdbimpl) InsertComment(c Comment) (Comment, error) {
	_, err := db.c.Exec("INSERT INTO comments (id, userId, photoid, photoOwner, content) VALUES (?,?,?,?)", c.CommentId, c.UserId, c.PhotoId, c.Content)
	if err != nil {
		return c, err
	}
	return c, nil

}

// Database function that removes a comment of a user
func (db *appdbimpl) RemoveComment(c Comment) error {
	res, err := db.c.Exec("DELETE FROM comments WHERE id = ? AND photoId = ? AND userId = ?", c.CommentId, c.PhotoId, c.UserId)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	} else if affected == 0 {
		return ErrCommentDoesNotExist
	}
	return nil
}

// Database function that removes each comments of a user if it is banned
func (db *appdbimpl) RemoveComments(user uint64, banned uint64) error {
	_, err := db.c.Exec("DELETE FROM comments WHERE userId = ? AND photoOwner = ?", user, banned)
	if err != nil {
		return err
	}
	return nil
}

// Database function that returns the count of comments on aphoto
func (db *appdbimpl) GetCommentsCount(photoid uint64) (int, error) {
	var count int

	if err := db.c.QueryRow("SELECT COUNT(*) FROM comments WHERE photoId = ?", photoid).Scan(&count); err != nil {
		if err == sql.ErrNoRows {
			return count, ErrCommentDoesNotExist
		}
	}
	return count, nil
}
