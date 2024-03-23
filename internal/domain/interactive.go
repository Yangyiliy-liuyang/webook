package domain

type Interactive struct {
	BizId int64
	Biz   string

	ReadCnt    int64 // 阅读数
	LikeCnt    int64 //点赞数
	CollectCnt int64 //收藏数
	CommentCnt int64 //评论数
	ShareCnt   int64 // 分享数
	Liked      bool  `json:"liked"`
	Collected  bool  `json:"collected"`
	Ctime      int64
	Utime      int64
}
