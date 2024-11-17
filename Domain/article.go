package Domain

type Article struct {
	Title   string
	Content string
	Author  Author
}

type Author struct {
	Id   int64
	Name string
}
