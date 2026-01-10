package model

import "time"

type MovieView struct {
	ID                  uint      `gorm:"column:id;primaryKey"`
	UserID              uint      `gorm:"column:user_id;not null"`
	MovieID             uint      `gorm:"column:movie_id;not null"`
	LastPositionSeconds int       `gorm:"column:last_position_seconds;default:0"`
	Completed           bool      `gorm:"column:completed;default:false"`
	UpdatedAt           time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (MovieView) TableName() string {
	return `"APP"."MOVIE_VIEW"`
}
