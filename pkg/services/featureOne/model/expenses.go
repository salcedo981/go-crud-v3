package mdlFeatureOne

type (
	// Request and response must be a pointer to make sure omitted values become null in db and in rest response
	ExpenseRequest struct {
		Title      *string  `json:"title"`
		Amount     *float64 `json:"amount"`
		CategoryId *int     `json:"categoryId"`
		Date       *string  `json:"date"`
		Notes      *string  `json:"notes"`
	}

	ExpenseResponse struct {
		Id        *int             `json:"id"`
		Title     *string          `json:"title"`
		Amount    *float64         `json:"amount"`
		Category  *ExpenseCategory `json:"category"`
		Date      *string          `json:"date"`
		Notes     *string          `json:"notes"`
		CreatedBy *int             `json:"createdBy"` 
		CreatedAt *string          `json:"createdAt"`
		UpdatedAt *string          `json:"updatedAt"`
	}

	ExpenseCategory struct {
		ID          *int    `json:"id"`
		Name        *string `json:"name"`
		Description *string `json:"description"`
		CreatedAt   *string `json:"createdAt"`
		UpdatedAt   *string `json:"updatedAt"`
	}

	GetExpensesResponse struct {
		Expenses *[]ExpenseResponse `json:"expenses"`
		Total    *int64             `json:"total"`
		Limit    *int               `json:"limit"`
		Offset   *int               `json:"offset"`
	}
)