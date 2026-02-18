package middleware

import (
	"context"
	"net/http"
	"strings"
	"log"

	"github.com/golang-jwt/jwt/v5"

	"github.com/djalben/epn-killer-mvp/internal/utils"
)

// Контекстный ключ для хранения ID пользователя после проверки токена
type ContextKey string
const UserIDKey ContextKey = "userID"

// JWTAuthMiddleware - Middleware для проверки JWT-токена в заголовке Authorization
func JWTAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 1. Получение заголовка Authorization
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		// 2. Проверка формата: "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Authorization format must be Bearer <token>", http.StatusUnauthorized)
			return
		}

		tokenString := parts[1]

		// 3. Парсинг и валидация токена
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Проверка, что используется ожидаемый алгоритм подписи
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				log.Printf("Unexpected signing method: %v", token.Header["alg"])
				return nil, jwt.ErrSignatureInvalid
			}
			return utils.GetJWTSecret(), nil
		})
		
		if err != nil || !token.Valid {
			// Логируем ошибку, чтобы увидеть, почему токен невалиден
			log.Printf("Token validation failed: %v", err)
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		// 4. Извлечение ID пользователя из Claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "Invalid token claims", http.StatusUnauthorized)
			return
		}
		
		// Извлечение user_id. В нашем случае это float64, нужно преобразовать в int.
		userIDFloat, ok := claims["user_id"].(float64)
		if !ok {
			http.Error(w, "User ID not found in token claims", http.StatusUnauthorized)
			return
		}
		
		userID := int(userIDFloat)

		// 5. Передача ID пользователя в контекст запроса
		ctx := context.WithValue(r.Context(), UserIDKey, userID)
		r = r.WithContext(ctx)

		// Передаем управление следующему обработчику (вашему роуту)
		next.ServeHTTP(w, r)
	})
}