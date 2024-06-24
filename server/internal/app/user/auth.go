package user

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/gommon/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// JWTCustomClaims представляет пользовательские утверждения для JWT токена.
type JWTCustomClaims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

// newToken создает новый JWT токен для пользователя.
func (ua *Users) newToken(user User, tokenExpiresAt uint, SecretKeyForToken string) (string, error) {

	token := ua.getTokensWithClaims(user, tokenExpiresAt)

	tokenString, err := token.SignedString([]byte(SecretKeyForToken))
	if err != nil {
		log.Errorf("error in signedString access token. error: %v", err)
		return "", err
	}

	return tokenString, nil
}

// getTokensWithClaims создает токен с пользовательскими утверждениями.
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

// TokenAuthMiddlewareGRPS является промежуточным слоем аутентификации для gRPC.
func (ua *Users) TokenAuthMiddlewareGRPS(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

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

		if ua.Cfg.SecretKeyForToken != "" {

			userID, err := extractUserIDFromToken(token, ua.Cfg.SecretKeyForToken)
			if err != nil {
				return nil, status.Error(codes.Unauthenticated, "missing token")
			}

			ctx = context.WithValue(ctx, "userID", userID)

			return handler(ctx, req)
		}

		return nil, status.Error(codes.Unauthenticated, "")
	}

	return handler(ctx, req)
}

// TokenAuthMiddlewareREST является промежуточным слоем аутентификации для REST API.
func (ua *Users) TokenAuthMiddlewareREST(ctx context.Context, req *http.Request) metadata.MD {

	md := metadata.Pairs()

	if req.URL.Path == "/register" || req.URL.Path == "/login" {
		return md
	}

	token := req.Header.Get("Authorization")
	if token == "" {
		md.Append("error", "Unauthorized")
		return md
	}

	if ua.Cfg.SecretKeyForToken != "" {

		userID, err := extractUserIDFromToken(token, ua.Cfg.SecretKeyForToken)
		if err != nil {
			md.Append("error", "Unauthorized")
			return md
		}

		md.Append("userID", userID)
	}

	md.Append("error", "Unauthorized")
	return md
}

// extractUserIDFromToken извлекает идентификатор пользователя из JWT токена.
func extractUserIDFromToken(token, secretKeyForToken string) (string, error) {

	if len(token) == 0 {
		return "", fmt.Errorf("no token")
	}

	valid, userClaims, err := parseToken(token, secretKeyForToken)

	if err != nil {
		return "", err
	}

	if !valid {
		return "", fmt.Errorf("token protuh")
	}

	if userClaims.UserID == "" {
		return "", fmt.Errorf("no userID")
	}

	return userClaims.UserID, nil
}

// parseToken разбирает JWT токен и возвращает его утверждения.
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

// GetUserID извлекает идентификатор пользователя из контекста.
func GetUserID(ctx context.Context) (string, error) {

	md, _ := metadata.FromIncomingContext(ctx)

	if val, ok := md["error"]; ok && val[0] == "Unauthorized" {
		// Возвращаем статус 401 Unauthorized
		slog.Error("Unauthorized")
		return "", fmt.Errorf("Unauthorized")
	}

	var userID string

	userID, ok := ctx.Value("userID").(string)

	if !ok {
		slog.Error("userID not found in context")
		return "", fmt.Errorf("userID not found in context")
	}

	return userID, nil
}
