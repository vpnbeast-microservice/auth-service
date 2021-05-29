package types

type User struct {
	Id                int
	Uuid              string
	UserName          string
	EncryptedPassword string
	Email             string
	VerificationCode  int
	AccessToken       string
	// AccessTokenExpiresAt	string
	RefreshToken string
	// RefreshTokenExpiresAt	string
	Enabled                bool
	EmailVerified          bool
	VerificationCodeUsable bool
	// VerificationCodeCreatedAt	string
	// VerificationCodeVerifiedAt	string
	FailedLoginAttempts int
	LastLogin           string
	Version             int
	CreatedAt           string
	UpdatedAt           string
}
