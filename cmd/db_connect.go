package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

var (
	dbHost = "aabgan6rulfuiy.cumylu5vzsok.eu-central-1.rds.amazonaws.com"
	dbPort = 5432
	dbUser = os.Getenv("SDP_SQL_USER")
	dbPass = os.Getenv("SDP_SQL_PASS")
	dbName = "sdp_data"
	dbObj  *sql.DB
)

type user struct {
	FirstName   *string `json:"first_name"`
	LastName    *string `json:"last_name"`
	SchoolEmail *string `json:"school_email"`
}

func init() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPass, dbName)
	var err error
	dbObj, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Panic(err)
	}
	if err = dbObj.Ping(); err != nil {
		log.Panic(err)
	}
}

func getUserID(googleID string) int32 {
	sqlStatement := `SELECT id FROM users WHERE google_id=$1;`
	row := dbObj.QueryRow(sqlStatement, googleID)
	var id int32
	switch err := row.Scan(&id); err {
	case sql.ErrNoRows:
		log.Println("User created.")
		id = createUser(googleID)
	case nil:
		log.Println("User found.")
	default:
		log.Panic(err)
	}
	return id
}

func createUser(googleID string) int32 {
	sqlStatement := `INSERT INTO users (google_id)
	VALUES ($1) RETURNING id`
	var id int32
	err := dbObj.QueryRow(sqlStatement, googleID).Scan(&id)
	if err != nil {
		log.Panic(err)
	}
	return id
}

func getUser(ctx *gin.Context) {
	sqlStatement := `SELECT first_name, last_name, school_email
	FROM users WHERE id=$1;`
	row := dbObj.QueryRow(sqlStatement, *getUserIDFromCookie(ctx))
	var userObj user
	if err := row.Scan(&userObj.FirstName, &userObj.LastName, &userObj.SchoolEmail); err != nil {
		log.Panic(err)
	}
	ctx.JSON(http.StatusOK, gin.H{
		"first_name":   userObj.FirstName,
		"last_name":    userObj.LastName,
		"school_email": userObj.SchoolEmail,
	})
}
