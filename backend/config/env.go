package config

import "os"

var (
	Port             = os.Getenv("PORT")
	Token            = os.Getenv("TOKEN")
	Secret           = os.Getenv("SECRET")
	DatabasePort     = os.Getenv("DATABASE_PORT")
	DatabaseName     = os.Getenv("DATABASE_NAME")
	DatabaseUser     = os.Getenv("DATABASE_USER")
	DatabasePassword = os.Getenv("DATABASE_PASSWORD")
	MySigningKey     = "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
	SaltDB           = os.Getenv("SALTDB")
)
