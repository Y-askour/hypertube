package model

import "time"

type Torrent struct {
	ID              uint      `gorm:"column:id;primaryKey"`
	MovieID         uint      `gorm:"column:movie_id;not null"`
	InfoHash        string    `gorm:"column:info_hash;size:40;not null"`
	FilePath        string    `gorm:"column:file_path;not null"`
	DownloadedBytes int64     `gorm:"column:downloaded_bytes;default:0"`
	Completed       bool      `gorm:"column:completed;default:false"`
	CreatedAt       time.Time `gorm:"column:created_at;autoCreateTime"`

	Movie Movie `gorm:"foreignKey:MovieID"`
}

func (Torrent) TableName() string {
	return `"APP"."TORRENT"`
}
