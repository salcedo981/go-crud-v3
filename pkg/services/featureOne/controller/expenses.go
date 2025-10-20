package ctrFeatureOne

import (
	"go_template_v3/pkg/global/utils"
	"net/http"

	v1 "github.com/FDSAP-Git-Org/hephaestus/helper/v1"
	"github.com/FDSAP-Git-Org/hephaestus/respcode"
	"github.com/gofiber/fiber/v3"
)

func AddExpense(c fiber.Ctx) error {
	// 1. Get user ID from JWT
	userId := utils.GetUserId(c)

	// 2. Parse request body
	var reqBody map[string]interface{}
	if err := c.Bind().Body(&reqBody); err != nil {
		return v1.JSONResponseWithError(c, respcode.ERR_CODE_400, "Invalid body", err, http.StatusBadRequest)
	}

	// 3. Add userId to the payload (controller logic)
	reqBody["userId"] = userId

	
	// 4. Execute the query
	return utils.ExecuteDBFunction(c, "SELECT add_expense_v3($1)", reqBody)
}

func UpdateExpense(c fiber.Ctx) error {
	// 1. Get user ID from JWT
	userId := utils.GetUserId(c)
	if userId == 0 {
		return v1.JSONResponseWithError(c, respcode.ERR_CODE_401, "User ID not found in token", nil, http.StatusUnauthorized)
	}

	// 2. Parse request body
	var reqBody map[string]interface{}
	if err := c.Bind().Body(&reqBody); err != nil {
		return v1.JSONResponseWithError(c, respcode.ERR_CODE_400, "Invalid body", err, http.StatusBadRequest)
	}

	// 3. Add userId and ensure we have the expense ID
	reqBody["userId"] = userId
	reqBody["expenseId"] = c.Params("id") // Get ID from URL params

	// 4. Execute the query
	return utils.ExecuteDBFunction(c, "SELECT update_expense_v3($1)", reqBody)
}

func GetExpenses(c fiber.Ctx) error {
	// 1. Get user ID from JWT
	userId := utils.GetUserId(c)

	// 2. Prepare payload from query parameters
	payload := map[string]interface{}{
		"userId": userId,
	}

	// Extract and add optional query parameters
	if title := fiber.Query[string](c, "title"); title != "" {
		payload["title"] = title
	}
	if amount := fiber.Query[float64](c, "amount"); amount != 0 {
		payload["amount"] = amount
	}
	if minAmount := fiber.Query[float64](c, "minAmount"); minAmount != 0 {
		payload["minAmount"] = minAmount
	}
	if maxAmount := fiber.Query[float64](c, "maxAmount"); maxAmount != 0 {
		payload["maxAmount"] = maxAmount
	}
	if categoryId := fiber.Query[int](c, "categoryId"); categoryId != 0 {
		payload["categoryId"] = categoryId
	}
	if limit := fiber.Query[int](c, "limit"); limit != 0 {
		payload["limit"] = limit
	}
	if offset := fiber.Query[int](c, "offset"); offset != 0 {
		payload["offset"] = offset
	}
	if startDate := fiber.Query[string](c, "startDate"); startDate != "" {
		payload["startDate"] = startDate
	}
	if endDate := fiber.Query[string](c, "endDate"); endDate != "" {
		payload["endDate"] = endDate
	}

	// 3. Execute the query
	return utils.ExecuteDBFunction(c, "SELECT get_expenses_v3($1)", payload)

}

func GetExpense(c fiber.Ctx) error {
	// 1. Get user ID from JWT
	userId := utils.GetUserId(c)
	if userId == 0 {
		return v1.JSONResponseWithError(c, respcode.ERR_CODE_401, "User ID not found in token", nil, http.StatusUnauthorized)
	}

	// 2. Build payload (controller logic)
	payload := map[string]interface{}{
		"userId":    userId,
		"expenseId": c.Params("id"),
	}

	// 3. Execute the query
	return utils.ExecuteDBFunction(c, "SELECT get_expense_v3($1)", payload)
}

func DeleteExpense(c fiber.Ctx) error {
	// 1. Get user ID from JWT
	userId := utils.GetUserId(c)
	if userId == 0 {
		return v1.JSONResponseWithError(c, respcode.ERR_CODE_401, "User ID not found in token", nil, http.StatusUnauthorized)
	}

	// 2. Build payload (controller logic)
	payload := map[string]interface{}{
		"userId": userId,
		"id":     c.Params("id"),
	}

	// 3. Execute the query
	return utils.ExecuteDBFunction(c, "SELECT delete_expense($1)", payload)
}

func AddCategory(c fiber.Ctx) error {
	// 1. Parse request body
	var reqBody map[string]interface{}
	if err := c.Bind().Body(&reqBody); err != nil {
		return v1.JSONResponseWithError(c, respcode.ERR_CODE_400, "Invalid body", err, http.StatusBadRequest)
	}

	// 2. No userId needed for categories (controller decision)
	// Just pass the request body as-is

	// 3. Execute the query
	return utils.ExecuteDBFunction(c, "SELECT add_category($1)", reqBody)
}

// GetExpenses godoc
// @Summary Get expenses with filtering and pagination
// @Description Get expenses with various filters and pagination
// @Tags expenses
// @Accept json
// @Produce json
// @Param title query string false "Filter by title (contains search)"
// @Param amount query number false "Filter by exact amount"
// @Param min_amount query number false "Filter by minimum amount"
// @Param max_amount query number false "Filter by maximum amount"
// @Param category_id query integer false "Filter by category ID"
// @Param limit query integer false "Limit results (default: 50)" default(50)
// @Param offset query integer false "Offset for pagination" default(0)
// @Param start_date query string false "Start date (YYYY-MM-DD)"
// @Param end_date query string false "End date (YYYY-MM-DD)"
// @Success 200 {object} mdlFeatureOne.GetExpensesResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/expenses [get]
