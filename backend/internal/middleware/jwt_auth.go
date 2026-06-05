package middleware

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func JWTAuth(secret string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Get Authorization header
			auth := c.Request().Header.Get("Authorization")
			if auth == "" {
				return c.JSON(401, map[string]string{
					"error": "Missing authorization header",
				})
			}

			// Parse Bearer token
			if len(auth) < 7 || auth[:7] != "Bearer " {
				return c.JSON(401, map[string]string{
					"error": "Invalid authorization format",
				})
			}

			tokenString := auth[7:]

			// Parse JWT
			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}
				return []byte(secret), nil
			})

			if err != nil || !token.Valid {
				return c.JSON(401, map[string]string{
					"error": "Invalid token",
				})
			}

			// Extract claims
			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				return c.JSON(401, map[string]string{
					"error": "Invalid token claims",
				})
			}

			// Set user info in context
			if userID, ok := claims["user_id"].(string); ok {
				c.Set("user_id", userID)
			}

			if username, ok := claims["username"].(string); ok {
				c.Set("username", username)
			}

			if role, ok := claims["role"].(string); ok {
				c.Set("role", role)
			}

			return next(c)
		}
	}
}
