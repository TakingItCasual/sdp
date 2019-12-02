package main

import (
	"crypto/rand"
	"encoding/base64"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
)

func googleLoginHandler(c *gin.Context) {
	oauthState := generateStateOauthCookie(c)
	loginURL := googleOauthConf.AuthCodeURL(oauthState)
	c.JSON(http.StatusOK, gin.H{"redirect": loginURL})
}

func googleCallbackHandler(c *gin.Context) {
	oauthState, _ := c.Cookie("googleOauthstate")
	// Confirm cookie and callback states are the same (prevents CSRF attacks)
	if c.Query("state") != oauthState {
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
	log.Printf("UserInfo: %s\n", data)
	c.Redirect(http.StatusTemporaryRedirect, "/profile")
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

func getGoogleUserData(c *gin.Context) ([]byte, error) {
	// Use code to get token and get user info from Google.
	tok, err := googleOauthConf.Exchange(oauth2.NoContext, c.Query("code"))
	if err != nil {
		return nil, err
	}

	client := googleOauthConf.Client(oauth2.NoContext, tok)
	userInfo, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		return nil, err
	}

	defer userInfo.Body.Close()
	data, _ := ioutil.ReadAll(userInfo.Body)
	c.Status(http.StatusOK)

	return data, nil
}
