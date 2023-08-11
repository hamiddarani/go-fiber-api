package dto

import "mime/multipart"

type FileFormRequest struct {
	File *multipart.FileHeader `json:"file" form:"file" validate:"required" swaggerignoer:"true"`
}

type UploadFileRequest struct {
	FileFormRequest
}

type CreateFileRequest struct {
	Name      string `json:"name"`
	Directory string `json:"directory"`
	MimeType  string `json:"mimeType"`
}

type UpdateFileRequest struct {
	Description string `json:"description"`
}

type FileResponse struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Directory   string `json:"directory"`
	MimeType    string `json:"mimeType"`
}

type FileFilter struct{}
