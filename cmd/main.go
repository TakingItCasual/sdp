package main

import (
	"log"
	"os"
	"strings"

	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	googleOauthConf *oauth2.Config
	jwtSecretKey    []byte
)

func init() {
	googleOauthConf = &oauth2.Config{
		ClientID:     os.Getenv("SDP_GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("SDP_GOOGLE_CLIENT_SECRET"),
		RedirectURL:  "http://127.0.0.1:9090/api/v1/auth/google/callback",
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}

	jwtSecretKey = []byte(os.Getenv("SDP_GOOGLE_CLIENT_SECRET"))
}

func authMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if getUserIDFromCookie(ctx) == nil {
			log.Println("JWT auth failed.")
			ctx.AbortWithStatus(401)
		}
		ctx.Next()
	}
}

func main() {
	router := gin.Default()

	// CORS is enabled to allow the backend and frontend to communicate
	//config := cors.DefaultConfig()
	//config.AllowOrigins = []string{"http://localhost:3000"}
	//router.Use(cors.New(config))

	// The React files are served
	router.Use(static.Serve("/", static.LocalFile("./ui/build", true)))
	// For manual url inputs, refreshes, page 404s, etc.
	router.NoRoute(func(c *gin.Context) {
		path := c.Request.URL.Path
		// Invalid API paths get a backend 404 rather than a frontend 404
		if !strings.HasPrefix(path, "/api/") {
			c.File("./ui/build/index.html")
		}
	})

	api := router.Group("/api/v1")
	{
		authorized := api.Group("/priv")
		authorized.Use(authMiddleware())
		{
			authorized.GET("/user", getUser)
		}

		auth := api.Group("/auth")
		{
			googleAuth := auth.Group("/google")
			{
				googleAuth.GET("/login", googleLoginHandler)
				googleAuth.GET("/callback", googleCallbackHandler)
			}
		}
	}

	// cd ui && npm run build && cd .. && go build ./cmd/... && cmd.exe
	router.Run(":9090")
}
