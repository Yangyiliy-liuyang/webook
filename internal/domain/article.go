package domain

type Article struct {
	ID      int64
	Title   string
	Content string
	Author
}

type Author struct {
	ID   int64
	Name string
}
