package mdlFeatureOne

type (
	CategoryRequest struct {
		Name        *string `json:"name"`
		Description *string `json:"description"`
	}

	CategoryResponse struct {
		ID          *int    `json:"id"`
		Name        *string `json:"name"`
		Description *string `json:"description"`
		CreatedAt   *string `json:"createdAt"`
		UpdatedAt   *string `json:"updatedAt"`
	}

	GetExpenseCategoriesResponse struct {
		Categories *[]CategoryResponse `json:"categories"`
		Total      *int                `json:"total"`
		Limit      *int                `json:"limit"`
		Offset     *int                `json:"offset"`
	}
)
