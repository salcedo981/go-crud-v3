package mdlEncryption

type EncryptCredentialsRequest struct {
	SecretKey string `json:"secretKey" validate:"required,min=32"`
	Host      string `json:"host" validate:"required"`
	DBName    string `json:"dbName" validate:"required"`
	Username  string `json:"username" validate:"required"`
	Password  string `json:"password" validate:"required"`
}

type EncryptedData struct {
	Host     string `json:"host,omitempty"`
	DBName   string `json:"dbName,omitempty"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}