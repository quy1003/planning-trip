package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Env struct {
	DatabaseURL          string
	Port                 string
	AuthSecret           string
	CloudinaryCloudName  string
	CloudinaryAPIKey     string
	CloudinaryAPISecret  string
	CloudinaryUploadPath string
}

const defaultPort = "8080"

func Load() (Env, error) {
	_ = godotenv.Overload()

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		return Env{}, fmt.Errorf("missing required environment variable: DATABASE_URL")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	return Env{
		DatabaseURL:          databaseURL,
		Port:                 port,
		AuthSecret:           getOrDefault("AUTH_SECRET", "planning-trip-dev-secret"),
		CloudinaryCloudName:  getOrDefault("CLOUDINARY_CLOUD_NAME", ""),
		CloudinaryAPIKey:     getOrDefault("CLOUDINARY_API_KEY", ""),
		CloudinaryAPISecret:  getOrDefault("CLOUDINARY_API_SECRET", ""),
		CloudinaryUploadPath: getOrDefault("CLOUDINARY_UPLOAD_PATH", "planning-trip"),
	}, nil
}

func getOrDefault(value string, fallback string) string {
	res := os.Getenv(value)
	if res == "" {
		return fallback
	}
	return res
}
