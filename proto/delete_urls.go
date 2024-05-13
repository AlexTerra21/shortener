package pb

import (
	context "context"

	"github.com/AlexTerra21/shortener/internal/app/storage/storagers"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Удаление записей для переданного UserID
func (s *GRPCServer) DeleteUrls(ctx context.Context, in *DeleteUrlsRequest) (*Empty, error) {
	userID := GetUserIDFromMetadata(ctx)
	if userID == 0 {
		return nil, status.Error(codes.Unauthenticated, "")
	}

	for _, urlID := range in.UrlsId {
		s.config.DelQueue.Push(storagers.UsersURL{
			UserID: userID,
			URLID:  urlID,
		})
	}

	return &Empty{}, nil
}
