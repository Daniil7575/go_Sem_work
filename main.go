package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

var conn *pgx.Conn

const PASSWORD_HASH_COST = 14
const ACCESS_TOKEN_LIFETIME_H = 24
const JWT_SECRET_KEY = "SuperSecretJWTKey"

type User struct {
	ID           string `json:"id"`
	Username     string `json:"username"`
	PasswordHash string `json:"password_hash"`
}

type UserAuth struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserClaim struct {
	jwt.RegisteredClaims
	Sub       string
	UserName string
	Iat      time.Time
	Exp      time.Time
}

type AccessToken struct {
	Value string `json:"access_token"`
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
	token = jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": user.Username,
		"sub":      user.ID,
		"iat":      iat.Unix(),
		"exp":      exp.Unix(),
	})
	tokenString, err := token.SignedString([]byte(JWT_SECRET_KEY))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func _ValidateAccessToken(tokenString string) (string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &UserClaim{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(JWT_SECRET_KEY), nil
	})
	payload, ok := token.Claims.(*UserClaim)
	if !ok {
		return "", errors.New("bad token")
	}
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	if time.Now().Local().After(payload.ExpiresAt.Time) {
		return "", errors.New("expired token")
	}
	return payload.Subject, err
}

func main() {
	var err error
	conn, err = pgx.Connect(context.Background(), "postgres://postgres:123@localhost:5432/go_sem")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connection to database: %v\n", err)
		os.Exit(1)
	}
	r := mux.NewRouter()
	r.HandleFunc("/signin", signIn).Methods("POST")
	r.HandleFunc("/signup", signUp).Methods("POST")
	r.HandleFunc("/access-token", ValidateAccessToken).Methods("POST")
	http.Handle("/", r)
	fmt.Println("Сервер запущен на порту 8001")
	http.ListenAndServe(":8001", nil)
}

func signIn(w http.ResponseWriter, r *http.Request) {
	var u UserAuth
	var id string
	var password_hash string
	err := json.NewDecoder(r.Body).Decode((&u))
	if err != nil {
		fmt.Errorf("ERROR!")
		return
	}
	err = conn.QueryRow(context.Background(), "select id, password_hash from users where username = $1", u.Username).Scan(&id, &password_hash)
	if err != nil {
		fmt.Errorf("ERROR!")
		return
	}
	// db_pass_hash, err := CheckPasswordHash(u.Password, password_hash)
	if !CheckPasswordHash(u.Password, password_hash) {
		http.Error(w, "wrong password", http.StatusForbidden)
		return
	}
	access_token, _ := CreateAccessToken(User{ID: id, Username: u.Username, PasswordHash: password_hash})

	fmt.Fprintf(w, access_token)
}

func signUp(w http.ResponseWriter, r *http.Request) {
	var u UserAuth
	err := json.NewDecoder(r.Body).Decode((&u))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	password_hash, _ := HashPassword(u.Password)

	id := uuid.New()
	_, err = conn.Exec(context.Background(), "insert into users(id, username, password_hash) values($1, $2, $3)", id, u.Username, password_hash)
	if err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}
	access_token, _ := CreateAccessToken(User{ID: id.String(), Username: u.Username, PasswordHash: password_hash})
	fmt.Fprintf(w, access_token)
}

func ValidateAccessToken(w http.ResponseWriter, r *http.Request) {
	var access_token AccessToken
	err := json.NewDecoder(r.Body).Decode((&access_token))
	if err != nil{
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}
	if access_token.Value == "" {
		http.Error(w, "no token", http.StatusForbidden)
		return
	}
	user_id, err := _ValidateAccessToken(access_token.Value)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}
	fmt.Fprintf(w, user_id)
}
