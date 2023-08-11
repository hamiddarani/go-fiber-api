package dto

type CreateUpdatePostRequest struct {
	Title       string `json:"title" validate:"required,alpha,min=3,max=20"`
	Description string `json:"description" validate:"required,alpha,max=200"`
	ImageId     int    `json:"imageId" validate:"required"`
}

type PostResponse struct {
	Id          int          `json:"id"`
	Title       string       `json:"title"`
	Description string       `json:"description"`
	Image       FileResponse `json:"image"`
}

type PostFilter struct{}
