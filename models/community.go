package models

import "time"

type Community struct {
	// 💡 重点：GORM 的 Scan 依赖 column 标签来匹配 SQL 结果
	ID   int64  `json:"id" gorm:"column:community_id"`
	Name string `json:"name" gorm:"column:community_name"`
}

type CommunityDetail struct {
	ID           int64     `json:"id" gorm:"column:community_id"`
	Name         string    `json:"name" gorm:"column:community_name"`
	Introduction string    `json:"introduction,omitempty" gorm:"column:introduction"`
	CreateTime   time.Time `json:"create_time" gorm:"column:create_time"`
}
