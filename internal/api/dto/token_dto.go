package dto

type TokenDto struct {
	UserId       int
	MobileNumber string
	Roles        []string
}

type GetOtpRequestDto struct {
	MobileNumber string `json:"mobileNumber" validate:"required,mobile,min=11,max=11"`
}

type TokenDetail struct {
	AccessToken            string `json:"accessToken"`
	RefreshToken           string `json:"refreshToken"`
	AccessTokenExpireTime  int64  `json:"accessTokenExpireTime"`
	RefreshTokenExpireTime int64  `json:"refreshTokenExpireTime"`
}

type RegisterLoginByMobileRequest struct {
	MobileNumber string `json:"mobileNumber" validate:"required,mobile,min=11,max=11"`
	Otp          string `json:"otp" validate:"required,min=6,max=6"`
}

type RefreshToken struct {
	RefreshToken string `json:"refreshToken" validate:"required"`
}

type AccessToken struct {
	AccessToken           string `json:"accessToken"`
	AccessTokenExpireTime int64  `json:"accessTokenExpireTime"`
}
