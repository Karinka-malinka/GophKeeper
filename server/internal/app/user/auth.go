package user

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/gommon/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type JWTCustomClaims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

type LoginResponse struct {
	AccessToken  string
	RefreshToken string
}

func (ua *Users) newToken(user User, tokenExpiresAt uint, SecretKeyForToken string) (string, error) {

	token := ua.getTokensWithClaims(user, tokenExpiresAt)

	tokenString, err := token.SignedString([]byte(SecretKeyForToken))
	if err != nil {
		log.Errorf("error in signedString access token. error: %v", err)
		return "", err
	}

	return tokenString, nil
}

func (ua *Users) getTokensWithClaims(user User, tokenExpiresAt uint) (token *jwt.Token) {

	tokenClaims := &JWTCustomClaims{
		UserID: user.UUID.String(),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(tokenExpiresAt) * time.Minute)),
		},
	}

	token = jwt.NewWithClaims(jwt.SigningMethodHS256, tokenClaims)

	return token
}

func (ua *Users) GetToken(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

	if info.FullMethod != "/gophkeeper.UserService/Login" && info.FullMethod != "/gophkeeper.UserService/Register" {

		var token string

		md, ok := metadata.FromIncomingContext(ctx)

		if ok {
			values := md.Get("access_token")
			if len(values) > 0 {
				token = values[0]
			}
		}

		if len(token) == 0 {
			return nil, status.Error(codes.Unauthenticated, "missing token")
		}

		if ua.cfg.SecretKeyForToken != "" {

			valid, userClaims, err := parseToken(token, ua.cfg.SecretKeyForToken)

			if err != nil {
				return nil, err
			}

			if !valid {
				return nil, status.Errorf(codes.Unauthenticated, "Действие токена доступа истекло. Выполните команду LOGIN")
			}

			if userClaims.UserID == "" {
				return nil, status.Errorf(codes.Unauthenticated, "no userID")
			}

			md.Append("userID", userClaims.UserID)

			return handler(ctx, req)
		}

		return nil, status.Error(codes.Unauthenticated, "")
	}

	return handler(ctx, req)

}

func parseToken(tokenstr, secretKey string) (bool, *JWTCustomClaims, error) {

	token, err := jwt.ParseWithClaims(tokenstr, &JWTCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(secretKey), nil

	})

	if err != nil {
		if !errors.Is(err, jwt.ErrTokenExpired) {
			log.Infof("error in parsing token. error: %v", err)
			return false, nil, err
		}
	}

	userClaims := token.Claims.(*JWTCustomClaims)

	return token.Valid, userClaims, nil
}

func GetUserID(ctx context.Context) string {

	var userID string

	md, ok := metadata.FromIncomingContext(ctx)

	if ok {
		values := md.Get("userID")
		if len(values) > 0 {
			userID = values[0]
		}
	}

	return userID
}
