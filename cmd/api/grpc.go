package main

import (
	"context"
	"net"

	"github.com/Wrehat/ewallet-wallet/internal/config"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func ServeGRPC(ctx context.Context, cfg *config.AppConfig, log *zap.Logger) {

	grpcServer := grpc.NewServer()

	grpcListener, err := net.Listen("tcp", ":"+cfg.AppGrpcPort)

	if err != nil {
		log.Error("failed open port gRPC", zap.Error(err))
		panic(err)
	}

	go func() {
		<-ctx.Done()
		log.Info("gRpc server gracful shutdown")
		grpcServer.GracefulStop()
	}()

	if err := grpcServer.Serve(grpcListener); err != nil {
		log.Error("error running gRPC server", zap.Error(err))
	}

}
