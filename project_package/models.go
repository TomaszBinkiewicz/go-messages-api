package project_package

// Message struct
type Message struct {
	Id          int    `json:"id"`
	Email       string `json:"email"`
	Title       string `json:"title"`
	Content     string `json:"content"`
	MagicNumber int    `json:"magic_number"`
	Created     int    `json:"created"`
}

// SendTo struct
type SendTo struct {
	MagicNumber int `json:"magic_number"`
}
