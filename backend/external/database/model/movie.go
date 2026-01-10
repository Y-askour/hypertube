package model

type Movie struct {
	ID       uint      `gorm:"column:id;primaryKey"`
	TMDBID   int       `gorm:"column:tmdb_id;unique;not null"`
	Torrents []Torrent `gorm:"foreignKey:MovieID"` // relation to torrents
}

func (Movie) TableName() string {
	return `"APP"."MOVIE"`
}
