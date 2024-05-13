package pb

import (
	context "context"
	"fmt"
	"net"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/AlexTerra21/shortener/internal/app/auth"
	"github.com/AlexTerra21/shortener/internal/app/config"
	"github.com/AlexTerra21/shortener/internal/app/logger"
)

// GRPCServer
type GRPCServer struct {
	UnimplementedShortenerServer
	server *grpc.Server
	config *config.Config
}

// Интерсептор для логирования
func logInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	logger.Log().Info("gRPC request",
		zap.String("method", info.FullMethod),
		zap.Any("request", req),
	)
	return handler(ctx, req)
}

// Интерсептор для авторизации
func authInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	var token string
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		values := md.Get("authorization")
		if len(values) > 0 {
			token = values[0]
		}
	}
	userID := auth.GetUserID(token)
	newCtx := ctx
	if userID > 0 {
		incomCtx, _ := metadata.FromIncomingContext(ctx)
		newMD := metadata.Pairs("userID", fmt.Sprintf("%d", userID))
		newCtx = metadata.NewIncomingContext(ctx, metadata.Join(incomCtx, newMD))
	}
	return handler(newCtx, req)
}

// Конструктор GRPCServer
func NewGRPCServer(config *config.Config) (*GRPCServer, error) {
	s := grpc.NewServer(grpc.ChainUnaryInterceptor(logInterceptor, authInterceptor))
	return &GRPCServer{
		server: s,
		config: config,
	}, nil
}

// Старт gRPC сервера
func (s *GRPCServer) Start() error {
	RegisterShortenerServer(s.server, s)

	addr := "localhost:3200"
	listen, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	logger.Log().Info("Running gRPC server", zap.String("address", addr))

	return s.server.Serve(listen)
}

// Остановка gRPC сервера
func (s *GRPCServer) Stop() error {
	logger.Log().Info("Stopping gRPC server")
	s.server.GracefulStop()
	return nil
}
