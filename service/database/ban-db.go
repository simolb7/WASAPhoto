package database

import "database/sql"

// Database fuction that allows a user to ban another user
func (db *appdbimpl) InserBan(b Ban) (Ban, error) {
	_, err := db.c.Exec("INSERT INTO bans (id, userId, bannedId) VALUES (?,?,?)", b.BanId, b.UserId, b.UserBannedId)
	if err != nil {
		return b, err
	}
	return b, nil

}

// Database fuction that removes a user ban by another one (banner)
func (db *appdbimpl) RemoveBan(b Ban) error {
	res, err := db.c.Exec("DELETE FROM bans WHERE id = ? AND userId = ? AND bannedId = ?", b.BanId, b.UserId, b.UserBannedId)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	} else if affected == 0 {
		return ErrBanDoesNotExist
	}
	return nil
}

// Database fuction that return an user ban
func (db *appdbimpl) GetBans(u User, token uint64) (Ban, error) {
	var ban Ban
	if err := db.c.QueryRow(`SELECT banId, bannedId, userId FROM bans WHERE bannedId = ? AND userId = ?`, u.Id, token).Scan(&ban.BanId, &ban.UserBannedId, &ban.UserId); err != nil {
		if err == sql.ErrNoRows {
			return ban, ErrBanDoesNotExist
		}
	}
	return ban, nil
}

// Database fuction that checks if the requesting user was banned by another 'user'.
// 'true' if is banned, 'false' otherwise
func (db *appdbimpl) BannedUserCheck(target User, request User) (bool, error) {

	var cnt int
	err := db.c.QueryRow("SELECT COUNT(*) FROM bans WHERE bannedId = ? AND userId = ?", target.Id, request.Id).Scan(&cnt)

	if err != nil {
		// Count always returns a row thanks to COUNT(*), so this situation should not happen
		return true, err
	}

	// If the counter is 1 then the user was banned
	if cnt == 1 {
		return true, nil
	}
	return false, nil
}
