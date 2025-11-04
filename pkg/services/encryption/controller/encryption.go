package ctrEncryption

import (
	mdlEncryption "go_template_v3/pkg/services/encryption/model"
	"log"
	"net/http"

	"github.com/FDSAP-Git-Org/hephaestus/encryption"
	v1 "github.com/FDSAP-Git-Org/hephaestus/helper/v1"
	"github.com/FDSAP-Git-Org/hephaestus/respcode"
	"github.com/gofiber/fiber/v3"
)

func EncryptDBCredentials(c fiber.Ctx) error {
	// 1. Parse request body
	var req mdlEncryption.EncryptCredentialsRequest
	if err := c.Bind().Body(&req); err != nil {
		return v1.JSONResponseWithError(
			c, 
			respcode.ERR_CODE_400, 
			"Invalid request body", 
			err, 
			http.StatusBadRequest,
		)
	}

	// 2. Validate required fields
	if req.SecretKey == "" {
		return v1.JSONResponseWithError(
			c, 
			respcode.ERR_CODE_400, 
			"Secret key is required", 
			nil, 
			http.StatusBadRequest,
		)
	}

	if len(req.SecretKey) != 32 {
		return v1.JSONResponseWithError(
			c, 
			respcode.ERR_CODE_400, 
			"Secret key must be exactly 32 characters long", 
			nil, 
			http.StatusBadRequest,
		)
	}

	// 3. Encrypt each field
	encryptedData := mdlEncryption.EncryptedData{}

	// Encrypt host
	if req.Host != "" {
		encryptedHost, err := encryption.Encrypt(req.Host, req.SecretKey)
		if err != nil {
			log.Printf("Error encrypting host: %v", err)
			return v1.JSONResponseWithError(
				c, 
				respcode.ERR_CODE_500, 
				"Failed to encrypt host", 
				err, 
				http.StatusInternalServerError,
			)
		}
		encryptedData.Host = encryptedHost
	}

	// Encrypt database name
	if req.DBName != "" {
		encryptedDBName, err := encryption.Encrypt(req.DBName, req.SecretKey)
		if err != nil {
			log.Printf("Error encrypting database name: %v", err)
			return v1.JSONResponseWithError(
				c, 
				respcode.ERR_CODE_500, 
				"Failed to encrypt database name", 
				err, 
				http.StatusInternalServerError,
			)
		}
		encryptedData.DBName = encryptedDBName
	}

	// Encrypt username
	if req.Username != "" {
		encryptedUsername, err := encryption.Encrypt(req.Username, req.SecretKey)
		if err != nil {
			log.Printf("Error encrypting username: %v", err)
			return v1.JSONResponseWithError(
				c, 
				respcode.ERR_CODE_500, 
				"Failed to encrypt username", 
				err, 
				http.StatusInternalServerError,
			)
		}
		encryptedData.Username = encryptedUsername
	}

	// Encrypt password
	if req.Password != "" {
		encryptedPassword, err := encryption.Encrypt(req.Password, req.SecretKey)
		if err != nil {
			log.Printf("Error encrypting password: %v", err)
			return v1.JSONResponseWithError(
				c, 
				respcode.ERR_CODE_500, 
				"Failed to encrypt password", 
				err, 
				http.StatusInternalServerError,
			)
		}
		encryptedData.Password = encryptedPassword
	}

	// 4. Prepare response data
	responseData := map[string]interface{}{
		"encrypted": encryptedData,
	}

	// 5. Return encrypted credentials using JSONResponseWithData
	return v1.JSONResponseWithData(
		c,
		respcode.SUC_CODE_200,
		"Credentials encrypted successfully",
		responseData,
		http.StatusOK,
	)
}

// DecryptDBCredentials godoc
// @Summary Decrypt database credentials
// @Description Decrypts database connection credentials using AES encryption
// @Tags encryption
// @Accept json
// @Produce json
// @Param request body mdlEncryption.DecryptCredentialsRequest true "Credentials to decrypt"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/decrypt/db-credentials [post]
func DecryptDBCredentials(c fiber.Ctx) error {
	// 1. Parse request body
	var req struct {
		SecretKey string `json:"secretKey"`
		Host      string `json:"host,omitempty"`
		DBName    string `json:"dbName,omitempty"`
		Username  string `json:"username,omitempty"`
		Password  string `json:"password,omitempty"`
	}

	if err := c.Bind().Body(&req); err != nil {
		return v1.JSONResponseWithError(
			c, 
			respcode.ERR_CODE_400, 
			"Invalid request body", 
			err, 
			http.StatusBadRequest,
		)
	}

	// 2. Validate required fields
	if req.SecretKey == "" {
		return v1.JSONResponseWithError(
			c, 
			respcode.ERR_CODE_400, 
			"Secret key is required", 
			nil, 
			http.StatusBadRequest,
		)
	}

	if len(req.SecretKey) != 32 {
		return v1.JSONResponseWithError(
			c, 
			respcode.ERR_CODE_400, 
			"Secret key must be exactly 32 characters long", 
			nil, 
			http.StatusBadRequest,
		)
	}

	// 3. Decrypt each provided field
	decryptedData := mdlEncryption.EncryptedData{}

	if req.Host != "" {
		decryptedHost, err := encryption.Decrypt(req.Host, req.SecretKey)
		if err != nil {
			log.Printf("Error decrypting host: %v", err)
			return v1.JSONResponseWithError(
				c, 
				respcode.ERR_CODE_400, 
				"Failed to decrypt host - invalid secret key or corrupted data", 
				err, 
				http.StatusBadRequest,
			)
		}
		decryptedData.Host = decryptedHost
	}

	if req.DBName != "" {
		decryptedDBName, err := encryption.Decrypt(req.DBName, req.SecretKey)
		if err != nil {
			log.Printf("Error decrypting database name: %v", err)
			return v1.JSONResponseWithError(
				c, 
				respcode.ERR_CODE_400, 
				"Failed to decrypt database name - invalid secret key or corrupted data", 
				err, 
				http.StatusBadRequest,
			)
		}
		decryptedData.DBName = decryptedDBName
	}

	if req.Username != "" {
		decryptedUsername, err := encryption.Decrypt(req.Username, req.SecretKey)
		if err != nil {
			log.Printf("Error decrypting username: %v", err)
			return v1.JSONResponseWithError(
				c, 
				respcode.ERR_CODE_400, 
				"Failed to decrypt username - invalid secret key or corrupted data", 
				err, 
				http.StatusBadRequest,
			)
		}
		decryptedData.Username = decryptedUsername
	}

	if req.Password != "" {
		decryptedPassword, err := encryption.Decrypt(req.Password, req.SecretKey)
		if err != nil {
			log.Printf("Error decrypting password: %v", err)
			return v1.JSONResponseWithError(
				c, 
				respcode.ERR_CODE_400, 
				"Failed to decrypt password - invalid secret key or corrupted data", 
				err, 
				http.StatusBadRequest,
			)
		}
		decryptedData.Password = decryptedPassword
	}

	// 4. Prepare response data
	responseData := map[string]interface{}{
		"decrypted": decryptedData,
	}

	// 5. Return decrypted credentials using JSONResponseWithData
	return v1.JSONResponseWithData(
		c,
		respcode.SUC_CODE_200,
		"Credentials decrypted successfully",
		responseData,
		http.StatusOK,
	)
}