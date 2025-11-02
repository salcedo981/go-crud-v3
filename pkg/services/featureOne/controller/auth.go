package ctrFeatureOne

import (
	"fmt"
	"go_template_v3/pkg/config"
	"go_template_v3/pkg/global/utils"
	mdlFeatureOne "go_template_v3/pkg/services/featureOne/model"
	"net/http"
	"net/smtp"
	"time"

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
func UpdateUser(c fiber.Ctx) error {
	// 1. Get user ID from JWT
	userId := utils.GetUserId(c)

	// 2. Parse request body
	var reqBody map[string]interface{}
	if err := c.Bind().Body(&reqBody); err != nil {
		return v1.JSONResponseWithError(c, respcode.ERR_CODE_400, "Invalid body", err, http.StatusBadRequest)
	}

	// 3. Add userId to the payload
	reqBody["userId"] = userId

	// 4. Execute the query
	return utils.ExecuteDBFunction(c, "SELECT update_user_v3($1)", reqBody)
}

func ForgotPassword(c fiber.Ctx) error {
	var req struct {
		Email *string `json:"email"`
	}
	if err := c.Bind().Body(&req); err != nil {
		return v1.JSONResponseWithError(c, respcode.ERR_CODE_400, "Invalid body", err, http.StatusBadRequest)
	}

	if req.Email == nil || !utils_v1.IsEmailValid(*req.Email) {
		return v1.JSONResponseWithError(c, respcode.ERR_CODE_400, "Valid email required", nil, http.StatusBadRequest)
	}

	// Check user exists and get user info
	var user struct {
		Id    int
		Name  string
		Email string
	}
	if err := config.DBConnList[0].Raw("SELECT id, name, email FROM users WHERE email = ? AND deleted_at IS NULL", *req.Email).Scan(&user).Error; err != nil || user.Id == 0 {
		return v1.JSONResponseWithData(c, respcode.SUC_CODE_200, "If email exists, reset link sent", nil, http.StatusOK)
	}

	// Generate and store token
	token := utils_v1.GenerateRandomStrings(32, []string{utils_v1.UpperString, utils_v1.LowerString, utils_v1.NumericString})
	tokenHash := utils_v1.HashDataSHA512(token)
	expiresAt := time.Now().Add(1 * time.Hour)

	// Invalidate old tokens and create new one
	config.DBConnList[0].Exec("UPDATE password_reset_tokens SET used_at = NOW() WHERE user_id = ? AND used_at IS NULL", user.Id)
	if err := config.DBConnList[0].Exec("INSERT INTO password_reset_tokens (user_id, token_hash, expires_at) VALUES (?, ?, ?)", user.Id, tokenHash, expiresAt).Error; err != nil {
		return v1.JSONResponseWithError(c, respcode.ERR_CODE_500, "Failed to create reset token", err, http.StatusInternalServerError)
	}

	// Send email with reset link
	go func() {
		if err := sendPasswordResetEmail(user.Email, user.Name, token); err != nil {
			// Log the error but don't fail the request for security reasons
			fmt.Printf("Failed to send reset email to %s: %v\n", user.Email, err)
		}
	}()

	// In development, return token for testing
	response := map[string]interface{}{"token": token}
	return v1.JSONResponseWithData(c, respcode.SUC_CODE_200, "Password reset initiated", response, http.StatusOK)
}

func VerifyResetToken(c fiber.Ctx) error {
	var req struct {
		Token *string `json:"token"`
	}
	if err := c.Bind().Body(&req); err != nil {
		return v1.JSONResponseWithError(c, respcode.ERR_CODE_400, "Invalid body", err, http.StatusBadRequest)
	}

	if req.Token == nil {
		return v1.JSONResponseWithError(c, respcode.ERR_CODE_400, "Token required", nil, http.StatusBadRequest)
	}

	// Verify token
	tokenHash := utils_v1.HashDataSHA512(*req.Token)
	var isValid int
	err := config.DBConnList[0].Raw(`
		SELECT 1 FROM password_reset_tokens prt 
		JOIN users u ON prt.user_id = u.id 
		WHERE prt.token_hash = ? AND prt.used_at IS NULL AND prt.expires_at > NOW() AND u.deleted_at IS NULL
	`, tokenHash).Scan(&isValid).Error

	if err != nil || isValid == 0 {
		return v1.JSONResponseWithError(c, respcode.ERR_CODE_400, "Invalid or expired token", nil, http.StatusBadRequest)
	}

	return v1.JSONResponseWithData(c, respcode.SUC_CODE_200, "Token is valid", nil, http.StatusOK)
}

func ResetPassword(c fiber.Ctx) error {
	var req struct {
		Token       *string `json:"token"`
		NewPassword *string `json:"newPassword"`
	}

	if err := c.Bind().Body(&req); err != nil {
		return v1.JSONResponseWithError(c, respcode.ERR_CODE_400, "Invalid body", err, http.StatusBadRequest)
	}

	if req.Token == nil || req.NewPassword == nil {
		return v1.JSONResponseWithError(c, respcode.ERR_CODE_400, "Token and new password required", nil, http.StatusBadRequest)
	}

	if !utils_v1.IsPasswordValid(*req.NewPassword) {
		return v1.JSONResponseWithError(c, respcode.ERR_CODE_400, "Password does not meet requirements", nil, http.StatusBadRequest)
	}

	// Verify token and get user ID
	tokenHash := utils_v1.HashDataSHA512(*req.Token)
	var tokenInfo struct{ TokenId, UserId int }

	err := config.DBConnList[0].Raw(`
		SELECT prt.id as token_id, prt.user_id 
		FROM password_reset_tokens prt 
		JOIN users u ON prt.user_id = u.id 
		WHERE prt.token_hash = ? AND prt.used_at IS NULL AND prt.expires_at > NOW() AND u.deleted_at IS NULL
	`, tokenHash).Scan(&tokenInfo).Error

	if err != nil || tokenInfo.TokenId == 0 {
		return v1.JSONResponseWithError(c, respcode.ERR_CODE_400, "Invalid or expired token", nil, http.StatusBadRequest)
	}

	// Update password and mark token used
	hashedPassword, _ := utils_v1.HashData(*req.NewPassword)
	config.DBConnList[0].Exec("UPDATE users SET password = ?, updated_at = NOW() WHERE id = ?", hashedPassword, tokenInfo.UserId)
	config.DBConnList[0].Exec("UPDATE password_reset_tokens SET used_at = NOW() WHERE id = ?", tokenInfo.TokenId)

	return v1.JSONResponseWithData(c, respcode.SUC_CODE_200, "Password reset successfully", nil, http.StatusOK)
}

// FORGOT PASSWORD SEND MAIL FUNCTIONS
func sendPasswordResetEmail(email, name, token string) error {
	// Get frontend URL from environment variables
	frontendURL := utils_v1.GetEnv("FRONTEND_URL")
	if frontendURL == "" {
		frontendURL = "http://localhost:3000" // default for development
	}

	// Create reset link with token as query parameter
	resetLink := fmt.Sprintf("%s/reset-password?token=%s", frontendURL, token)

	// Email content
	subject := "Password Reset Request"

	// HTML email template
	htmlBody := fmt.Sprintf(`
	<!DOCTYPE html>
	<html>
	<head>
		<style>
			body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
			.container { max-width: 600px; margin: 0 auto; padding: 20px; }
			.button { display: inline-block; padding: 12px 24px; background-color: #007bff; 
					color: white !important; text-decoration: none; border-radius: 4px; margin: 20px 0; }
			.footer { margin-top: 30px; font-size: 12px; color: #666; }
		</style>
	</head>
	<body>
		<div class="container">
			<h2>Password Reset Request</h2>
			<p>Hello %s,</p>
			<p>You requested to reset your password. Click the button below to create a new password:</p>
			<p><a href="%s" class="button">Reset Password</a></p>
			<p>Or copy and paste this link in your browser:</p>
			<p><code>%s</code></p>
			<p>This link will expire in 1 hour for security reasons.</p>
			<p>If you didn't request this reset, please ignore this email.</p>
			<div class="footer">
				<p>This is an automated message, please do not reply to this email.</p>
			</div>
		</div>
	</body>
	</html>
	`, name, resetLink, resetLink)

	fmt.Printf("=== PASSWORD RESET EMAIL ===\n")
	fmt.Printf("To: %s\n", email)
	fmt.Printf("Subject: %s\n", subject)
	fmt.Printf("Reset Link: %s\n", resetLink)
	fmt.Printf("=======================\n")

	// Send via SMTP
	return sendWithSMTP(email, subject, htmlBody)
}

func sendWithSMTP(to, subject, htmlContent string) error {
	smtpHost := utils_v1.GetEnv("SMTP_HOST")
	smtpPort := utils_v1.GetEnv("SMTP_PORT")
	smtpUser := utils_v1.GetEnv("SMTP_USER")
	smtpPass := utils_v1.GetEnv("SMTP_PASS")
	from := utils_v1.GetEnv("EMAIL_FROM")

	auth := smtp.PlainAuth("", smtpUser, smtpPass, smtpHost)

	msg := []byte("To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"MIME-Version: 1.0\r\n" +
		"Content-Type: text/html; charset=UTF-8\r\n" +
		"\r\n" + htmlContent)

	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, msg)
	return err
}
