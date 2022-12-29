package requests

type BlogRequest struct {
	Title       string `json:"title"`
	Image       string `json:"image"`
	Description string `json:"description"`
	Tags        string `json:"tags"`
	Category    string `json:"category"`
}

type FileHandler struct {
	Filename string `json:"filename"`
	Size     string `json:"size"`
	Header   string `json:"header"`
}
