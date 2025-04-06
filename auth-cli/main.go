package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/iulian509/realtime-messaging/config"
)

func main() {
	action := flag.String("action", "", "Action to perform (create)")
	username := flag.String("username", "", "Username for the token")

	flag.Parse()

	if *action != "create" {
		flag.Usage()
		os.Exit(1)
	}

	if *username == "" {
		log.Fatal("username is required for create action")
	}

	token, err := generateJWT(*username)
	if err != nil {
		log.Fatalf("failed to generate JWT: %v", err)
	}

	fmt.Printf("generated JWT for user %s: %s\n", *username, token)
}

func generateJWT(username string) (string, error) {
	config, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load YAML configuration")
	}

	claims := jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Hour * 1).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.JWT.SecretKey))
}
