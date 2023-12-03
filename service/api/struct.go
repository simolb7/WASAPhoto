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
	CommentId  uint64 `json:"id"`
	UserId     uint64 `json:"userId"`
	PhotoId    uint64 `json:"photoId"`
	PhotoOwner uint64 `json:"photoOwner"`
	Content    string `json:"content"`
}

//Functions for talking with the db

func (u *User) FromDatabase(user database.User) {
	u.Id = user.Id
	u.Username = user.Username
}

func (u *User) ToDatabase() database.User {
	return database.User{
		Id:       u.Id,
		Username: u.Username,
	}
}

func (f *Follow) FollowFromDatabase(follow database.Follow) {
	f.FollowId = follow.Id
	f.UserFollowedId = follow.UserFollowedID
	f.UserId = follow.UserId
}

func (f *Follow) FollowToDatabase() database.Follow {
	return database.Follow{
		Id:             f.FollowId,
		UserId:         f.UserId,
		UserFollowedID: f.UserFollowedId,
	}
}

func (b *Ban) BanFromDatabase(ban database.Ban) {
	b.BanId = ban.BanId
	b.BannedId = ban.UserBannedId
	b.UserId = ban.UserId
}

func (b *Ban) BanToDatabase() database.Ban {
	return database.Ban{
		BanId:        b.BanId,
		UserBannedId: b.BannedId,
		UserId:       b.UserId,
	}
}

func (p *Photo) PhotoFromDatabase(photo database.Photo) {
	p.Id = photo.PhotoId
	p.UserId = photo.UserId
	p.File = photo.File
	p.DateTime = photo.Date
	p.LikeNumber = photo.LikeNumber
	p.CommentNumber = photo.CommentNumber
}

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

func (l *Like) LikeFromDatabase(like database.Like) {
	l.LikeId = like.LikeId
	l.UserId = like.UserId
	l.PhotoId = like.PhotoId
	l.PhotoOwner = like.PhotoOwner

}

func (l *Like) LikeToDatabase() database.Like {
	return database.Like{
		LikeId:     l.LikeId,
		UserId:     l.UserId,
		PhotoId:    l.PhotoId,
		PhotoOwner: l.PhotoOwner,
	}
}

func (c *Comment) CommentFromDatabase(comment database.Comment) {
	c.CommentId = comment.CommentId
	c.UserId = comment.UserId
	c.PhotoId = comment.PhotoId
	c.Content = comment.Content
}

func (c *Comment) CommentToDatabase() database.Comment {
	return database.Comment{
		CommentId:  c.CommentId,
		UserId:     c.UserId,
		PhotoId:    c.PhotoId,
		PhotoOwner: c.PhotoOwner,
		Content:    c.Content,
	}
}
