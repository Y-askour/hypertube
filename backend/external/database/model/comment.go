package model

import "time"

type Comment struct {
	ID        uint      `gorm:"column:id;primaryKey"`
	MovieID   uint      `gorm:"column:movie_id;not null"`
	UserID    uint      `gorm:"column:user_id;not null"`
	Content   string    `gorm:"column:content;not null"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
}

func (Comment) TableName() string {
	return `"APP"."COMMENT"`
}
