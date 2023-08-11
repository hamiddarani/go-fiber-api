package models

type User struct {
	BaseModel
	MobileNumber string `gorm:"type:string;size:11;null;unique;default:null"`
	Enabled      bool   `gorm:"default:true"`
	UserRoles    *[]UserRole
}

type Role struct {
	BaseModel
	Name      string `gorm:"type:string;size:10;not null,unique"`
	UserRoles *[]UserRole
}

type UserRole struct {
	BaseModel
	User   User `gorm:"foreignKey:UserId;constraint:OnUpdate:NO ACTION;OnDelete:NO ACTION"`
	Role   Role `gorm:"foreignKey:RoleId;constraint:OnUpdate:NO ACTION;OnDelete:NO ACTION"`
	UserId int
	RoleId int
}
