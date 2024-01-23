package storagers

type ShortenedURL struct {
	UUID        string `json:"uuid"`
	IdxShortURL string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}
