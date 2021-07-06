package constant

const (
	HOST = "localhost"
	PORT = "8080"
	URL  = "http://" + HOST + ":" + PORT

	PATH_AUTH_GOOGLE          = "/auth/google"
	PATH_AUTH_GOOGLE_CALLBACK = PATH_AUTH_GOOGLE + "/cb"

	URL_AUTH_GOOGLE          = URL + PATH_AUTH_GOOGLE
	URL_AUTH_GOOGLE_CALLBACK = URL + PATH_AUTH_GOOGLE_CALLBACK

	SESSION_FLASH                = "flash"
	SESSION_AUTH                 = "gid"
	SESSION_AUTH_CONTENT_SESS_ID = "session_id"
)

var (
	SKIP_USER_SESSION_MIDDLEWARE = []string{
		PATH_AUTH_GOOGLE,
		PATH_AUTH_GOOGLE_CALLBACK,
	}
)
