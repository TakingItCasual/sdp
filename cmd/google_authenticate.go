package main

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
)

type googleUser struct {
	Sub     string `json:"sub"`
	Picture string `json:"picture"`
	Email   string `json:"email"`
}

func googleLoginHandler(ctx *gin.Context) {
	oauthState := generateStateOauthCookie(ctx)
	loginURL := googleOauthConf.AuthCodeURL(oauthState)
	ctx.JSON(http.StatusOK, gin.H{"redirect": loginURL})
}

func googleCallbackHandler(ctx *gin.Context) {
	oauthState, _ := ctx.Cookie("googleOauthstate")
	// Confirm cookie and callback states are the same (prevents CSRF attacks)
	if ctx.Query("state") != oauthState {
		log.Println("invalid oauth google state")
		ctx.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}

	data, err := getGoogleUserData(ctx)
	if err != nil {
		log.Println(err.Error())
		ctx.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}

	userID := getUserID(data.Sub)
	log.Printf("User ID: %d\n", userID)

	createJWT(ctx, userID)
	ctx.Redirect(http.StatusTemporaryRedirect, "/profile")
}

func generateStateOauthCookie(ctx *gin.Context) string {
	b := make([]byte, 32)
	rand.Read(b)
	randData := base64.URLEncoding.EncodeToString(b)
	cookie := http.Cookie{
		Name:     "googleOauthstate",
		Value:    randData,
		Expires:  time.Now().Add(time.Hour), // Must log in in under an hour
		HttpOnly: true,                      // Cookie only accessible from backend
	}
	http.SetCookie(ctx.Writer, &cookie)

	return randData
}

func getGoogleUserData(ctx *gin.Context) (*googleUser, error) {
	// Use code to get token and get user info from Google.
	tok, err := googleOauthConf.Exchange(oauth2.NoContext, ctx.Query("code"))
	if err != nil {
		return nil, err
	}

	client := googleOauthConf.Client(oauth2.NoContext, tok)
	data, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		return nil, err
	}

	defer data.Body.Close()
	userData := &googleUser{}
	json.NewDecoder(data.Body).Decode(userData)
	ctx.Status(http.StatusOK)

	return userData, nil
}
