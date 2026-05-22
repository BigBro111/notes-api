package model

type Note struct {
	ID     string   `json:"id"`
	Title  string   `json:"title"`
	Tags   []string `json:"tags"`
	Status string   `json:"status"`
}

type CreateNoteRequest struct {
	Title string   `json:"title"`
	Tags  []string `json:"tags"`
}
