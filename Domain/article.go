package Domain

type Article struct {
	Id      int64
	Title   string
	Content string
	Author  Author
	Status  ArticleStatus
}

type Author struct {
	Id   int64
	Name string
}

type ArticleStatus uint8

func (a ArticleStatus) ToUint8() uint8 {
	return uint8(a)
}

const (
	// ArticleStatusUnKnow 未知状态，无意义
	ArticleStatusUnKnow = iota
	// ArticleStatusUnpublished 未发布
	ArticleStatusUnpublished
	// ArticleStatusPublished 已发布
	ArticleStatusPublished
	// ArticleStatusPrivate 隐藏
	ArticleStatusPrivate
)
