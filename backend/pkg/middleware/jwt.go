package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
	"github.com/makskozlov/linguist-hub/internal/constants"
	"github.com/makskozlov/linguist-hub/pkg/apierror"
	"github.com/makskozlov/linguist-hub/pkg/httputil"
)

func JWTAuth(secret string) fiber.Handler {
	return func(c fiber.Ctx) error {
		authHeader := c.Get(constants.AuthorizationHeader)
		if authHeader == "" {
			return httputil.RespondError(c, apierror.ErrUnauthorized)
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != constants.BearerPrefix {
			return httputil.RespondError(c, apierror.ErrUnauthorized)
		}

		tokenString := parts[1]
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, apierror.ErrUnauthorized
			}
			return []byte(secret), nil
		})

		if err != nil || !token.Valid {
			return httputil.RespondError(c, apierror.ErrUnauthorized)
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return httputil.RespondError(c, apierror.ErrUnauthorized)
		}

		userID, ok := claims[constants.JWTContextKey].(float64)
		if !ok {
			return httputil.RespondError(c, apierror.ErrUnauthorized)
		}

		c.Locals("userID", uint(userID))
		return c.Next()
	}
}
