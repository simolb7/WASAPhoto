package database

import (
	"database/sql"
	"errors"
)

// Database fuction that allows a user to ban another user
func (db *appdbimpl) InsertBan(b Ban) (Ban, error) {
	res, err := db.c.Exec("INSERT INTO bans (userId, bannedId) VALUES (?,?)", b.UserId, b.UserBannedId)
	if err != nil {
		return b, err
	}
	lastInsertID, err := res.LastInsertId()

	if err != nil {
		return b, err
	}
	b.BanId = uint64(lastInsertID)
	return b, nil

}

// Database fuction that removes a user ban
func (db *appdbimpl) RemoveBan(b Ban) error {
	res, err := db.c.Exec("DELETE FROM bans WHERE banId = ? AND userId = ? AND bannedId = ?", b.BanId, b.UserId, b.UserBannedId)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return ErrBanDoesNotExist
	}
	return nil
}

// Database fuction that return an user ban
func (db *appdbimpl) GetBan(u User, token uint64) (Ban, error) {
	var ban Ban
	if err := db.c.QueryRow(`SELECT * FROM bans WHERE bannedId = ? AND userId = ?`, u.Id, token).Scan(&ban.BanId, &ban.UserBannedId, &ban.UserId); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ban, ErrBanDoesNotExist
		}
	}
	return ban, nil
}
