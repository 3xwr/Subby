package model

type UploadResponse struct {
	FileAddress string `json:"file_address"`
	Uploaded    bool   `json:"uploaded"`
}
