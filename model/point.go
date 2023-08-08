package model

type Point struct {
	Id     int `gorm:"primaryKey" json:"id"`
	UserId int `gorm:"foreignKey:Id" json:"user_id"`
	Amount int `gorm:"type:integer;not null;default:0" json:"amount"`
}
