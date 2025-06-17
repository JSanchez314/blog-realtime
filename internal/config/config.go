package config

import (
	"log"
	"os"
)

type Config struct {
	Port        string
	PgURL       string
	MongoURI    string
	MongoCAPath string
	JwtSecret   string
}

func Load() Config {
	c := Config{
		Port:        os.Getenv("PORT"),
		PgURL:       os.Getenv("PG_URL"),
		MongoURI:    os.Getenv("MONGO_URI"),
		MongoCAPath: os.Getenv("MONGO_CA_PATH"),
		JwtSecret:   os.Getenv("JWT_SECRET"),
	}

	if c.Port == "" || c.PgURL == "" || c.MongoURI == "" || c.JwtSecret == "" {
		log.Fatal("Missing environment variables")
	}

	return c
}
