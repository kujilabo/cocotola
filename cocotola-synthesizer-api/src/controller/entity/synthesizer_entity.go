package entity

type SynthesizeParameter struct {
	Lang2 string `json:"lang2" binding:"required,len=2"`
	Text  string `json:"text" binding:"required"`
}

type AudioResponse struct {
	ID      int    `json:"id"`
	Lang2   string `json:"lang2"`
	Text    string `json:"text"`
	Content string `json:"content"`
}
