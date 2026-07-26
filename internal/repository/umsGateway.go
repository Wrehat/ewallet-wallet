package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/Wrehat/ewallet-wallet/internal/domain"
	"github.com/Wrehat/ewallet-wallet/pkg/tokenvalidation"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type umsGateway struct {
	client tokenvalidation.TokenValidationClient
}

func NewUMSGateway(host string) (domain.UMSGateway, error) {
	conn, err := grpc.NewClient(host, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect UMS gRPC: %w", err)
	}
	return &umsGateway{
		client: tokenvalidation.NewTokenValidationClient(conn),
	}, nil
}

func (g *umsGateway) ValidateToken(ctx context.Context, token string) (*domain.TokenData, error) {

	res, err := g.client.ValidateToken(ctx, &tokenvalidation.TokenRequest{
		Token: token,
	})
	if err != nil {
		return nil, err
	}

	if res.Message != "success" {
		return nil, errors.New(res.Message)
	}

	return &domain.TokenData{
		UserID:   int(res.Data.UserId),
		Username: res.Data.Username,
		FullName: res.Data.FullName,
		Email:    res.Data.Email,
	}, nil

}
