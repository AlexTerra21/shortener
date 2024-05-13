package pb

import (
	context "context"
	"errors"
	"fmt"
	"strconv"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/AlexTerra21/shortener/internal/app/auth"
	"github.com/AlexTerra21/shortener/internal/app/errs"
	"github.com/AlexTerra21/shortener/internal/app/logger"
	"github.com/AlexTerra21/shortener/internal/app/storage/storagers"
	"github.com/AlexTerra21/shortener/internal/app/utils"
)

// Добавление одной записи
func (s *GRPCServer) Shorten(ctx context.Context, in *ShortenRequest) (*ShortenResponse, error) {

	var response ShortenResponse

	if in.Url == "" {
		return nil, status.Error(codes.InvalidArgument, "invalid url")
	}

	userID := GetUserIDFromMetadata(ctx)

	if userID == 0 {
		userID = (auth.GenerateUserId())
	}
	token, err := auth.BuildJWTString(userID)
	if err == nil {
		header := metadata.Pairs("Authorization", token)
		grpc.SendHeader(ctx, header)
	}

	id := utils.RandSeq(8)
	if err := s.config.Storage.S.Set(ctx, id, in.Url, userID); err != nil {
		logger.Log().Debug("Error adding new url", zap.Error(err))

		if errors.Is(err, errs.ErrConflict) {
			db, ok := s.config.Storage.S.(*storagers.DB)
			if ok {
				id, _ := db.GetShortURL(ctx, in.Url, userID)
				response.Result = s.config.BaseURL + "/" + id
			}
		} else {
			return nil, status.Error(codes.Internal, err.Error())

		}
	} else {
		response.Result = s.config.BaseURL + "/" + id
	}

	return &response, nil
}

// Получить UserID из метаданных
func GetUserIDFromMetadata(ctx context.Context) int {
	var userID int

	if md, ok := metadata.FromIncomingContext(ctx); ok {
		values := md.Get("userID")
		fmt.Printf("values = %v\n", values)

		if len(values) > 0 {
			userID, _ = strconv.Atoi(values[0])
		}
	}

	return userID
}
