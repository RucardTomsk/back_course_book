package service

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/RucardTomsk/course_book/model"
	"github.com/RucardTomsk/course_book/pkg/repository"
	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
	"gopkg.in/gomail.v2"
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
	user.Password = encryptString(user.Password)
	if flag {
		return s.repo.CreateUser(user)
	} else {
		return errors.New("The user exists")
	}
}

func (s *AuthService) GenerateToken(email string, password string) (string, string, error) {
	user, err := s.repo.GetUser(email, encryptString(password))
	if err != nil {
		return "", "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		user.Guid,
	})

	refreshToken, err := s.GenerateRefreshToken()
	if err != nil {
		return "", "", err
	}

	if err := s.IssueSessionUser(user, refreshToken); err != nil {
		return "", "", err
	}

	str_token, err := token.SignedString([]byte(signingKey))

	return str_token, refreshToken, err
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

func (s *AuthService) GenerateRefreshToken() (string, error) {
	b := make([]byte, 32)

	_s := rand.NewSource(time.Now().Unix())
	r := rand.New(_s)
	_, err := r.Read(b)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", b), nil
}

func (s *AuthService) IssueSessionUser(user model.User, refreshToken string) error {
	return s.repo.IssueSessionUser(user, refreshToken)
}
func (s *AuthService) GetUserToRefreshToken(refreshToken string) (model.User, error) {
	return s.repo.GetUserToRefreshToken(refreshToken)
}

func encryptString(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}

func (s *AuthService) CreateResetPassword(user model.User) error {
	code, err := s.repo.CreateResetPassword(user)
	if err != nil {
		return err
	}
	msg := gomail.NewMessage()
	msg.SetHeader("From", "course-book@tsu.ru")
	msg.SetHeader("To", user.Email)
	msg.SetHeader("Subject", "Восстановление аккаунта")
	msg.SetBody("text/html", fmt.Sprintf(`<p>Ваш код восстановление пароля для аккаунта "%s"</p>
										<p><strong>%s</strong></p>`, user.Email, code))

	if err := godotenv.Load(); err != nil {
		panic(err)
	}
	n := gomail.NewDialer("smtp.gmail.com", 465, "www.carat.ru@gmail.com", os.Getenv("PASS_EMAIL"))

	// Send the email
	if err := n.DialAndSend(msg); err != nil {
		panic(err)
	}

	return nil
}

func (s *AuthService) CheckResetPassword(code string, user model.User) error {
	return s.repo.CheckResetPassword(code, user)
}

func (s *AuthService) UserResetPassword(user model.User, newPassword string) error {
	return s.repo.UserResetPassword(user, encryptString(newPassword))
}

func (s *AuthService) GetUserByEmail(email string) (model.User, error) {
	return s.repo.GetUserByEmail(email)
}
