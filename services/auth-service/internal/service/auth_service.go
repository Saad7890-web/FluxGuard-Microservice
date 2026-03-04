package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/Saad7890-web/FluxGuard/internal/repository"
	"github.com/Saad7890-web/FluxGuard/internal/token"
)

type AuthService struct {
	repo   *repository.UserRepository
	tokens *token.Manager
}

func NewAuthService(repo *repository.UserRepository, tokens *token.Manager) *AuthService {
	return &AuthService{repo: repo, tokens: tokens}
}

func hashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}

func (s *AuthService) Register(ctx context.Context, email, password, deviceID string, accessTTL, refreshTTL int64) (string, string, int64, int64, error) {

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", "", 0, 0, err
	}

	userID, err := s.repo.Create(ctx, email, string(hash))
	if err != nil {
		return "", "", 0, 0, err
	}

	access, accessExp, err := s.tokens.Generate(userID, []string{"user"}, deviceID, accessTTL)
	if err != nil {
		return "", "", 0, 0, err
	}

	refresh, refreshExp, err := s.tokens.Generate(userID, []string{"user"}, deviceID, refreshTTL)
	if err != nil {
		return "", "", 0, 0, err
	}

	refreshHash := hashToken(refresh)
	if err := s.repo.StoreRefreshToken(ctx, userID, refreshHash, refreshExp); err != nil {
		return "", "", 0, 0, err
	}

	return access, refresh, accessExp, refreshExp, nil
}