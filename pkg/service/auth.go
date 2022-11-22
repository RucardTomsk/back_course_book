package service

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"time"

	"github.com/RucardTomsk/course_book/model"
	"github.com/RucardTomsk/course_book/pkg/repository"
	"github.com/golang-jwt/jwt"
)

const (
	salt       = "nsfgnstg45s5fbnsfdg"
	signingKey = "qwerqwerGS#jjsS"
)

type tokenClaims struct {
	jwt.StandardClaims
	UserGuid string `json:"user_guid"`
}
type AuthService struct {
	repo repository.Authorization
}

func NewAuthService(repo repository.Authorization) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) CreateUser(user model.User) error {
	flag, err := s.repo.CheckAbsentEmail(user.Email)
	if err != nil {
		return err
	}
	if flag {
		return s.repo.CreateUser(user)
	} else {
		return errors.New("The user exists")
	}
}

func (s *AuthService) GenerateToken(email string, password string) (string, error) {
	user, err := s.repo.GetUser(email, password)
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(12 * time.Hour).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		user.Guid,
	})

	return token.SignedString([]byte(signingKey))
}

func (s *AuthService) GetUserNotAccess(guid_node string) ([]model.User, error) {
	return s.repo.GetUserNotAccess(guid_node)
}

func (s *AuthService) GetUserFioByGuid(guid string) (string, error) {
	fio, err := s.repo.GetUserFIOByGuid(guid)
	if err != nil {
		return "", err
	}
	return fio, nil
}

func (s *AuthService) ParseToken(accessToken string) (string, error) {
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}

		return []byte(signingKey), nil
	})
	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return "", errors.New("token claims are not of type *tokenClaims")
	}

	return claims.UserGuid, nil
}

func encryptString(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}
