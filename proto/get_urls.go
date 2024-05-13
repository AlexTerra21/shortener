package pb

import (
	context "context"

	"github.com/AlexTerra21/shortener/internal/app/logger"
	"github.com/AlexTerra21/shortener/internal/app/storage/storagers"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Получение всех записей для переданного UserID
func (s *GRPCServer) GetUrls(ctx context.Context, in *Empty) (*GetUrlsResponse, error) {
	userID := GetUserIDFromMetadata(ctx)
	if userID == 0 {
		return nil, status.Error(codes.Unauthenticated, "")
	}

	db, ok := s.config.Storage.S.(*storagers.DB)
	if !ok {
		return nil, status.Error(codes.Internal, "Database not supported")
	}
	dbResponse, err := db.GetAll(ctx, s.config.BaseURL, userID)
	if err != nil {
		logger.Log().Debug("error get all urls from DB", zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())

	}
	if dbResponse == nil {
		logger.Log().Debug("Empty DB")
		return nil, status.Error(codes.NotFound, "Empty DB")
	}
	response := &GetUrlsResponse{
		Urls: []*UrlsInfo{},
	}
	for _, v := range dbResponse {
		url := &UrlsInfo{
			ShortUrl:    v.IdxShortURL,
			OriginalUrl: v.OriginalURL,
		}
		response.Urls = append(response.Urls, url)
	}

	return response, nil

}
