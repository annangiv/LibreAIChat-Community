package handlers

import (
	"LibreAI/internal/authentication"
	"LibreAI/internal/database"
	"LibreAI/models"
	"LibreAI/utils"
	"encoding/json"
	"net/url"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/github"
	"github.com/markbates/goth/providers/google"
)

var store = session.New()

func InitOAuth() {
	goth.UseProviders(
		github.New(
			os.Getenv("GITHUB_CLIENT_ID"),
			os.Getenv("GITHUB_CLIENT_SECRET"),
			os.Getenv("BASE_URL")+"/auth/github/callback",
			"user:email",
		),
		google.New(
			os.Getenv("GOOGLE_CLIENT_ID"),
			os.Getenv("GOOGLE_CLIENT_SECRET"),
			os.Getenv("BASE_URL")+"/auth/google/callback",
			"email", "profile",
		),
	)
}

func AuthPageHandler(c *fiber.Ctx) error {
	return c.Render("auth", fiber.Map{
		"Title": "Sign In - LibreAI",
	}, "layouts/main")
}

// GitHubLogin starts the GitHub OAuth flow
func GitHubLogin(c *fiber.Ctx) error {
	provider, err := goth.GetProvider("github")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Invalid provider")
	}

	sess, err := store.Get(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Session error")
	}

	authURL, err := provider.BeginAuth(sess.ID())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("OAuth error: " + err.Error())
	}

	url, _ := authURL.GetAuthURL()
	sess.Set("github", authURL.Marshal())
	sess.Save()
	return c.Redirect(url)
}

// GoogleLogin starts the Google OAuth flow
func GoogleLogin(c *fiber.Ctx) error {
	provider, err := goth.GetProvider("google")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Invalid provider")
	}

	sess, err := store.Get(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Session error")
	}

	authURL, err := provider.BeginAuth(sess.ID())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("OAuth error: " + err.Error())
	}

	url, _ := authURL.GetAuthURL()
	sess.Set("google", authURL.Marshal())
	sess.Save()
	return c.Redirect(url)
}

// GitHubCallback handles the GitHub OAuth callback
func GitHubCallback(c *fiber.Ctx) error {
	return handleOAuthCallback(c, "github")
}

// GoogleCallback handles the Google OAuth callback
func GoogleCallback(c *fiber.Ctx) error {
	return handleOAuthCallback(c, "google")
}

// handleOAuthCallback handles the OAuth callback for both providers
func handleOAuthCallback(c *fiber.Ctx, providerName string) error {
	sess, err := store.Get(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Session error")
	}

	provider, err := goth.GetProvider(providerName)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Provider error")
	}

	raw := sess.Get(providerName)
	if raw == nil {
		return c.Status(fiber.StatusBadRequest).SendString("Missing session")
	}

	authSess, err := provider.UnmarshalSession(raw.(string))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Unmarshal error: " + err.Error())
	}

	params := url.Values{}
	c.Request().URI().QueryArgs().VisitAll(func(key, value []byte) {
		params.Add(string(key), string(value))
	})

	_, err = authSess.Authorize(provider, params)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).SendString("Authorization failed: " + err.Error())
	}

	oauthUser, err := provider.FetchUser(authSess)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("User fetch error: " + err.Error())
	}

	// ✅ Check if user already exists
	db := database.Get()
	var existingUser models.User
	if err := db.Where("provider = ? AND provider_id = ?", providerName, oauthUser.UserID).First(&existingUser).Error; err == nil {
		// Use authentication helper to set session/cookie
		authentication.AuthStore(c, existingUser.ID)
		return c.Redirect("/ask")
	}

	// Not found → proceed to consent

	// Serialize user info
	userJSON, err := json.Marshal(oauthUser)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to serialize user data")
	}

	sess.Set("oauth_user_data", string(userJSON))
	sess.Set("oauth_provider", providerName)
	sess.Save()

	// Render consent page
	return c.Render("consent", fiber.Map{
		"UserInfo": fiber.Map{
			"Provider":  providerName,
			"Username":  oauthUser.NickName,
			"Name":      oauthUser.Name,
			"Email":     utils.MaskEmail(oauthUser.Email),
			"AvatarURL": oauthUser.AvatarURL,
		},
	}, "layouts/main")
}

func ProcessConsent(c *fiber.Ctx) error {
	sess, err := store.Get(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Session error")
	}

	// Check if consent is given
	consent := c.FormValue("consent")
	if consent == "" {
		// No consent given, redirect back to auth
		return c.Redirect("/auth?error=consent_required")
	}

	// Retrieve stored OAuth user data
	rawUserData := sess.Get("oauth_user_data")
	providerName := sess.Get("oauth_provider")

	if rawUserData == nil || providerName == nil {
		return c.Status(fiber.StatusBadRequest).SendString("Missing OAuth session data")
	}

	// Unmarshal user data
	var user goth.User
	if err := json.Unmarshal([]byte(rawUserData.(string)), &user); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to parse user data")
	}

	db := database.Get()
	var existing models.User

	// Check if user already exists
	if err := db.Where("provider_id = ? AND provider = ?", user.UserID, providerName).First(&existing).Error; err != nil {
		// Create new user
		newUser := models.User{
			Name:       user.Name,
			Email:      user.Email,
			Username:   user.NickName,
			Image:      user.AvatarURL,
			AvatarURL:  user.AvatarURL,
			ProviderID: user.UserID,
			Provider:   providerName.(string),
			CreatedAt:  time.Now(),
		}

		// Fallback username if none provided
		if newUser.Username == "" {
			newUser.Username = "user" + user.UserID
		}

		// Create the user
		if err := db.Create(&newUser).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Failed to create user")
		}

		// Authenticate the new user
		authentication.AuthStore(c, newUser.ID)
	} else {
		// User exists, authenticate
		authentication.AuthStore(c, existing.ID)
	}

	// Clear temporary OAuth session data
	sess.Delete("oauth_user_data")
	sess.Delete("oauth_provider")
	sess.Save()

	// Redirect to ask page
	return c.Redirect("/ask")
}

func CancelOAuth(c *fiber.Ctx) error {
	sess, err := store.Get(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Session error")
	}

	// Clear all OAuth-related session data
	sess.Delete("oauth_user_data")
	sess.Delete("oauth_provider")
	sess.Save()

	// Redirect back to auth page
	return c.Redirect("/auth")
}

func Logout(c *fiber.Ctx) error {
	sess, err := store.Get(c)
	if err == nil {
		sess.Destroy()
	}
	authentication.AuthClear(c)
	return c.Redirect("/")
}
