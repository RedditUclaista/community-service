package middleware

import (
		"net/http"
	"strings"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v5"
)

type CustomClaims struct {
	jwt.RegisteredClaims
	Email     string
	Role      string
	SessionID string
	UserID    string
}

func JWTMiddleware(secretKey string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, "missing authorization header")
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid authorization header format")
			}

			tokenString := parts[1]
			
			token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
				return []byte(secretKey), nil
			})

			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid token")
			}

			if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
				c.Set("user_email", claims.Email)
				c.Set("user_role", claims.Role)
				return next(c)
			}
			return echo.NewHTTPError(http.StatusUnauthorized, "invalid token claims")
		}
	}
}

func RequireRole(roles ...string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			userRole := c.Get("user_role")
			if userRole == nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
			}
			
			roleStr := userRole.(string)
			for _, r := range roles {
				if roleStr == r {
					return next(c)
				}
			}
			return echo.NewHTTPError(http.StatusForbidden, "forbidden: insufficient permissions")
		}
	}
}
