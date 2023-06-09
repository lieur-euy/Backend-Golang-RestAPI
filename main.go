package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	googleOauthConfig *oauth2.Config
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	googleOauthConfig = &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("GOOGLE_REDIRECT_URI"),
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
		},
		Endpoint: google.Endpoint,
	}
}

func main() {
	app := fiber.New()

	app.Get("/auth/google", handleGoogleLogin)
	app.Get("/auth/google/callback", handleGoogleCallback)

	log.Fatal(app.Listen(":3000"))
}

func handleGoogleLogin(c *fiber.Ctx) error {
	url := googleOauthConfig.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	return c.Redirect(url)
}

func handleGoogleCallback(c *fiber.Ctx) error {
	code := c.Query("code")

	token, err := googleOauthConfig.Exchange(oauth2.NoContext, code)
	if err != nil {
		return err
	}

	client := googleOauthConfig.Client(oauth2.NoContext, token)
	email, err := getUserEmail(client)
	if err != nil {
		return err
	}

	// Lakukan proses autentikasi di sini sesuai kebutuhan Anda
	// Misalnya, periksa apakah email pengguna ada dalam database

	return c.SendString(fmt.Sprintf("Logged in with email: %s", email))
}

func getUserEmail(client *http.Client) (string, error) {
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var data struct {
		Email string `json:"email"`
	}

	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return "", err
	}

	return data.Email, nil
}
