package web

const ErrUnknown string = "Unknown error occurred at the backend!"
const ErrInvalidPass string = "Invalid password!"
const ErrUserNotFound string = "User not found!"

const SqlSelectUsernamePass string = "SELECT encrypted_password, user_name FROM users WHERE user_name='%s'"
const SqlUpdateUser string = "UPDATE users SET version = version + 1, last_login='%v', access_token='%v', " +
	"access_token_expires_at='%v', refresh_token='%s', refresh_token_expires_at='%v' WHERE user_name='%s'"
const SqlSelectUserAll string = "SELECT uuid, id, encrypted_password, created_at, updated_at, version, user_name, " +
	"email, last_login, enabled, email_verified, access_token, access_token_expires_at, refresh_token, " +
	"refresh_token_expires_at FROM users WHERE user_name='%s'"
