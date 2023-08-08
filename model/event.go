package model

type Event struct {
	Id         int    `gorm:"primaryKey" json:"id"`
	UserId     int    `gorm:"foreignKey:Id" json:"user_id"`
	EventName  string `gorm:"type:varchar(100)" json:"event_name"`
	EventDate  string `gorm:"type:varchar" json:"event_date"`
	URL        string `gorm:"type:varchar(255)" json:"url"`
	IsVerified bool   `gorm:"type:boolean;not null;default:false" json:"is_verified"`
	Points     int    `gorm:"type:integer;not null;default:0" json:"points"`
}
