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
	Title      string    `json:"title,omitempty" db:"title"`
	AuthorName string    `json:"authorName,omitempty" db:"username"`
	Content    string    `json:"content,omitempty" db:"content"`
	PostID     int64     `json:"postID,omitempty" db:"postID"`
	Score      int64     `json:"score,omitempty" db:"score"`
	CommID     int16     `json:"commID,omitempty" db:"commID"`
	Status     int16     `json:"status,omitempty" db:"status"`
	CreateTime time.Time `json:"createTime,omitempty" db:"create_time"`
}
