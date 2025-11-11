package models

type PasteInput struct {
	Content  string `json:"content"`
	Language string `json:"language"`
	Password string `json:"password"`
}
type PasteOutput struct {
	ID       int    `json:"id"`
	Slug     string `json:"slug"`
	Content  string `json:"content"`
	Language string `json:"language"`
}
