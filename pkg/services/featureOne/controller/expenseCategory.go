package ctrFeatureOne

import (
	"encoding/json"
	"go_template_v3/pkg/config"
	mdlFeatureOne "go_template_v3/pkg/services/featureOne/model"
	scpFeatureOne "go_template_v3/pkg/services/featureOne/script"
	"log"
	"net/http"

	v1 "github.com/FDSAP-Git-Org/hephaestus/helper/v1"
	"github.com/FDSAP-Git-Org/hephaestus/respcode"
	"github.com/gofiber/fiber/v3"
)

func AddExpenseCategory(c fiber.Ctx) error {
	var req mdlFeatureOne.CategoryRequest
	if err := c.Bind().Body(&req); err != nil {
		return v1.JSONResponseWithError(c, respcode.ERR_CODE_400, "Invalid body", err, http.StatusBadRequest)
	}

	// Validate required fields
	if req.Name == nil || *req.Name == "" {
		return v1.JSONResponseWithError(c, respcode.ERR_CODE_400, "Category name is required", nil, http.StatusBadRequest)
	}

	query := scpFeatureOne.AddExpenseCategory
	var resultStr string
	err := config.DBConnList[0].Raw(
		query,
		req.Name,
		req.Description,
	).Scan(&resultStr).Error // Scan into string
	if err != nil {
		return v1.JSONResponseWithError(c, respcode.ERR_CODE_500, "DB error", err, http.StatusInternalServerError)
	}

	// Convert string to []byte for JSON unmarshalling
	resultJSON := []byte(resultStr)
	var res mdlFeatureOne.CategoryResponse
	if err := json.Unmarshal(resultJSON, &res); err != nil {
		return v1.JSONResponseWithError(c, respcode.ERR_CODE_500, "Parse error", err, http.StatusInternalServerError)
	}

	return v1.JSONResponseWithData(c, respcode.SUC_CODE_200, "Category created", res, http.StatusOK)
}

func GetExpenseCategories(c fiber.Ctx) error {
	// Parse query parameters
	limit := fiber.Query[int](c, "limit")
	offset := fiber.Query[int](c, "offset")

    if limit <= 0 {limit = 50}
    if offset < 0 {offset = 0}

	var resultStr string
	err := config.DBConnList[0].Debug().Raw(
		scpFeatureOne.GetExpenseCategories,
		limit,
		offset,
	).Scan(&resultStr).Error

	if err != nil {
		log.Printf("Database error: %v", err)
		return v1.JSONResponseWithError(c, respcode.ERR_CODE_500, "DB error", err, http.StatusInternalServerError)
	}

	// Convert string to []byte for JSON unmarshalling
	resultJSON := []byte(resultStr)
	var res mdlFeatureOne.GetExpenseCategoriesResponse
	if err := json.Unmarshal(resultJSON, &res); err != nil {
		log.Printf("JSON unmarshal error: %v, raw: %s", err, resultStr)
		return v1.JSONResponseWithError(c, respcode.ERR_CODE_500, "Parse error", err, http.StatusInternalServerError)
	}

	return v1.JSONResponseWithData(c, respcode.SUC_CODE_200, "Categories retrieved", res, http.StatusOK)
}
