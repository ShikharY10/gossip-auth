package config

import (
	"os"

	"github.com/joho/godotenv"
)

type ENV struct {
	MONGODB_CONNECTION_METHOD     string
	MONGODB_PORT                  string
	MONGODB_HOST                  string
	MONGODB_USERNAME              string
	MONGODB_PASSWORD              string
	MONGODB_CONNECTION_STRING     string
	REDIS_HOST                    string
	REDIS_PORT                    string
	AUTH_WEBSERVER_PORT           string
	JWT_ACCESS_TOKEN_SECRET_KEY   string
	JWT_REFRESH_TOKEN_SECRET_KEY  string
	GIN_MODE                      string
	CLOUDINARY_API_KEY            string
	CLOUDINARY_API_SECRET         string
	CLOUDINARY_CLOUD_NAME         string
	CLOUDINARY_AVATAR_FOLDER_NAME string
	LogServerHost                 string // 127.0.0.1
	LogServerPort                 string // 6002
	SERVICE_NAME                  string
}

func LoadENV() *ENV {
	godotenv.Load()
	var env ENV
	var value string
	var found bool

	value, found = os.LookupEnv("MONGODB_CONNECTION_METHOD")
	if found {
		env.MONGODB_CONNECTION_METHOD = value
	} else {
		env.MONGODB_CONNECTION_METHOD = "manual"
	}

	value, found = os.LookupEnv("MONGODB_PORT")
	if found {
		env.MONGODB_PORT = value
	} else {
		env.MONGODB_PORT = "27017"
	}

	value, found = os.LookupEnv("MONGODB_HOST")
	if found {
		env.MONGODB_HOST = value
	} else {
		env.MONGODB_HOST = "127.0.0.1"
	}

	value, found = os.LookupEnv("MONGODB_USERNAME")
	if found {
		env.MONGODB_USERNAME = value
	} else {
		env.MONGODB_USERNAME = "rootuser"
	}

	value, found = os.LookupEnv("MONGODB_PASSWORD")
	if found {
		env.MONGODB_PASSWORD = value
	} else {
		env.MONGODB_PASSWORD = "rootpass"
	}

	value, found = os.LookupEnv("MONGODB_CONNECTION_STRING")
	if found {
		env.MONGODB_CONNECTION_STRING = value
	} else {
		env.MONGODB_CONNECTION_STRING = ""
	}

	value, found = os.LookupEnv("REDIS_HOST")
	if found {
		env.REDIS_HOST = value
	} else {
		env.REDIS_HOST = "127.0.0.1"
	}

	value, found = os.LookupEnv("REDIS_PORT")
	if found {
		env.REDIS_PORT = value
	} else {
		env.REDIS_PORT = "6379"
	}

	value, found = os.LookupEnv("JWT_ACCESS_TOKEN_SECRET_KEY")
	if found {
		env.JWT_ACCESS_TOKEN_SECRET_KEY = value
	} else {
		env.JWT_ACCESS_TOKEN_SECRET_KEY = "982u3923jhdwhe3fjdw30fj02j3ijwef023jfijwjf802j300"
	}

	value, found = os.LookupEnv("JWT_REFRESH_TOKEN_SECRET_KEY")
	if found {
		env.JWT_REFRESH_TOKEN_SECRET_KEY = value
	} else {
		env.JWT_REFRESH_TOKEN_SECRET_KEY = "wkjjecu9qenc203foqwneocn2qifno2qenfidh2qp3hfpd2j"
	}

	value, found = os.LookupEnv("API_WEBSERVER_PORT")
	if found {
		env.AUTH_WEBSERVER_PORT = value
	} else {
		env.AUTH_WEBSERVER_PORT = "10220"
	}

	value, found = os.LookupEnv("GIN_MODE")
	if found {
		env.GIN_MODE = value
	} else {
		env.GIN_MODE = "debug"
	}

	value, found = os.LookupEnv("LOG_SERVER_HOST")
	if found {
		env.LogServerHost = value
	} else {
		env.LogServerHost = "127.0.0.1"
	}

	value, found = os.LookupEnv("LOG_SERVER_PORT")
	if found {
		env.LogServerPort = value
	} else {
		env.LogServerPort = "10223"
	}

	return &env
}
