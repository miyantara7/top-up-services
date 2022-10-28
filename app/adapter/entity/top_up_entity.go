package entity

type TopUp struct {
	BaseModel
	NoKartu string `gorm:"column:no_kartu;size:20"`
	UserId  string `gorm:"column:user_id;size:100"`
	Product string `gorm:"column:product;size:100"`
	Price   string `gorm:"column:price;size:100"`
	Fee     string `gorm:"column:fee;size:100"`
}

func (TopUp) TableName() string {
	return "t_top_up"
}
