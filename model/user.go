package model

type User struct {
	Id       int    `gorm:"primaryKey" json:"id"`
	Role     string `gorm:"type:varchar(50)" json:"role"`
	Name     string `gorm:"type:varchar(50)" json:"name"`
	Email    string `gorm:"type:varchar(50)" json:"email"`
	Password string `gorm:"type:varchar(50)" json:"password"`
}
