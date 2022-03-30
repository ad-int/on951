package models

const TableArticles = "articles"

type ArticleBriefInfo struct {
	Id    uint `gorm:"primaryKey"`
	Title string
}

type ArticleWithComments struct {
	Id          uint `gorm:"primaryKey"`
	Title       string
	Description string
	Comments    []Comment `gorm:"foreignKey:ArticleId;references:Id"`
}

type Article struct {
	Id          uint `gorm:"primaryKey"`
	Title       string
	Description string
}
