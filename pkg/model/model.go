package model

type User struct {
	Id                     uint `gorm:"primary_key,AUTO_INCREMENT"`
	Uuid                   string
	UserName               string
	EncryptedPassword      string
	Email                  string
	VerificationCode       uint
	AccessToken            string
	AccessTokenExpiresAt   string
	RefreshToken           string
	RefreshTokenExpiresAt  string
	Enabled                bool
	EmailVerified          bool
	VerificationCodeUsable bool
	// check https://gorm.io/docs/models.html#Field-Level-Permission for permissions
	VerificationCodeCreatedAt  string `gorm:"->"`
	VerificationCodeVerifiedAt string `gorm:"->"`
	FailedLoginAttempts        uint
	LastLogin                  string
	Version                    uint
	CreatedAt                  string
	UpdatedAt                  string
	Roles                      []*Role `gorm:"many2many:users_roles;->"`
}

type Role struct {
	Id        uint    `gorm:"primary_key,AUTO_INCREMENT" json:"id"`
	Name      string  `json:"name"`
	Version   uint    `json:"version"`
	CreatedAt string  `json:"createdAt"`
	UpdatedAt string  `json:"updatedAt"`
	Users     []*User `gorm:"many2many:users_roles" json:"users"`
}
