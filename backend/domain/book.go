package domain

type Book struct {
	ID       string `gorm:"primaryKey"`
	Title    string
	Author   string
	Provider Provider `gorm:"embedded"`
}
