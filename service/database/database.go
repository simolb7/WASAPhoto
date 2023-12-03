/*
Package database is the middleware between the app database and the code. All data (de)serialization (save/load) from a
persistent database are handled here. Database specific logic should never escape this package.
Then you can initialize the AppDatabase and pass it to the api package.
*/
package database

import (
	"database/sql"
	"errors"
	"fmt"
)

var ErrUserDoesNotExist = errors.New("User does not exist")
var ErrPhotoDoesNotExist = errors.New("Photo does not exist")
var ErrBanDoesNotExist = errors.New("Ban does not exist")
var ErrFollowDoesNotExist = errors.New("Follow does not exist")
var ErrCommentDoesNotExist = errors.New("Comment does not exist")
var ErrLikeDoesNotExist = errors.New("Like does not exist")
var ErrUsernameAlreadyExists = errors.New("Username already exist")

type User struct {
	Id       uint64 `json:"identifier"`
	Username string `json:"username"`
}

/*
	type PhotoStream struct {
		Id            uint64 `json:"id"`
		UserID        uint64 `json:"userId"`
		File          []byte `json:"file"`
		Date          string `json:"date"`
		LikeNumber    int    `json:"likeNumber"`
		CommentNumber int    `json:"commentNumber"`
	}
*/

type Stream struct {
	Id     uint64  `json:"identifier"`
	Photos []Photo `json:"Photo"`
}

type Follow struct {
	Id             uint64 `json:"followId"`
	UserId         uint64 `json:"userId"`
	UserFollowedID uint64 `json:"followedId"`
}

type Ban struct {
	BanId        uint64 `json:"banId"`
	UserId       uint64 `json:"userId"`
	UserBannedId uint64 `json:"bannedId"`
}

type Photo struct {
	PhotoId       uint64 `json:"id"`
	UserId        uint64 `json:"userId"`
	File          []byte `json:"file"`
	Date          string `json:"date"`
	LikeNumber    int    `json:"likeNumber"`
	CommentNumber int    `json:"commentNumber"`
}

type Photos struct {
	RequestUser uint64  `json:"requestUser"`
	Identifier  uint64  `json:"identifier"`
	Photos      []Photo `json:"photos"`
}

type Like struct {
	PhotoOwner uint64 `json:"photoOwner"`
	PhotoId    uint64 `json:"photoIdentifier"`
	LikeId     uint64 `json:"likeId"`
	UserId     uint64 `json:"identifier"`
}

type Comment struct {
	PhotoOwner uint64 `json:"photoOwner"`
	PhotoId    uint64 `json:"photoId"`
	CommentId  uint64 `json:"id"`
	UserId     uint64 `json:"userId"`
	Content    string `json:"content"`
}

type Comments struct {
	RequestId  uint64    `json:"requestId"`
	PhotoId    uint64    `json:"photoId"`
	PhotoOwner uint64    `json:"identifier"`
	Comments   []Comment `json:"comments"`
}

// AppDatabase is the high level interface for the DB
type AppDatabase interface {
	CreateUser(User) (User, error)
	SetUsername(User, string) (User, error)
	GetUserId(string) (User, error)
	CheckUserById(User) (User, error)
	CheckUserByUsername(User) (User, error)
	CheckUser(User) (User, error)
	GetStream(u User) ([]Photo, error)
	//GetUserById(id int) (User, error)

	InsertPhoto(p Photo) (Photo, error)
	RemovePhoto(id uint64) error
	GetPhotos(u User, token uint64) ([]Photo, error)
	GetPhotosCount(id uint64) (int, error)
	CheckPhoto(p Photo) (Photo, error)

	InsertLike(l Like) (Like, error)
	RemoveLike(l Like) error
	RemoveLikes(user uint64, banned uint64) error
	GetLikeCount(photoid uint64) (int, error)
	GetLike(photoid uint64, token uint64) (Like, error)

	InsertFollow(f Follow) (Follow, error)
	RemoveFollow(FollowId uint64, UserId uint64, FollowedId uint64) error
	GetFollowersCount(id uint64) (int, error)
	GetFollowedCount(id uint64) (int, error)
	GetFollower(u User, token uint64) (Follow, error)

	InsertComment(c Comment) (Comment, error)
	RemoveComment(c Comment) error
	RemoveComments(user uint64, banned uint64) error
	GetCommentsCount(photoid uint64) (int, error)
	GetComments(photoid uint64) ([]Comment, error)

	InsertBan(b Ban) (Ban, error)
	RemoveBan(b Ban) error
	GetBan(u User, token uint64) (Ban, error)
	BannedUserCheck(target User, request User) (bool, error)

	Ping() error
}

type appdbimpl struct {
	c *sql.DB
}

// New returns a new instance of AppDatabase based on the SQLite connection `db`.
// `db` is required - an error will be returned if `db` is `nil`.
func New(db *sql.DB) (AppDatabase, error) {
	// Check if db is nil
	if db == nil {
		return nil, errors.New("database is required when building a AppDatabase")
	}
	// Enable foreign keys
	_, err := db.Exec("PRAGMA foreign_keys = ON")
	if err != nil {
		return nil, err
	}
	// Check if table exists. If not, the database is empty, and we need to create the structure
	var tableName string
	err = db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='users';`).Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) {
		usersDatabase := `CREATE TABLE users (
			Id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
			Username TEXT NOT NULL UNIQUE
			);`
		photosDatabase := `CREATE TABLE photos (
			Id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, 
			userId INTEGER NOT NULL,
			photo BLOB,
			date TEXT ,
			FOREIGN KEY (userId) REFERENCES users(Id)
			);`
		likesDatabase := `CREATE TABLE likes (
			Id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
			userId INTEGER NOT NULL,
			photoId INTEGER NOT NULL,
			photoOwner INTEGER NOT NULL,
			FOREIGN KEY (userId) REFERENCES users(Id),
			FOREIGN KEY (photoId) REFERENCES photos(Id)
			);`
		commentsDatabase := `CREATE TABLE comments (
			Id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
			userId INTEGER NOT NULL,
			photoId INTEGER NOT NULL,
			photoOwner INTEGER NOT NULL,
			content TEXT NOT NULL,
			FOREIGN KEY (userId) REFERENCES users(Id),
			FOREIGN KEY (photoId) REFERENCES photos(Id)
			);`
		bansDatabase := `CREATE TABLE bans (
			banId INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
			bannedId INTEGER NOT NULL,
			userId INTEGER NOT NULL,
			FOREIGN KEY (userId) REFERENCES users(Id)
			);`
		followersDatabase := `CREATE TABLE followers (
			Id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
			followerId INTEGER NOT NULL,
			userId INTEGER NOT NULL,
			FOREIGN KEY (userId) REFERENCES users(Id)
			);`
		_, err = db.Exec(usersDatabase)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure: %w", err)
		}
		_, err = db.Exec(photosDatabase)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure: %w", err)
		}
		_, err = db.Exec(likesDatabase)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure: %w", err)
		}
		_, err = db.Exec(commentsDatabase)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure: %w", err)
		}
		_, err = db.Exec(bansDatabase)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure: %w", err)
		}
		_, err = db.Exec(followersDatabase)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure: %w", err)
		}
	}

	return &appdbimpl{
		c: db,
	}, nil

}

func (db *appdbimpl) Ping() error {
	return db.c.Ping()
}
