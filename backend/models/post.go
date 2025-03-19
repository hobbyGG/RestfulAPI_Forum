package models

import (
	"time"
)

type Post struct {
	PostID    int64 `db:"postID"`
	AuthorUID int64 `db:"authorUID"`
	Score     int64 `db:"score"`
	Status    int16 `db:"status"`
	ParamCreatePost
	CreateTime time.Time `db:"create_time"`
	UpdateTime time.Time `db:"update_time"`
}

type PostPreview struct {
	Title      string `json:"title" db:"title"`
	AuthorName string `json:"authorName" db:"username"`
	Content    string `json:"content" db:"content"`
	PostID     int64  `json:"postID" db:"postID"`
	Score      int64  `json:"score" db:"score"`
	CommID     int16  `json:"commID" db:"commID"`
	Status     int16  `json:"status" db:"status"`
}
