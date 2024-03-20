package domain

type Article struct {
	Id      int64  `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
	Author
	Status ArticleStatus `json:"status"`
	Ctime  int64         `json:"ctime"`
	Utime  int64         `json:"utime"`
	//Ctime *timestamppb.Timestamp `json:"ctime"`
	//Utime *timestamppb.Timestamp `json:"utime"`
}

type Author struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

func (a Article) Abstract() string {
	if len(a.Content) > 1024 {
		a.Content = a.Content[:1024]
	}
	return a.Content
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
