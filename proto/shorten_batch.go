package pb

import (
	context "context"

	"github.com/AlexTerra21/shortener/internal/app/auth"
	"github.com/AlexTerra21/shortener/internal/app/logger"
	"github.com/AlexTerra21/shortener/internal/app/models"
	"github.com/AlexTerra21/shortener/internal/app/utils"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// Пакетная запись данных
func (s *GRPCServer) ShortenBatch(ctx context.Context, in *ShortenBatchRequest) (*ShortenBatchResponse, error) {

	userID := GetUserIDFromMetadata(ctx)

	if userID == 0 {
		userID = (auth.GenerateUserID())
	}
	token, err := auth.BuildJWTString(userID)
	if err == nil {
		header := metadata.Pairs("Authorization", token)
		grpc.SendHeader(ctx, header)
	}

	response := &ShortenBatchResponse{
		Urls: []*ShortenBatchOut{},
	}
	batchStor := make([]models.BatchStore, 0)
	for _, v := range in.Urls {
		if v.OriginalUrl == "" {
			continue
		}
		id := utils.RandSeq(8)
		resp := &ShortenBatchOut{
			CorrelationId: v.CorrelationId,
			ShortUrl:      s.config.BaseURL + "/" + id,
		}
		batch := models.BatchStore{
			OriginalURL: v.OriginalUrl,
			IdxShortURL: id,
		}
		batchStor = append(batchStor, batch)
		response.Urls = append(response.Urls, resp)
	}
	if err := s.config.Storage.S.BatchSet(ctx, &batchStor, userID); err != nil {
		logger.Log().Debug("Error adding new url", zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}
	return response, nil
}
