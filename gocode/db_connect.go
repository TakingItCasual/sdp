package gocode

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
	dbHost = "localhost"
	dbPort = 5432
	dbUser = os.Getenv("SDP_SQL_USER")
	dbPass = os.Getenv("SDP_SQL_PASS")
	dbName = "postgres"
	dbObj  *sql.DB
)

type user struct {
	FirstName   *string `json:"first_name" binding:"required"`
	LastName    *string `json:"last_name" binding:"required"`
	SchoolEmail *string `json:"school_email" binding:"required"`
}

func panicIfErr(err error) {
	if err != nil {
		log.Panic(err)
	}
}

func init() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPass, dbName)
	var err error
	dbObj, err = sql.Open("postgres", psqlInfo)
	panicIfErr(err)
	err = dbObj.Ping()
	panicIfErr(err)
	createTable()
}

func createTable() {
	sqlStatement := `CREATE TABLE IF NOT EXISTS users
	(
		id serial PRIMARY KEY,
		google_id VARCHAR (50) NOT NULL,
		first_name VARCHAR (50),
		last_name VARCHAR (50),
		school_email VARCHAR (50),
	)`
	_, err := dbObj.Exec(sqlStatement)
	panicIfErr(err)
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
	panicIfErr(err)
	return id
}

// GetUser used as route
func GetUser(ctx *gin.Context) {
	sqlStatement := `SELECT first_name, last_name, school_email
	FROM users WHERE id=$1;`
	row := dbObj.QueryRow(sqlStatement, *GetUserIDFromCookie(ctx))
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

// GetUsers used as route
func GetUsers(ctx *gin.Context) {
	rows, err := dbObj.Query("SELECT first_name, last_name, school_email FROM users")
	panicIfErr(err)
	defer rows.Close()
	userObjs := []user{}
	for rows.Next() {
		var userObj user
		err = rows.Scan(&userObj.FirstName, &userObj.LastName, &userObj.SchoolEmail)
		panicIfErr(err)
		userObjs = append(userObjs, userObj)
	}
	err = rows.Err()
	panicIfErr(err)

	ctx.JSON(http.StatusOK, userObjs)
}

// PutUser used as route
func PutUser(ctx *gin.Context) {
	var form user
	if err := ctx.ShouldBindJSON(&form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if *form.FirstName == "" || *form.LastName == "" || *form.SchoolEmail == "" {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	sqlStatement := `UPDATE users
	SET first_name = $2, last_name = $3, school_email = $4
	WHERE id = $1;`
	_, err := dbObj.Exec(sqlStatement, *GetUserIDFromCookie(ctx), *form.FirstName, *form.LastName, *form.SchoolEmail)
	panicIfErr(err)
}
