package mdlFeatureOne

type (
	AddExpenseRequest struct {
		Title      string  `json:"title"`
		Amount     float64 `json:"amount"`
		CategoryID *int    `json:"categoryId"`
		Date       *string `json:"date"`
		Notes      *string `json:"notes"`
		ImageURL    *string `json:"imageUrl"` 
	}
)
