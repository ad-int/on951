package models

type Comment struct {
	Id        uint `gorm:"primaryKey"`
	ArticleId uint
	Content   string
}
