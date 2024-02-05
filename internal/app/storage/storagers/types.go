package storagers

type ShortenedURL struct {
	UUID        int    `json:"user_id"`
	IdxShortURL string `json:"short_url"`
	OriginalURL string `json:"original_url"`
	DeletedFlag bool   `json:"is_deleted"`
}

type UsersURL struct {
	UserID int
	URLID  string
}
