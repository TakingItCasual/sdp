package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// User is a retrieved and authentiacted user.
type User struct {
	Sub           string `json:"sub"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Profile       string `json:"profile"`
	Picture       string `json:"picture"`
	Email         string `json:"email"`
	EmailVerified string `json:"email_verified"`
	Gender        string `json:"gender"`
}

var conf *oauth2.Config

func randToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.StdEncoding.EncodeToString(b)
}

func init() {
	conf = &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRECT"),
		RedirectURL:  "http://127.0.0.1:9090/callback",
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}
}

func googleCallbackHandler(c *gin.Context) {
	// Read oauthState from Cookie
	oauthState, _ := c.Cookie("googleOauthstate")

	if c.Query("googleOauthstate") != oauthState {
		log.Println("invalid oauth google state")
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}

	data, err := getGoogleUserData(c)
	if err != nil {
		log.Println(err.Error())
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}

	// GetOrCreate User in your db.
	// Redirect or response with a token.
	// More code .....
	fmt.Fprintf(c.Writer, "UserInfo: %s\n", data)
}

func getGoogleUserData(c *gin.Context) ([]byte, error) {
	// Use code to get token and get user info from Google.
	tok, err := conf.Exchange(oauth2.NoContext, c.Query("code"))
	if err != nil {
		return nil, err
	}

	client := conf.Client(oauth2.NoContext, tok)
	email, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		return nil, err
	}

	defer email.Body.Close()
	data, _ := ioutil.ReadAll(email.Body)
	log.Println("Email body: ", string(data))
	c.Status(http.StatusOK)

	return data, nil
}

func googleLoginHandler(c *gin.Context) {
	oauthState := generateStateOauthCookie(c)
	loginURL := conf.AuthCodeURL(oauthState)
	c.Redirect(http.StatusSeeOther, loginURL)
}

func generateStateOauthCookie(c *gin.Context) string {
	b := make([]byte, 32)
	rand.Read(b)
	randData := base64.URLEncoding.EncodeToString(b)
	cookie := http.Cookie{
		Name:     "googleOauthstate",
		Value:    randData,
		Expires:  time.Now().Add(time.Hour), // Must log in in under an hour
		HttpOnly: true,                      // Cookie only accessible from backend
	}
	http.SetCookie(c.Writer, &cookie)

	return randData
}

func main() {
	router := gin.Default()

	// CORS is enabled to allow the backend and frontend to communicate
	//config := cors.DefaultConfig()
	//config.AllowOrigins = []string{"http://google.com"}
	//router.Use(cors.New(config))

	// The React files are served
	router.Use(static.Serve("/", static.LocalFile("./ui/build", true)))
	// For manual url inputs, refreshes, page 404s, etc.
	// TODO: Figure out how to not apply this redirect for /api/ routes
	router.NoRoute(func(c *gin.Context) {
		c.File("./ui/build/index.html")
	})

	api := router.Group("/api/v1")
	{
		api.GET("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "pong",
			})
		})

		auth := api.Group("/auth/google")
		{
			auth.GET("/login", googleLoginHandler)
			auth.GET("/callback", googleCallbackHandler)
		}
	}

	// cd ui && npm run build && cd .. && go build cmd/main.go && main.exe
	router.Run(":9090")
}
