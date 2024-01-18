package api

import (
	"github.com/simolb7/WASAPhoto/service/database"
)

type Profile struct {
	RequestId       uint64 `json:"requestId"`
	UserId          uint64 `json:"id"`
	Username        string `json:"username"`
	NumberFollowers int    `json:"NumberFollowers"`
	NumberFollowed  int    `json:"NumberFollowed"`
	PhotoCount      int    `json:"photoCount"`
}

type User struct {
	Id       uint64 `json:"id"`
	Username string `json:"username"`
}

type Follow struct {
	FollowId       uint64 `json:"followId"`
	UserFollowedId uint64 `json:"userFollowedId"`
	UserId         uint64 `json:"userId"`
}

type Ban struct {
	BanId    uint64 `json:"banId"`
	BannedId uint64 `json:"bannedId"`
	UserId   uint64 `json:"userId"`
}

type Photo struct {
	Id            uint64 `json:"id"`
	UserId        uint64 `json:"userId"`
	File          []byte `json:"file"`
	DateTime      string `json:"DateTime"`
	LikeNumber    int    `json:"likeNumber"`
	CommentNumber int    `json:"CommentNumber"`
}

type Like struct {
	LikeId     uint64 `json:"likeId"`
	UserId     uint64 `json:"identifier"`
	PhotoId    uint64 `json:"photoIdentifier"`
	PhotoOwner uint64 `json:"photoOwner"`
}

type Comment struct {
	CommentId    uint64 `json:"id"`
	UserId       uint64 `json:"userId"`
	PhotoId      uint64 `json:"photoId"`
	PhotoOwnerID uint64 `json:"photoOwnerID"`
	Content      string `json:"content"`
}

// Functions for talking from/to the db

// USER
// convert a user object from databse format to struct format
func (u *User) FromDatabase(user database.User) {
	u.Id = user.Id
	u.Username = user.Username
}

// convert a user object from struct format to databse format
func (u *User) ToDatabase() database.User {
	return database.User{
		Id:       u.Id,
		Username: u.Username,
	}
}

// FOLLOW
// convert a follow object from databse format to struct format
func (f *Follow) FollowFromDatabase(follow database.Follow) {
	f.FollowId = follow.Id
	f.UserFollowedId = follow.UserFollowedID
	f.UserId = follow.UserId
}

// convert a follow object from struct format to databse format
func (f *Follow) FollowToDatabase() database.Follow {
	return database.Follow{
		Id:             f.FollowId,
		UserId:         f.UserId,
		UserFollowedID: f.UserFollowedId,
	}
}

// BAN
// convert a ban object from databse format to struct format
func (b *Ban) BanFromDatabase(ban database.Ban) {
	b.BanId = ban.BanId
	b.BannedId = ban.UserBannedId
	b.UserId = ban.UserId
}

// convert a ban object from struct format to databse format
func (b *Ban) BanToDatabase() database.Ban {
	return database.Ban{
		BanId:        b.BanId,
		UserBannedId: b.BannedId,
		UserId:       b.UserId,
	}
}

// PHOTO
// convert a Photo object from databse format to struct format
func (p *Photo) PhotoFromDatabase(photo database.Photo) {
	p.Id = photo.PhotoId
	p.UserId = photo.UserId
	p.File = photo.File
	p.DateTime = photo.Date
	p.LikeNumber = photo.LikeNumber
	p.CommentNumber = photo.CommentNumber
}

// convert a Photo object from struct format to databse format
func (p *Photo) PhotoToDatabase() database.Photo {
	return database.Photo{
		PhotoId:       p.Id,
		UserId:        p.UserId,
		File:          p.File,
		Date:          p.DateTime,
		LikeNumber:    p.LikeNumber,
		CommentNumber: p.CommentNumber,
	}
}

// LIKE
// convert a Like object from databse format to struct format
func (l *Like) LikeFromDatabase(like database.Like) {
	l.LikeId = like.LikeId
	l.UserId = like.UserId
	l.PhotoId = like.PhotoId
	l.PhotoOwner = like.PhotoOwner

}

// convert a Like object from struct format to databse format
func (l *Like) LikeToDatabase() database.Like {
	return database.Like{
		LikeId:     l.LikeId,
		UserId:     l.UserId,
		PhotoId:    l.PhotoId,
		PhotoOwner: l.PhotoOwner,
	}
}

// COMMENT
// convert a comment object from databse format to struct format
func (c *Comment) CommentFromDatabase(comment database.Comment) {
	c.CommentId = comment.CommentId
	c.UserId = comment.UserId
	c.PhotoId = comment.PhotoId
	c.Content = comment.Content
}

// convert a comment object from struct format to databse format
func (c *Comment) CommentToDatabase() database.Comment {
	return database.Comment{
		CommentId:    c.CommentId,
		UserId:       c.UserId,
		PhotoId:      c.PhotoId,
		PhotoOwnerID: c.PhotoOwnerID,
		Content:      c.Content,
	}
}
