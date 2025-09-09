package dto

type User struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Fullname string `json:"fullname"`
	Password string `json:"password"`
}

type Note struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Text        string `json:"text"`
	User        User   `json:"user"`
}
