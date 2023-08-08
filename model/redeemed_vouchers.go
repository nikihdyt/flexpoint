package model

type RedeemedVoucher struct {
	Id           int    `gorm:"primaryKey" json:"id"`
	UserID       int    `gorm:"foreignKey:Id" json:"user_id"`
	VoucherCode  string `gorm:"type:varchar(50)" json:"voucher_code"`
	ReedemedDate string `gorm:"type:varchar(50)" json:"reedemed_date"`
	PointUsed    int    `gorm:"type:integer;not null;default:0" json:"point_used"`
}
