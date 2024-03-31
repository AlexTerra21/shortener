package storagers

// Информация об URL
// Используется при сохранении и получении записи из базы
type ShortenedURL struct {
	UUID        int    `json:"user_id"`      // id пользователя
	IdxShortURL string `json:"short_url"`    // индекс сокращенного URL
	OriginalURL string `json:"original_url"` // оригинальный URL
	DeletedFlag bool   `json:"is_deleted"`   // флаг удаленной записи
}

// Информация об удаляемой записи
type UsersURL struct {
	UserID int    // id пользователя
	URLID  string // индекс сокращенного URL
}
