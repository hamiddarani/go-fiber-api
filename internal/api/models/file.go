package models

type File struct {
	BaseModel
	Name        string `gorm:"size:100;type:string;not null"`
	Directory   string `gorm:"size:100;type:string;not null"`
	MimeType    string `gorm:"size:20;type:string;not null"`
}
