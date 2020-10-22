package models

import (
	"time"
)

type PostColor struct {
	Font string
	Back string
	Tag  string
}

type Post struct {
	ID       int
	AuthorID string `json:"-"`
	Author   User
	Tags     string
	Title    string
	Type     string
	Content  string
	Price    int
	Clr      PostColor `gorm:"embedded;embeddedPrefix:clr_"`
	UnitID   int       `json:"-"`
	Unit     Unit      `json:"-"`
	// TODO 컬럼 생성 방지
	IsFavorite int `gorm:"<-:false"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type IPostStore struct{}

var PostStore IPostStore

func (store IPostStore) AddPost(post Post) {
	db.Create(&post)
}

func (store IPostStore) GetPost(postID int) Post {
	var post Post
	db.Where("id = ?", postID).Preload("Author").Find(&post)
	return post
}

func (store IPostStore) GetPostListByUnitID(userID string, unitID int) []Post {
	var posts []Post
	db.Model(&Post{}).Raw("select posts.*, (f.post_id is NOT NULL) as is_favorite from posts LEFT OUTER JOIN favorites f ON f.user_id = ? and posts.id = f.post_id WHERE posts.unit_id = ? order by posts.created_at desc ", userID, unitID).Preload("Author").Find(&posts)
	//db.Order("created_at desc").Where("unit_id = ?", unitID).Preload("Author").Find(&posts)
	return posts
}

func (store IPostStore) GetFavorites(user User) []Post {
	var posts []Post
	db.Model(&Post{}).Raw("select posts.*, (f.post_id is NOT NULL) as is_favorite from posts INNER JOIN favorites f ON f.user_id = ? and posts.id = f.post_id order by posts.created_at desc ", user.ID).Preload("Author").Find(&posts)
	return posts
}

func (store IPostStore) AddFavorite(uuid string, postID int) {
	db.Exec("INSERT INTO favorites (user_id, post_id) VALUES (?, ?);", uuid, postID)
}

func (store IPostStore) DeleteFavorite(uuid string, postID int) {
	db.Exec("DELETE FROM favorites where user_id = ? and post_id = ?", uuid, postID)
}