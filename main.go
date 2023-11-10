package main

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

const PASSWORD_HASH_COST = 14
const ACCESS_TOKEN_LIFETIME_H = 24
const JWT_SECRET_KEY = "SuperSecretJWTKey"

type User struct {
	ID           string `json:"id"`
	Username     string `json:"username"`
	PasswordHash string `json:"password_hash"`
}

type UserClaim struct {
	jwt.RegisteredClaims
	ID        int
	UserName  string
	ExpiresAt time.Time
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), PASSWORD_HASH_COST)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func CreateAccessToken(user User) (string, error) {
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

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/signin", signIn).Methods("GET")
	http.Handle("/", r)
	fmt.Println("Сервер запущен на порту 8001")
	http.ListenAndServe(":8001", nil)
}

func signIn(w http.ResponseWriter, r *http.Request) {
	// get User by username from db
	var db_pass string = "12345"
	db_pass_hash, err := HashPassword(db_pass)
	if err != nil {
		fmt.Errorf("ERROR!")
	}
	plain_pass := r.URL.Query().Get("ppass")
	if !CheckPasswordHash(plain_pass, db_pass_hash) {
		fmt.Errorf("Bad creds")
	}

	username := r.URL.Query().Get("un")
	w.WriteHeader(http.StatusOK)

	fmt.Fprintf(w, "Hello!: %v\n", username)
}
