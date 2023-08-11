package models

type Post struct {
	BaseModel
	Title       string `gorm:"size:15;type:string;not null,unique;"`
	Image       File   `gorm:"foreignKey:ImageId;constraint:OnUpdate:NO ACTION;OnDelete:NO ACTION"`
	ImageId     int
	Description string `gorm:"size:15;type:string;not null,unique;"`
}
