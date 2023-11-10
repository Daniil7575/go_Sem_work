package implementation

import (
	usersvc "authservice/auth"
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/gofrs/uuid"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

const PASSWORD_HASH_COST = 14
const ACCESS_TOKEN_LIFETIME_H = 24
const JWT_SECRET_KEY = "SuperSecretJWTKey"

// It is a const value

type service struct {
	repo   usersvc.UserRepository
	logger log.Logger
}

type UserClaim struct {
	jwt.RegisteredClaims
	ID        int
	UserName  string
	ExpiresAt time.Time
}

func NewService(repo usersvc.UserRepository) usersvc.Service {
	return &service{
		repo: repo,
	}
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), PASSWORD_HASH_COST)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func CreateAccessToken(user usersvc.User) (string, error) {
	var (
		token *jwt.Token
	)
	var iat time.Time = time.Now().Local()
	var exp time.Time = iat.Local().Add(
		time.Hour * time.Duration(ACCESS_TOKEN_LIFETIME_H),
	)
	token = jwt.NewWithClaims(jwt.SigningMethodES256, jwt.MapClaims{
		"username": user.Username,
		"sub":      user.ID,
		"iat":      iat,
		"exp":      exp,
	})
	tokenString, err := token.SignedString([]byte(JWT_SECRET_KEY))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func ValidateAccessToken(tokenString string) (UserClaim, error) {
	var userClaim UserClaim
	_, err := jwt.ParseWithClaims(tokenString, &userClaim, func(token *jwt.Token) (interface{}, error) {
		return []byte(JWT_SECRET_KEY), nil
	})
	if userClaim.ExpiresAt.After(time.Now().Local()) || err != nil {
		return UserClaim{}, errors.New("Bad token")
	}
	return userClaim, nil
}

func (s *service) SignUp(ctx context.Context, username string, password string) (string, error) {
	var (
		user usersvc.User
		err  error
	)
	uuid, _ := uuid.NewV4()
	user.ID = uuid.String()
	user.Username = username
	user.PasswordHash, err = HashPassword(password)
	if err != nil {
		return "", err
	}

	if err := s.repo.CreateUser(ctx, user); err != nil {
		return "", errors.New("User already exists")
	}

	return CreateAccessToken(user)
}

func (s *service) SignIn(ctx context.Context, username string, password string) (string, error) {
	// logger := log.
	var (
		user usersvc.User
		err  error
	)
	user, err = s.repo.GetUserByUsername(ctx, username)
	if err != nil {
		return "", err
	}

	return CreateAccessToken(user)
}
