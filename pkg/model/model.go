package model

type User struct {
	Id                         int
	Uuid                       string
	UserName                   string
	EncryptedPassword          string
	Email                      string
	VerificationCode           int
	AccessToken                string
	AccessTokenExpiresAt       string
	RefreshToken               string
	RefreshTokenExpiresAt      string
	Enabled                    bool
	EmailVerified              bool
	VerificationCodeUsable     bool
	// check https://gorm.io/docs/models.html#Field-Level-Permission for permissions
	VerificationCodeCreatedAt  string `gorm:"->"`
	VerificationCodeVerifiedAt string `gorm:"->"`
	FailedLoginAttempts        int
	LastLogin                  string
	Version                    int
	CreatedAt                  string
	UpdatedAt                  string
}
