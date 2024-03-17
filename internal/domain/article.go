package domain

type Article struct {
	Id      int64
	Title   string
	Content string
	Author
	Status ArticleStatus
}

type Author struct {
	Id   int64
	Name string
}

type ArticleStatus uint8

const (
	ArticleStatusUnKnown     ArticleStatus = 0 // 未知default
	ArticleStatusUnPublished ArticleStatus = 1 // 未发布
	ArticleStatusPublished   ArticleStatus = 2 // 已发布
	ArticleStatusPrivate     ArticleStatus = 3 // 私密
)

func (a ArticleStatus) ToUint8() uint8 {
	return uint8(a)
}
