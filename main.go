package main

import (
	"context"
	"os/user"
	// "crypto/ecdsa"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"authservice/auth"
	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/golang-jwt/jwt/v5"
)



// type User struct {
// 	Username string `json:"username"`
// 	Password string `json:"password"`
// }



type service struct {
	repo UserRepository
}

// type AuthService interface {
// 	SignIn(ctx context.Context, username, password string) (string, error)
// }

// type authService struct{}

// func (svc authService) SignIn(ctx context.Context, username, password string) (string, error) {
// 	var (
// 		token *jwt.Token
// 	)
// 	// Get if id from db
// 	user_id := "aaaaaa-aaaaa-aaaaaaa"
// 	if username == "example" && password == "password" {
// 		var iat time.Time = time.Now().Local()
// 		var exp time.Time = iat.Local().Add(
// 			time.Hour * time.Duration(ACCESS_TOKEN_LIFETIME_H),
// 		)
// 		token = jwt.NewWithClaims(jwt.SigningMethodES256, jwt.MapClaims{
// 			"username": username,
// 			"sub": user_id,
// 			"iat": iat,
// 			"exp": exp,
// 		})
// 		tokenString, err := token.SignedString(JWT_SECRET_KEY)
// 		if err != nil {
// 			return "", err
// 		}
// 		return tokenString, nil
// 	}
// 	return "", fmt.Errorf("invalid username or password")
// }

// func makeSignInEndpoint(svc AuthService) endpoint.Endpoint {
// 	return func(ctx context.Context, request interface{}) (interface{}, error) {
// 		req := request.(User)
// 		token, err := svc.SignIn(ctx, req.Username, req.Password)
// 		if err != nil {
// 			return nil, err
// 		}
// 		return map[string]string{"token": token}, nil
// 	}
// }

// func decodeSignInRequest(_ context.Context, r *http.Request) (interface{}, error) {
// 	var user User
// 	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
// 		return nil, err
// 	}
// 	return user, nil
// }

// func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
// 	return json.NewEncoder(w).Encode(response)
// }

// func main() {
// 	authService := authService{}

// 	signInHandler := httptransport.NewServer(
// 		makeSignInEndpoint(authService),
// 		decodeSignInRequest,
// 		encodeResponse,
// 		http.MethodPost,
// 	)

// 	http.Handle("/signin", signInHandler)

// 	fmt.Println("Сервер запущен на порту 8080")
// 	http.ListenAndServe(":8080", nil)
// }

// func main() {
// 	var iat time.Time = time.Now().Local()
// 	var exp time.Time = iat.Local().Add(time.Hour * time.Duration(ACCESS_TOKEN_LIFETIME_H))
// 	fmt.Println(iat)
// 	fmt.Println(exp)
// }
