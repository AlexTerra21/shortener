package models

// Модель запроса для хэндлера shortenURL
type Request struct {
	URL string `json:"url"`
}

// Модель ответа для хэндлера shortenURL
type Response struct {
	Result string `json:"result"`
}

// Модель запроса для хэндлера batch
type BatchReq struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

// Модель ответа для хэндлера batch
type BatchResp struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

// Модель для пакетной записи в базу данных
type BatchStore struct {
	IdxShortURL string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}
