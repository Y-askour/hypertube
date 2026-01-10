package model

// ============================
// USER
// ============================
type User struct {
	ID                 uint   `gorm:"column:id;primaryKey"`
	Language           string `gorm:"column:LANGUAGE;size:10;default:'en'"`
	Email              string `gorm:"column:email;size:255;unique;not null"`
	HashedPassword     string `gorm:"column:hashed_password;size:255"`
	FirstName          string `gorm:"column:first_name;size:100;not null"`
	LastName           string `gorm:"column:last_name;size:100;not null"`
	Strategy           string `gorm:"column:strategy;size:30;not null"`
	ProfilePicturePath string `gorm:"column:profile_picture_path;size:255"`
}

func (User) TableName() string {
	return `"APP"."USER"`
}
