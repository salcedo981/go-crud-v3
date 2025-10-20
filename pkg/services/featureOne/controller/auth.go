package ctrFeatureOne

import (
	"go_template_v3/pkg/config"
	mdlFeatureOne "go_template_v3/pkg/services/featureOne/model"
	"net/http"

	v1 "github.com/FDSAP-Git-Org/hephaestus/helper/v1"
	"github.com/FDSAP-Git-Org/hephaestus/respcode"
	utils_v1 "github.com/FDSAP-Git-Org/hephaestus/utils/v1"
	"github.com/gofiber/fiber/v3"
)

func Register(c fiber.Ctx) error {
	var req mdlFeatureOne.RegisterRequest
	if err := c.Bind().Body(&req); err != nil {
		return v1.JSONResponseWithError(c, respcode.ERR_CODE_400,
			"Invalid body", err, http.StatusBadRequest)
	}

	// Validate email
	if !utils_v1.IsEmailValid(*req.Email) {
		return v1.JSONResponseWithError(c, respcode.ERR_CODE_400,
			"Invalid email format", nil, http.StatusBadRequest)
	}

	// Validate password
	if !utils_v1.IsPasswordValid(*req.Password) {
		return v1.JSONResponseWithError(c, respcode.ERR_CODE_400,
			"Password must be 8+ chars with uppercase, lowercase, and special char",
			nil, http.StatusBadRequest)
	}

	// Perform a quick raw query to check if the ID exists.
	var exists int
	existenceQuery := "SELECT 1 FROM users WHERE email = ? AND deleted_at IS NULL LIMIT 1"
	err := config.DBConnList[0].Raw(existenceQuery, *req.Email).Scan(&exists).Error

	if err != nil {
		return v1.JSONResponseWithError(c, respcode.ERR_CODE_500, "DB error during existence check", err, http.StatusInternalServerError)
	}
	if exists == 1 {
		return v1.JSONResponseWithError(c, respcode.ERR_CODE_404, "Email already exists", nil, http.StatusNotFound)
	}

	// Hash password
	hashedPassword, err := utils_v1.HashData(*req.Password)
	if err != nil {
		return v1.JSONResponseWithError(c, respcode.ERR_CODE_500,
			"Failed to hash password", err, http.StatusInternalServerError)
	}

	// Call DB function (create register_user function in PostgreSQL)
	var resultStr string
	err = config.DBConnList[0].Raw(
		"SELECT register_user(?, ?, ?)",
		req.Email,
		hashedPassword,
		req.Name,
	).Scan(&resultStr).Error

	if err != nil {
		return v1.JSONResponseWithError(c, respcode.ERR_CODE_500,
			"Registration failed", err, http.StatusInternalServerError)
	}

	return v1.JSONResponseWithData(c, respcode.SUC_CODE_201,
		"User registered successfully", nil, http.StatusCreated)
}

func Login(c fiber.Ctx) error {
	var req mdlFeatureOne.LoginRequest
	if err := c.Bind().Body(&req); err != nil {
		return v1.JSONResponseWithError(c, respcode.ERR_CODE_400,
			"Invalid body", err, http.StatusBadRequest)
	}

	// Validate request
	if req.Email == nil || *req.Email == "" {
		return v1.JSONResponseWithError(c, respcode.ERR_CODE_400,
			"Email is required", nil, http.StatusBadRequest)
	}
	if req.Password == nil || *req.Password == "" {
		return v1.JSONResponseWithError(c, respcode.ERR_CODE_400,
			"Password is required", nil, http.StatusBadRequest)
	}

	// Perform a quick raw query to check if the email exists
	var exists int
	existenceQuery := "SELECT 1 FROM users WHERE email = ? AND deleted_at IS NULL LIMIT 1"
	err := config.DBConnList[0].Raw(existenceQuery, *req.Email).Scan(&exists).Error

	if err != nil {
		return v1.JSONResponseWithError(c, respcode.ERR_CODE_500, "DB error during existence check", err, http.StatusInternalServerError)
	}
	if exists == 0 {
		return v1.JSONResponseWithError(c, respcode.ERR_CODE_401, "Invalid Credentials", nil, http.StatusNotFound)
	}

	var userResult mdlFeatureOne.User
	getUserQuery := "SELECT * FROM users WHERE email = ?"
	err = config.DBConnList[0].Raw(getUserQuery, *req.Email).Scan(&userResult).Error

	// Check if user was found
	if err != nil {
		return v1.JSONResponseWithError(c, respcode.ERR_CODE_500, respcode.ERR_CODE_500_MSG, nil, http.StatusInternalServerError)
	}

	// Verify password
	if !utils_v1.CheckHashData(*req.Password, *userResult.Password) {
		return v1.JSONResponseWithError(c, respcode.ERR_CODE_401,
			"Invalid credentials", nil, http.StatusUnauthorized)
	}

	// Generate JWT token
	claims := map[string]interface{}{
		"userId": *userResult.Id,
		"email":  *userResult.Email,
		"name":   *userResult.Name,
	}

	token, err := utils_v1.GenerateJWTSignedString(
		[]byte(utils_v1.GetEnv("JWT_SECRET")),
		24, // 24 hours
		claims,
	)

	if err != nil {
		return v1.JSONResponseWithError(c, respcode.ERR_CODE_500,
			"Token generation failed", err, http.StatusInternalServerError)
	}

	return v1.JSONResponseWithData(c, respcode.SUC_CODE_200,
		"Login successful",
		map[string]string{"token": token},
		http.StatusOK)
}

func Logout(c fiber.Ctx) error {
	// With JWT, logout is handled client-side by removing the token
	// Optionally implement token blacklisting here
	return v1.JSONResponseWithData(c, respcode.SUC_CODE_200,
		"Logout successful", nil, http.StatusOK)
}
