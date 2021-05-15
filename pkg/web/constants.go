package web

const errUnknown string = "Unknown error occurred at the backend!"
const errInvalidPass string = "Invalid password!"
const errUserNotFound string = "User not found!"

const sqlSelectUsernamePass string = "SELECT encrypted_password, user_name FROM users WHERE user_name='%s'"
const sqlUpdateUser string = "UPDATE users SET version = version + 1, last_login='%v', access_token='%v', " +
	"access_token_expires_at='%v', refresh_token='%s', refresh_token_expires_at='%v' WHERE user_name='%s'"
const sqlSelectUserAll string = "SELECT uuid, id, encrypted_password, created_at, updated_at, version, user_name, " +
	"email, last_login, enabled, email_verified, access_token, access_token_expires_at, refresh_token, " +
	"refresh_token_expires_at FROM users WHERE user_name='%s'"
