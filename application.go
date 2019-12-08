package main

import (
	"strings"

	"github.com/TakingItCasual/sdp/gocode"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	// The React files are served
	router.Use(static.Serve("/", static.LocalFile("./client/build", true)))
	// For manual url inputs, refreshes, page 404s, etc.
	router.NoRoute(func(c *gin.Context) {
		path := c.Request.URL.Path
		// Invalid API paths get a backend 404 rather than a frontend 404
		if !strings.HasPrefix(path, "/api/") {
			c.File("./client/build/index.html")
		}
	})

	api := router.Group("/api/v1")
	{
		authorized := api.Group("/priv")
		authorized.Use(gocode.authMiddleware())
		{
			authorized.GET("/user", gocode.GetUser)
			authorized.PUT("/user", gocode.PutUser)
			authorized.GET("/users", gocode.GetUsers)
		}

		auth := api.Group("/auth")
		{
			googleAuth := auth.Group("/google")
			{
				googleAuth.GET("/login", gocode.GoogleLoginHandler)
				googleAuth.GET("/callback", gocode.GoogleCallbackHandler)
			}
		}
	}

	// cd client && npm run build && cd .. && go build ./application.go && application.exe
	router.Run(":5000")
}
