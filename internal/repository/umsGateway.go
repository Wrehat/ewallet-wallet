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
	umsHost string
}

func NewUMSGateway(host string) domain.UMSGateway {
	return &umsGateway{
		umsHost: host,
	}
}

func (g *umsGateway) ValidateToken(ctx context.Context, token string) (*domain.TokenData, error) {
	conn, err := grpc.NewClient(g.umsHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to UMS: %v", err)
	}
	defer conn.Close()

	client := tokenvalidation.NewTokenValidationClient(conn)

	res, err := client.ValidateToken(ctx, &tokenvalidation.TokenRequest{
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
