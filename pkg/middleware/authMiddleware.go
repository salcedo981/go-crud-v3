package middleware

import (
	"net/http"
	"strings"

	v1 "github.com/FDSAP-Git-Org/hephaestus/helper/v1"
	"github.com/FDSAP-Git-Org/hephaestus/respcode"
	utils_v1 "github.com/FDSAP-Git-Org/hephaestus/utils/v1"
	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v4"
)

func AuthMiddleware(c fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return v1.JSONResponseWithError(c, respcode.ERR_CODE_401,
			"Authorization header missing", nil, http.StatusUnauthorized)
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenString == authHeader {
		return v1.JSONResponseWithError(c, respcode.ERR_CODE_401,
			"Invalid token format", nil, http.StatusUnauthorized)
	}

	// Parse and validate JWT
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(utils_v1.GetEnv("JWT_SECRET")), nil
	})

	if err != nil || !token.Valid {
		return v1.JSONResponseWithError(c, respcode.ERR_CODE_401,
			"Invalid or expired token", err, http.StatusUnauthorized)
	}

	// Extract claims and store user info in context
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		c.Locals("userId", claims["body"].(map[string]interface{})["userId"])
		c.Locals("email", claims["body"].(map[string]interface{})["email"])
	}

	return c.Next()
}
