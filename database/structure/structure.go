package structure

const TableArticles = "articles"

type User struct {
	Id       uint
	Name     string
	Password string
	Articles []Article `gorm:"foreignKey:AuthorId; references:Id"`
	Comments []Comment `gorm:"foreignKey:UserId; references:Id"`
}

type ArticleBriefInfo struct {
	Id    uint `gorm:"primaryKey"`
	Title string
}

type ArticleWithComments struct {
	Id          uint `gorm:"primaryKey"`
	AuthorId    uint
	Title       string
	Description string
	Comments    []Comment `gorm:"foreignKey:ArticleId;references:Id"`
}

type Article struct {
	Id          uint `gorm:"primaryKey"`
	AuthorId    uint
	Title       string
	Description string
}

type Comment struct {
	Id        uint `gorm:"primaryKey"`
	UserId    uint
	ArticleId uint
	Content   string
}
