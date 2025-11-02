package ctrFeatureOne

import (
	"encoding/json"
	"fmt"
	"go_template_v3/pkg/config"
	"go_template_v3/pkg/global/utils"
	mdlFeatureOne "go_template_v3/pkg/services/featureOne/model"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	v1 "github.com/FDSAP-Git-Org/hephaestus/helper/v1"
	"github.com/FDSAP-Git-Org/hephaestus/respcode"
	utils_v1 "github.com/FDSAP-Git-Org/hephaestus/utils/v1"
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

func AddExpenseV2Old(c fiber.Ctx) error {
	userId := utils.GetUserId(c)
	reqBody := mdlFeatureOne.AddExpenseRequest{}

	if err := c.Bind().Body(&reqBody); err != nil {
		return v1.JSONResponseWithError(
			c, respcode.ERR_CODE_400,
			"Invalid request body", err,
			http.StatusBadRequest,
		)
	}

	// 3. Field validation
	if strings.TrimSpace(reqBody.Title) == "" {
		return v1.JSONResponseWithError(c, respcode.ERR_CODE_400, "Title is required", nil, http.StatusBadRequest)
	}

	if reqBody.Amount <= 0 {
		return v1.JSONResponseWithError(c, respcode.ERR_CODE_400, "Amount must be greater than 0", nil, http.StatusBadRequest)
	}

	if reqBody.Date != nil && strings.TrimSpace(*reqBody.Date) != "" {
		if _, err := time.Parse("2006-01-02", *reqBody.Date); err != nil {
			return v1.JSONResponseWithError(
				c,
				respcode.ERR_CODE_400,
				"Invalid date format (expected YYYY-MM-DD)",
				err,
				http.StatusBadRequest,
			)
		}
	}

	// 4. Prepare payload for DB
	payload := map[string]interface{}{
		"userId": userId,
		"title":  reqBody.Title,
		"amount": reqBody.Amount,
	}

	if reqBody.CategoryID != nil {
		payload["categoryId"] = *reqBody.CategoryID
	}
	if reqBody.Date != nil && strings.TrimSpace(*reqBody.Date) != "" {
		payload["date"] = *reqBody.Date
	}
	if reqBody.Notes != nil && strings.TrimSpace(*reqBody.Notes) != "" {
		payload["notes"] = *reqBody.Notes
	}

	// 5. Execute the PostgreSQL function
	return utils.ExecuteDBFunction(c, "SELECT add_expense_v2($1)", payload)
}

func AddExpenseV2(c fiber.Ctx) error {
	userId := utils.GetUserId(c)

	// Parse multipart form
	form, err := c.MultipartForm()
	if err != nil {
		return v1.JSONResponseWithError(
			c, respcode.ERR_CODE_400,
			"Invalid form data", err,
			http.StatusBadRequest,
		)
	}

	// Bind handles multipart forms too
	reqBody := mdlFeatureOne.AddExpenseRequest{}
	if err := c.Bind().Body(&reqBody); err != nil {
		return v1.JSONResponseWithError(
			c, respcode.ERR_CODE_400,
			"Invalid request body", err,
			http.StatusBadRequest,
		)
	}

	// Field validation
	if strings.TrimSpace(reqBody.Title) == "" {
		return v1.JSONResponseWithError(c, respcode.ERR_CODE_400, "Title is required", nil, http.StatusBadRequest)
	}

	if reqBody.Amount <= 0 {
		return v1.JSONResponseWithError(c, respcode.ERR_CODE_400, "Amount must be greater than 0", nil, http.StatusBadRequest)
	}

	if reqBody.Date != nil && strings.TrimSpace(*reqBody.Date) != "" {
		if _, err := time.Parse("2006-01-02", *reqBody.Date); err != nil {
			return v1.JSONResponseWithError(
				c,
				respcode.ERR_CODE_400,
				"Invalid date format (expected YYYY-MM-DD)",
				err,
				http.StatusBadRequest,
			)
		}
	}

	var imageURL *string

	// Handle file upload
	if files, ok := form.File["image"]; ok && len(files) > 0 {
		fileHeader := files[0]
		config := utils.DefaultFileUploadConfig()

		uploadedPath, err := utils.UploadFile(c, fileHeader, config)
		if err != nil {
			return v1.JSONResponseWithError(
				c, respcode.ERR_CODE_400,
				"Failed to upload image", err,
				http.StatusBadRequest,
			)
		}

		// Construct full URL if needed

		fullURL := utils_v1.GetEnv("BASE_URL") + uploadedPath
		imageURL = &fullURL
	}

	// Prepare payload for DB
	payload := map[string]interface{}{
		"userId": userId,
		"title":  reqBody.Title,
		"amount": reqBody.Amount,
	}

	if reqBody.CategoryID != nil {
		payload["categoryId"] = *reqBody.CategoryID
	}
	if reqBody.Date != nil && strings.TrimSpace(*reqBody.Date) != "" {
		payload["date"] = *reqBody.Date
	}
	if reqBody.Notes != nil && strings.TrimSpace(*reqBody.Notes) != "" {
		payload["notes"] = *reqBody.Notes
	}
	if imageURL != nil {
		payload["imageUrl"] = *imageURL
	} else if reqBody.ImageURL != nil {
		// Use existing image URL if provided and no file uploaded
		payload["imageUrl"] = *reqBody.ImageURL
	}

	// Execute the PostgreSQL function
	return utils.ExecuteDBFunction(c, "SELECT add_expense_v2($1)", payload)
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

func DeleteExpenseOld(c fiber.Ctx) error {
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
	return utils.ExecuteDBFunction(c, "SELECT delete_expense($1)", payload)
}

func DeleteExpense(c fiber.Ctx) error {
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

	// 3. Execute the query and get the raw response
	result, err := utils.ExecuteDBFunctionRaw("SELECT delete_expense($1)", payload)
	if err != nil {
		return v1.JSONResponseWithError(c, respcode.ERR_CODE_500, "Database error", err, http.StatusInternalServerError)
	}

	// 4. Check if the operation was successful
	success, _ := result["success"].(bool)
	code, _ := result["code"].(float64)
	codeInt := int(code)
	codeStr := strconv.Itoa(codeInt)
	message, _ := result["message"].(string)

	if message == "" {
		message = utils.CodeMessageMap[codeStr]
	}

	if !success {
		return v1.JSONResponseWithError(c, codeStr, message, nil, codeInt)
	}

	// 5. If successful, delete the associated image file
	if data, ok := result["data"].(map[string]interface{}); ok {
		if imageUrl, exists := data["imageUrl"]; exists {
			if imageUrlStr, ok := imageUrl.(string); ok && imageUrlStr != "" {
				// Delete the image file from filesystem
				if err := utils.DeleteUploadedFile(imageUrlStr); err != nil {
					// Log the error but don't fail the request since the expense is already deleted
					fmt.Printf("Warning: Failed to delete image file %s: %v\n", imageUrlStr, err)
					// You might want to log this to a proper logging system
				}
			}
		}
	}

	// 6. Return success response
	return v1.JSONResponseWithData(c, codeStr, message, nil, codeInt)
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

func BatchUpdateExpenses(c fiber.Ctx) error {
	var successfulCount = 0
	var failedCount = 0
	// 1. Get user ID from JWT
	userId := utils.GetUserId(c)
	if userId == 0 {
		return v1.JSONResponseWithError(c, respcode.ERR_CODE_401, "User ID not found in token", nil, http.StatusUnauthorized)
	}

	// 2. Parse request body
	var req []map[string]any
	if err := c.Bind().Body(&req); err != nil {
		return v1.JSONResponseWithError(c, respcode.ERR_CODE_400, "Invalid body", err, http.StatusBadRequest)
	}

	// 3. Validate the batch update request
	if len(req) == 0 {
		return v1.JSONResponseWithError(c, respcode.ERR_CODE_400, "No updates provided", nil, http.StatusBadRequest)
	}

	if len(req) > 100 { // Limit batch size to prevent abuse
		return v1.JSONResponseWithError(c, respcode.ERR_CODE_400, "Batch size too large. Maximum 100 updates allowed", nil, http.StatusBadRequest)
	}

	// 4. Process batch updates
	results := make([]map[string]interface{}, 0, len(req))
	hasErrors := false

	for i, update := range req {
		// Create a clean payload for individual expense update
		expensePayload := make(map[string]interface{})

		// Copy only the relevant fields for individual expense update
		if expenseId, exists := update["expenseId"]; exists {
			expensePayload["expenseId"] = expenseId
		} else {
			results = append(results, map[string]interface{}{
				"index":     i,
				"message":   "Expense ID is required",
				"expenseId": expensePayload["expenseId"],
			})

			hasErrors = true
			continue
		}

		// Add userId
		expensePayload["userId"] = userId

		// Copy updateable fields
		if title, exists := update["title"]; exists {
			expensePayload["title"] = title
		}
		if amount, exists := update["amount"]; exists {
			expensePayload["amount"] = amount
		}
		if categoryId, exists := update["categoryId"]; exists {
			expensePayload["categoryId"] = categoryId
		}
		if date, exists := update["date"]; exists {
			expensePayload["date"] = date
		}
		if notes, exists := update["notes"]; exists {
			expensePayload["notes"] = notes
		}

		// Execute the update for this expense using the individual update function
		result, err := utils.ExecuteDBFunctionRaw("SELECT update_expense_v3($1)", expensePayload)
		if err != nil {
			log.Println(err.Error())
		}

		if result["success"] == true {
			successfulCount++
		} else {
			failedCount++
			results = append(results, map[string]interface{}{
				"index":     i,
				"expenseId": expensePayload["expenseId"],
				"message":   result["message"],
			})
		}

		if success, ok := result["success"].(bool); !ok || !success {
			hasErrors = true
		}
	}

	// 5. Return batch response
	responseData := map[string]interface{}{
		"results":    results,
		"total":      len(req),
		"successful": successfulCount,
		"failed":     failedCount,
	}

	if hasErrors {
		return v1.JSONResponseWithData(c, respcode.SUC_CODE_200,
			"Batch update completed with some errors", responseData, http.StatusMultiStatus)
	}

	return v1.JSONResponseWithData(c, respcode.SUC_CODE_200,
		"All expenses updated successfully", responseData, http.StatusOK)
}

func BatchUpdateExpensesAsync(c fiber.Ctx) error {
	// 1. Get user ID from JWT
	userId := utils.GetUserId(c)
	if userId == 0 {
		return v1.JSONResponseWithError(c, respcode.ERR_CODE_401, "User ID not found in token", nil, http.StatusUnauthorized)
	}

	// 2. Parse request body
	var req []map[string]interface{}
	if err := c.Bind().Body(&req); err != nil {
		return v1.JSONResponseWithError(c, respcode.ERR_CODE_400, "Invalid body", err, http.StatusBadRequest)
	}

	// 3. Validate the batch update request
	if len(req) == 0 {
		return v1.JSONResponseWithError(c, respcode.ERR_CODE_400, "No updates provided", nil, http.StatusBadRequest)
	}

	if len(req) > 100 {
		return v1.JSONResponseWithError(c, respcode.ERR_CODE_400, "Batch size too large. Maximum 100 updates allowed", nil, http.StatusBadRequest)
	}

	// 4. Create batch job record
	var jobId int
	err := config.DBConnList[0].Raw(
		"SELECT create_batch_job($1, $2, $3)",
		userId,
		"expense_batch_update",
		len(req),
	).Scan(&jobId).Error

	if err != nil {
		log.Printf("Error creating batch job: %v", err)
		return v1.JSONResponseWithError(c, respcode.ERR_CODE_500, "Failed to create batch job", err, http.StatusInternalServerError)
	}

	// 5. Process updates in background
	go processBatchUpdatesAsync(jobId, userId, req)

	// 6. Return immediately with job ID
	return v1.JSONResponseWithData(c, respcode.SUC_CODE_200,
		"Batch update job created successfully",
		map[string]interface{}{
			"jobId":      jobId,
			"totalItems": len(req),
			"status":     "pending",
		},
		http.StatusAccepted)
}

// GetBatchJobStatus - Endpoint to check job status
func GetBatchJobStatus(c fiber.Ctx) error {
	// Get user ID from JWT
	userId := utils.GetUserId(c)
	if userId == 0 {
		return v1.JSONResponseWithError(c, respcode.ERR_CODE_401, "User ID not found in token", nil, http.StatusUnauthorized)
	}

	// Get job ID from params
	jobId := c.Params("jobId")
	if jobId == "" {
		return v1.JSONResponseWithError(c, respcode.ERR_CODE_400, "Job ID is required", nil, http.StatusBadRequest)
	}

	// Query job status
	var resultStr string
	err := config.DBConnList[0].Raw(`
		SELECT jsonb_build_object(
			'jobId', id,
			'userId', user_id,
			'jobType', job_type,
			'status', status,
			'totalItems', total_items,
			'processedItems', processed_items,
			'successfulItems', successful_items,
			'failedItems', failed_items,
			'results', results,
			'createdAt', created_at,
			'updatedAt', updated_at,
			'completedAt', completed_at
		)
		FROM batch_jobs
		WHERE id = $1 AND user_id = $2
	`, jobId, userId).Scan(&resultStr).Error

	if err != nil {
		log.Printf("Error fetching job status: %v", err)
		return v1.JSONResponseWithError(c, respcode.ERR_CODE_500, "Failed to fetch job status", err, http.StatusInternalServerError)
	}

	if resultStr == "" {
		return v1.JSONResponseWithError(c, respcode.ERR_CODE_404, "Job not found", nil, http.StatusNotFound)
	}

	var jobData map[string]interface{}
	if err := json.Unmarshal([]byte(resultStr), &jobData); err != nil {
		return v1.JSONResponseWithError(c, respcode.ERR_CODE_500, "Failed to parse job data", err, http.StatusInternalServerError)
	}

	return v1.JSONResponseWithData(c, respcode.SUC_CODE_200,
		"Job status retrieved successfully",
		jobData,
		http.StatusOK)
}

// processBatchUpdatesAsync - Background worker to process batch updates
func processBatchUpdatesAsync(jobId int, userId int, updates []map[string]interface{}) {
	var successfulCount = 0
	var failedCount = 0
	results := make([]map[string]interface{}, 0)

	// Update status to processing
	updateJobStatus(jobId, "processing", 0, 0, 0, nil)

	// Process each update
	for i, update := range updates {
		// Create payload for individual expense update
		expensePayload := make(map[string]interface{})

		// Validate expense ID
		if expenseId, exists := update["expenseId"]; exists {
			expensePayload["expenseId"] = expenseId
		} else {
			failedCount++
			results = append(results, map[string]interface{}{
				"index":     i,
				"expenseId": nil,
				"message":   "Expense ID is required",
			})
			updateJobProgress(jobId, i+1, successfulCount, failedCount, results)
			continue
		}

		// Add userId
		expensePayload["userId"] = userId

		// Copy updateable fields
		if title, exists := update["title"]; exists {
			expensePayload["title"] = title
		}
		if amount, exists := update["amount"]; exists {
			expensePayload["amount"] = amount
		}
		if categoryId, exists := update["categoryId"]; exists {
			expensePayload["categoryId"] = categoryId
		}
		if date, exists := update["date"]; exists {
			expensePayload["date"] = date
		}
		if notes, exists := update["notes"]; exists {
			expensePayload["notes"] = notes
		}

		// Execute the update
		result, err := utils.ExecuteDBFunctionRaw("SELECT update_expense_v3($1)", expensePayload)
		if err != nil {
			log.Printf("Error executing update for expense %v: %v", expensePayload["expenseId"], err)
			failedCount++
			results = append(results, map[string]interface{}{
				"index":     i,
				"expenseId": expensePayload["expenseId"],
				"message":   err.Error(),
			})
		} else if result["success"] == true {
			successfulCount++
			// For successful updates, you might want to store minimal info or skip
			results = append(results, map[string]interface{}{
				"index":     i,
				"expenseId": expensePayload["expenseId"],
				"message":   "Successfully updated",
			})
		} else {
			failedCount++
			// Extract only the message from the result
			message := "Update failed"
			if msg, exists := result["message"]; exists {
				message = msg.(string)
			} else if errMsg, exists := result["error"]; exists {
				message = errMsg.(string)
			}

			results = append(results, map[string]interface{}{
				"index":     i,
				"expenseId": expensePayload["expenseId"],
				"message":   message,
			})
		}

		// Update progress after each item
		updateJobProgress(jobId, i+1, successfulCount, failedCount, results)

		// Add small delay to prevent overwhelming the database
		time.Sleep(1 * time.Minute)
	}

	// Mark job as completed
	finalStatus := "completed"
	if failedCount == len(updates) {
		finalStatus = "failed"
	}

	updateJobStatus(jobId, finalStatus, len(updates), successfulCount, failedCount, results)
}

// updateJobStatus - Helper to update job status
func updateJobStatus(jobId int, status string, processed, successful, failed int, results interface{}) {
	resultsJSON, _ := json.Marshal(results)
	err := config.DBConnList[0].Exec(
		"SELECT update_batch_job_progress($1, $2, $3, $4, $5, $6)",
		jobId,
		status,
		processed,
		successful,
		failed,
		string(resultsJSON),
	).Error

	if err != nil {
		log.Printf("Error updating job status for job %d: %v", jobId, err)
	}
}

// updateJobProgress - Helper to update job progress (only counts)
func updateJobProgress(jobId int, processed, successful, failed int, results interface{}) {
	resultsJSON, _ := json.Marshal(results)

	err := config.DBConnList[0].Exec(
		"SELECT update_batch_job_progress($1, NULL, $2, $3, $4, $5)",
		jobId,
		processed,
		successful,
		failed,
		string(resultsJSON),
	).Error

	if err != nil {
		log.Printf("Error updating job progress for job %d: %v", jobId, err)
	}
}

func TestInternalSendRequest(c fiber.Ctx) error {
	// Hardcoded payload
	payload := map[string]interface{}{
		"title":       "New Produkto",
		"price":       10,
		"description": "A description",
		"categoryId":  1,
		"images":      []string{"https://placehold.co/600x400"},
	}

	body, _ := json.Marshal(payload)

	headers := map[string]string{
		"Content-Type": "application/json",
	}

	// Make API call
	response, statusCode, err := utils_v1.SendRequestWithCode(
		"https://api.escuelajs.co/api/v1/products/",
		"POST",
		body,
		headers,
		30,
	)

	if err != nil {
		return v1.JSONResponseWithError(
			c,
			respcode.ERR_CODE_500,
			"Failed to create product",
			err,
			http.StatusInternalServerError,
		)
	}

	// Extract status code number
	statusCodeNum := http.StatusInternalServerError
	var statusCodeStr string
	parts := strings.Fields(*statusCode)
	if len(parts) > 0 {
		statusCodeStr = parts[0]
		if code, err := strconv.Atoi(parts[0]); err == nil {
			statusCodeNum = code
		}
	}

	// Check if status code is 2xx
	if strings.HasPrefix(*statusCode, "2") {
		return v1.JSONResponseWithData(
			c,
			statusCodeStr,
			"Product created successfully",
			response,
			statusCodeNum, // Use the extracted status code!
		)
	} else {
		// Handle non-2xx status codes
		errorMessage := "External API error"
		if response != nil {
			if respMap, ok := response.(map[string]interface{}); ok {
				if msg, exists := respMap["message"]; exists {
					errorMessage = fmt.Sprintf("%v", msg)
				}
			}
		}

		return v1.JSONResponseWithError(
			c,
			statusCodeStr,
			errorMessage,
			nil,
			statusCodeNum,
		)
	}
}
