package models

const TableArticles = "articles"

type Article struct {
	Id          uint64
	Title       string
	Description string
}

type ArticleBriefInfo struct {
	Id    uint64
	Title string
}
