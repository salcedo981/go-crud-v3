package mdlFeatureOne

type (
    RegisterRequest struct {
        Email    *string `json:"email"`
        Password *string `json:"password"`
        Name     *string `json:"name"`
    }

    LoginRequest struct {
        Email    *string `json:"email"`
        Password *string `json:"password"`
    }

	 User struct {
		Id       *int    `json:"id"`
		Email    *string `json:"email"`
		Password *string `json:"password"`
		Name     *string `json:"name"`
	}
)