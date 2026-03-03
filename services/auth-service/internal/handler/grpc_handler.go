package handler

import (
	"context"

	"github.com/Saad7890-web/FluxGuard/internal/service"
	authv1 "github.com/Saad7890-web/FluxGuard/proto/auth/v1"
)

type GRPCHandler struct {
	authv1.UnimplementedAuthServiceServer
	svc *service.AuthService
}

func NewGRPCHandler(svc *service.AuthService) *GRPCHandler {
	return &GRPCHandler{svc: svc}
}

func (h *GRPCHandler) Register(ctx context.Context, req *authv1.RegisterRequest) (*authv1.AuthResponse, error) {

	access, refresh, accessExp, refreshExp, err :=
		h.svc.Register(ctx, req.Email, req.Password, req.DeviceId, 900, 604800)

	if err != nil {
		return nil, err
	}

	return &authv1.AuthResponse{
		AccessToken:           access,
		RefreshToken:          refresh,
		AccessTokenExpiresAt:  accessExp,
		RefreshTokenExpiresAt: refreshExp,
	}, nil
}