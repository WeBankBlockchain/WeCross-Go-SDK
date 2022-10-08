package common

const (
	APPLICATION_CONFIG_FILE = "classpath:application.toml"
	TEMPLATE_URL            = "http://127.0.0.1:8250/"
	XA_TRANSACTION_ID_KEY   = "XA_TRANSACTION_ID"
	XA_TRANSACTION_SEQ_KEY  = "XA_TRANSACTION_SEQ"

	// alphabet(upper or lower case) + digit + character("_-") ，length in (4,16)
	USERNAME_PATTERN = "^[a-zA-Z0-9_-]{4,16}$"

	//alphabet(upper or lower case) + digit + character("@+!%*#?") ，length in (4,16)
	PASSWORD_PATTERN = "^[A-Za-z0-9@+!%*#?]{1,16}$"

	SSL_OFF            = 2
	SSL_ON             = 1
	SSL_ON_CLIENT_AUTH = 0
)
