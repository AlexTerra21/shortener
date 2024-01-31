package storagers

type ShortenedURL struct {
	UUID        string `json:"uuid"`
	IdxShortURL string `json:"short_url"`
	OriginalURL string `json:"original_url"`
	UserID      int    `json:"user_id"`
}
