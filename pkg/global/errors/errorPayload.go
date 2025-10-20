package errors

type (
	ErrorModel struct {
		Message string      `json:"message"`
		Error   error       `json:"error"`
		Data    interface{} `json:"data,omitempt"`
	}
)
