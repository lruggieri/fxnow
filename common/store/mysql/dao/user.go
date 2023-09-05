package dao

type User struct {
	ID        uint64 `gorm:"column:id"`
	UserID    string `gorm:"column:user_id"`
	FirstName string `gorm:"column:first_name"`
	LastName  string `gorm:"column:last_name"`
	Email     string `gorm:"column:email"`
}

func (*User) TableName() string {
	return "user"
}
