package models

import "time"

type Post struct {
	ID          int64     `json:"id" gorm:"column:post_id"`
	AuthorID    int64     `json:"author_id" gorm:"column:author_id"`
	CommunityID int64     `json:"community_id" gorm:"column:community_id" binding:"required"`
	Status      int32     `json:"status" gorm:"column:status"`
	Title       string    `json:"title" gorm:"column:title" binding:"required"`
	Content     string    `json:"content" gorm:"column:content" binding:"required"`
	CreateTime  time.Time `json:"create_time" gorm:"column:create_time"`
}

type PostDetailVO struct {
	AuthorName string `json:"author_name" gorm:"column:author_name"`
	*Post
	*CommunityDetail `json:"community"`
}
