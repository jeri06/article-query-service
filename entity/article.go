package entity

type Article struct {
	ID      int64  `json:"id"`
	Author  string `json:"author"`
	Title   string `json:"title"`
	Body    string `json:"body" `
	Created string `json:"created"`
}
