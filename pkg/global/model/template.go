package model

import "time"

type (
	TemplateDetails struct {
		Id             int       `json:"id"`
		UploadType     string    `json:"uploadType"`
		TemplateSerial string    `json:"templateSerial"`
		FileData       []byte    `json:"fileData"`
		FileExt        string    `json:"fileExt"`
		EncodedBy      string    `json:"encodedBy"`
		CreatedAt      time.Time `json:"createdAt"`
	}

	ViewTemplateDetails struct {
		UploadType     string    `json:"uploadType"`
		TemplateSerial string    `json:"templateSerial"`
		FileExt        string    `json:"fileExt"`
		EncodedBy      string    `json:"encodedBy"`
		CreatedAt      time.Time `json:"createdAt"`
	}
)
