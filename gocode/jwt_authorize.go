package gocode

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var jwtSecretKey []byte

func init() {
	jwtSecretKey = []byte(os.Getenv("SDP_GOOGLE_CLIENT_SECRET"))
}

type customClaims struct {
	UserID    int32 `json:"id"`
	RefreshAt int64 `json:"refresh"`
	jwt.StandardClaims
}

func createJWT(ctx *gin.Context, id int32) {
	refreshTimer, refreshExpire := genExpires()
	claims := customClaims{
		id,
		refreshTimer.Unix(),
		jwt.StandardClaims{
			ExpiresAt: refreshExpire.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(jwtSecretKey)
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	http.SetCookie(ctx.Writer, &http.Cookie{
		Name:    "auth_token",
		Value:   tokenString,
		Path:    "/",
		Expires: refreshExpire, // Same as token's
	})
}

func getUserIDFromCookie(ctx *gin.Context) *int32 {
	tokenString, err := ctx.Cookie("auth_token")
	if err != nil {
		log.Printf("JWT: Error with cookie: %v\n", err.Error())
		return nil
	}

	token, err := jwt.ParseWithClaims(tokenString, &customClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return jwtSecretKey, nil
		},
	)
	if err != nil {
		log.Printf("JWT: Error parsing claim: %v\n", err.Error())
		return nil
	}

	// Signing method doesn't match
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		log.Printf("JWT: Signing method mismatch: %v\n", err.Error())
		return nil
	}
	if claims, ok := token.Claims.(*customClaims); ok && token.Valid {
		if time.Unix(claims.RefreshAt, 0).Sub(time.Now()) < 0 {
			refreshJWT(ctx, *claims)
		}
		return &claims.UserID
	}
	log.Println("Token not valid.")
	return nil
}

func refreshJWT(ctx *gin.Context, oldClaims customClaims) {
	createJWT(ctx, oldClaims.UserID)
}

func genExpires() (time.Time, time.Time) {
	shortExpire := time.Now().Add(15 * time.Minute)  // JWT must be refreshed
	longExpire := time.Now().Add(7 * 24 * time.Hour) // User must re-login
	return shortExpire, longExpire
}
