package mdlFeatureOne

type (
	SampleModel struct {
		Id        int    `json:"id"`
		Code      string `json:"code"`
		Name      string `json:"name"`
		EncodedBy string `json:"encodedBy"`
		CreatedAt string `json:"createdAt"`
		UpdatedAt string `json:"updatedAt"`
	}
)
