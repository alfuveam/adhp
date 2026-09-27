package config

import "os"

var (
	Port              = os.Getenv("PORT")
	Token             = os.Getenv("TOKEN")
	Secret            = os.Getenv("SECRET")
	DatabaseHost      = os.Getenv("DATABASE_HOST")
	DatabasePort      = os.Getenv("DATABASE_PORT")
	DatabaseName      = os.Getenv("DATABASE_NAME")
	DatabaseUser      = os.Getenv("DATABASE_USER")
	DatabasePassword  = os.Getenv("DATABASE_PASSWORD")
	MySigningKey      = os.Getenv("MY_SIGNING_KEY")
	SaltDB            = os.Getenv("SALTDB")
	CodeHandlerApiUrl = os.Getenv("CODE_HANDLER_API_URL")
	CodeHandlerKey    = os.Getenv("CODE_HANDLER_KEY")
	CorsOrigin        = os.Getenv("CORS_ORIGIN")
)