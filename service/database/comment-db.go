package database

import (
	"database/sql"
	"errors"
)

// Database function that adds a comment of a user to a photo
func (db *appdbimpl) InsertComment(c Comment) (Comment, error) {

	res, err := db.c.Exec("INSERT INTO comments (userId, photoid, photoOwnerID, content) VALUES (?,?,?,?)", c.UserId, c.PhotoId, c.PhotoOwnerID, c.Content)
	if err != nil {
		return c, err
	}

	lastInsertID, err := res.LastInsertId()

	if err != nil {
		return c, err
	}
	c.CommentId = uint64(lastInsertID)
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
	}
	if affected == 0 {
		return ErrCommentDoesNotExist
	}
	return nil
}

// Database function that removes each comments of a user due to a ban
func (db *appdbimpl) RemoveComments(user uint64, banned uint64) error {
	_, err := db.c.Exec("DELETE FROM comments WHERE userId = ? AND photoOwnerID = ?", banned, user)
	if err != nil {
		return err
	}
	_, err = db.c.Exec("DELETE FROM comments WHERE userId = ? AND photoOwnerID = ?", user, banned)
	if err != nil {
		return err
	}
	return nil
}

// Database function that returns the count of comments on a photo
func (db *appdbimpl) GetCommentsCount(photoid uint64) (int, error) {
	var count int

	err := db.c.QueryRow("SELECT COUNT(*) FROM comments WHERE photoId = ?", photoid).Scan(&count)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrCommentDoesNotExist
		}
		return count, err
	}
	return count, nil
}

// return all comments from a photo
func (db *appdbimpl) GetComments(photoid uint64) ([]Comment, error) {
	var comments []Comment
	rows, err := db.c.Query(`SELECT Id, userId, photoId, photoOwnerID, content FROM comments WHERE photoId = ?`, photoid)
	if err != nil {
		return comments, ErrPhotoDoesNotExist
	}
	defer func() { _ = rows.Close() }()
	for rows.Next() {
		var comment Comment
		err = rows.Scan(&comment.CommentId, &comment.UserId, &comment.PhotoId, &comment.PhotoOwnerID, &comment.Content)
		if err != nil {
			return nil, ErrCommentDoesNotExist
		}
		comments = append(comments, comment)
	}
	if err := rows.Err(); err != nil {
		return nil, ErrCommentDoesNotExist
	}
	return comments, nil
}
