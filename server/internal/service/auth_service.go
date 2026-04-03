package service

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/lunancy1992/jianghu-server/internal/config"
	"github.com/lunancy1992/jianghu-server/internal/model"
	jwtpkg "github.com/lunancy1992/jianghu-server/internal/pkg/jwt"
	"github.com/lunancy1992/jianghu-server/internal/repo"
	"github.com/lunancy1992/jianghu-server/internal/sms"
)

type AuthUser struct {
	ID          int64  `json:"id"`
	Phone       string `json:"phone"`
	Nickname    string `json:"nickname"`
	Avatar      string `json:"avatar"`
	CoinBalance int64  `json:"coin_balance"`
	IsVerified  bool   `json:"is_verified"`
	Role        string `json:"role"`
}

type AuthService struct {
	userRepo *repo.UserRepo
	coinRepo *repo.CoinRepo
	cfg      *config.Config
	sms      sms.Provider
}

func NewAuthService(userRepo *repo.UserRepo, coinRepo *repo.CoinRepo, cfg *config.Config) *AuthService {
	var smsProvider sms.Provider
	if cfg.Auth.SMSProvider == "aliyun" && cfg.Auth.SMSAccessKey != "" {
		smsProvider = sms.NewAliyunSMS(cfg)
		log.Println("Using Aliyun SMS provider")
	} else {
		smsProvider = sms.NewStubSMS()
		log.Println("Using Stub SMS provider (development mode)")
	}

	return &AuthService{
		userRepo: userRepo,
		coinRepo: coinRepo,
		cfg:      cfg,
		sms:      smsProvider,
	}
}

// SendSMS sends a verification code.
func (s *AuthService) SendSMS(ctx context.Context, phone string) error {
	var code string

	// 开发环境：指定手机号使用固定验证码
	mockPhones := map[string]string{
		"18612293920":  "1234",
		"15811449974":  "1234",
		"13521935523":  "1234",
	}

	if mockCode, ok := mockPhones[phone]; ok {
		code = mockCode
		log.Printf("[SMS MOCK] Using fixed code %s for phone %s", code, phone)
	} else {
		code = fmt.Sprintf("%06d", rand.Intn(1000000))
	}

	expiresAt := time.Now().Add(5 * time.Minute)

	if err := s.userRepo.SaveVerification(ctx, phone, code, expiresAt); err != nil {
		return err
	}

	return s.sms.Send(ctx, phone, code)
}

// LoginWithSMS verifies code and returns JWT.
func (s *AuthService) LoginWithSMS(ctx context.Context, phone, code string) (string, *AuthUser, error) {
	v, err := s.userRepo.FindVerification(ctx, phone, code)
	if err != nil {
		return "", nil, err
	}
	if v == nil {
		return "", nil, fmt.Errorf("invalid or expired verification code")
	}

	if err := s.userRepo.MarkVerificationUsed(ctx, v.ID); err != nil {
		return "", nil, err
	}

	user, err := s.userRepo.FindByPhone(ctx, phone)
	if err != nil {
		return "", nil, err
	}

	if user == nil {
		// New user registration
		user = &model.User{
			Phone:    phone,
			Nickname: "江湖侠客" + phone[len(phone)-4:],
			Role:     "user",
			Status:   0,
		}
		id, err := s.userRepo.Create(ctx, user)
		if err != nil {
			return "", nil, err
		}

		// Re-fetch to get created_at etc.
		user, err = s.userRepo.FindByID(ctx, id)
		if err != nil {
			return "", nil, err
		}

		// Create coin account with initial grant
		if err := s.coinRepo.CreateAccount(ctx, id, 10); err != nil {
			return "", nil, err
		}
	}

	token, err := s.generateToken(user)
	if err != nil {
		return "", nil, err
	}

	profile, err := s.buildAuthUser(ctx, user)
	if err != nil {
		return "", nil, err
	}

	return token, profile, nil
}

// LoginWithOAuth handles third-party login.
func (s *AuthService) LoginWithOAuth(ctx context.Context, provider, providerID, unionID, nickname, avatar string) (string, *AuthUser, error) {
	oauth, err := s.userRepo.FindOAuth(ctx, provider, providerID)
	if err != nil {
		return "", nil, err
	}

	var user *model.User
	if oauth != nil {
		user, err = s.userRepo.FindByID(ctx, oauth.UserID)
		if err != nil {
			return "", nil, err
		}
	} else {
		user = &model.User{
			Nickname: nickname,
			Avatar:   avatar,
			Role:     "user",
			Status:   0,
		}
		id, err := s.userRepo.Create(ctx, user)
		if err != nil {
			return "", nil, err
		}
		user.ID = id

		_, err = s.userRepo.CreateOAuth(ctx, &model.UserOAuth{
			UserID:     id,
			Provider:   provider,
			ProviderID: providerID,
			UnionID:    unionID,
		})
		if err != nil {
			return "", nil, err
		}

		if err := s.coinRepo.CreateAccount(ctx, id, 10); err != nil {
			return "", nil, err
		}
	}

	token, err := s.generateToken(user)
	if err != nil {
		return "", nil, err
	}

	profile, err := s.buildAuthUser(ctx, user)
	if err != nil {
		return "", nil, err
	}

	return token, profile, nil
}

// RefreshToken generates a new token from an existing valid token.
func (s *AuthService) RefreshToken(ctx context.Context, tokenStr string) (string, error) {
	claims, err := jwtpkg.Validate(s.cfg.Auth.JWTSecret, tokenStr)
	if err != nil {
		return "", err
	}

	user, err := s.userRepo.FindByID(ctx, claims.UserID)
	if err != nil || user == nil {
		return "", fmt.Errorf("user not found")
	}

	return s.generateToken(user)
}

// GetProfile returns user info by ID.
func (s *AuthService) GetProfile(ctx context.Context, userID int64) (*AuthUser, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil || user == nil {
		return nil, err
	}
	return s.buildAuthUser(ctx, user)
}

func (s *AuthService) buildAuthUser(ctx context.Context, user *model.User) (*AuthUser, error) {
	account, err := s.coinRepo.GetAccount(ctx, user.ID)
	if err != nil {
		return nil, err
	}
	var balance int64
	if account != nil {
		balance = account.Balance
	}
	return &AuthUser{
		ID:          user.ID,
		Phone:       user.Phone,
		Nickname:    user.Nickname,
		Avatar:      user.Avatar,
		CoinBalance: balance,
		IsVerified:  false,
		Role:        user.Role,
	}, nil
}

func (s *AuthService) generateToken(user *model.User) (string, error) {
	return jwtpkg.Generate(
		s.cfg.Auth.JWTSecret,
		user.ID,
		user.Role,
		user.Nickname,
		s.cfg.Auth.JWTExpireHours,
	)
}
