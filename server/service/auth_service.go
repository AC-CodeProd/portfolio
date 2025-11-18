package service

import (
	"portfolio/config"
	"portfolio/domain"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	cfg *config.JWTConfig
}

type tokenData struct {
	Token     string
	UserID    int
	ExpiresAt time.Time
	ExpiresIn int
}

func NewAuthService(cfg *config.JWTConfig) *AuthService {
	service := &AuthService{
		cfg: cfg,
	}

	return service
}

func (as *AuthService) HashPassword(password, salt string) (string, error) {
	if len(password+salt) > 72 {
		return "", domain.NewValidationError("Password and salt length exceeds 72 bytes (bcrypt limitation)", "password", nil)
	}
	bytes, err := bcrypt.GenerateFromPassword([]byte(password+salt), bcrypt.DefaultCost)
	return string(bytes), err
}

func (as *AuthService) CheckPassword(password, salt, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password+salt)) == nil
}

func (as *AuthService) GenerateToken(userID int) (*tokenData, error) {
	expDuration, _ := time.ParseDuration(as.cfg.Expiration)
	now := time.Now()
	expiresAt := now.Add(expDuration)

	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     expiresAt.Unix(),
		"iss":     as.cfg.Issuer,
		"aud":     as.cfg.Audience,
	}
	token := jwt.NewWithClaims(jwt.GetSigningMethod(as.cfg.SigningMethod), claims)
	tokenString, err := token.SignedString([]byte(as.cfg.Secret))
	if err != nil {
		return nil, err
	}

	tokenData := &tokenData{
		Token:     tokenString,
		UserID:    userID,
		ExpiresAt: expiresAt,
		ExpiresIn: int(expDuration.Seconds()),
	}

	return tokenData, nil
}
