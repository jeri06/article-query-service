package model

type Article struct {
	ID      int64  `json:"id"`
	Author  string `json:"author"`
	Title   string `json:"title"`
	Body    string `json:"body"`
	Created string `json:"created"`
}

type GetManyArticleParams struct {
	Keyword string `validate:"-"`
	Author  string `validate:"-"`
	Page    int64  `validate:"min=1"`
	Size    int64  `validate:"min=1"`
}
