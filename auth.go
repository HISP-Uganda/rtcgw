package main

import (
	"encoding/base64"
	"github.com/gin-gonic/gin"
	"rtcgw/db"
	"rtcgw/models"
	"strings"
)

func BasicAuth() gin.HandlerFunc {

	return func(c *gin.Context) {
		c.Set("dbConn", db.GetDB())
		c.Set("asynqClient", client)
		auth := strings.SplitN(c.Request.Header.Get("Authorization"), " ", 2)

		if len(auth) != 2 || (auth[0] != "Basic" && auth[0] != "Token:") {
			RespondWithError(401, "Unauthorized", c)
			return
		}
		tokenAuthenticated, userUID := AuthenticateUserToken(auth[1])
		if auth[0] == "Token:" {
			if !tokenAuthenticated {
				RespondWithError(401, "Unauthorized", c)
				return
			}
			c.Set("currentUser", userUID)
			c.Next()
			return
		}

		payload, _ := base64.StdEncoding.DecodeString(auth[1])
		pair := strings.SplitN(string(payload), ":", 2)

		basicAuthenticated, userUID := AuthenticateUser(pair[0], pair[1])

		if len(pair) != 2 || !basicAuthenticated {
			RespondWithError(401, "Unauthorized", c)
			// c.Writer.Header().Set("WWW-Authenticate", "Basic realm=Restricted")
			return
		}
		c.Set("currentUser", userUID)

		c.Next()
	}
}

func AuthenticateUser(username, password string) (bool, int64) {
	// log.Printf("Username:%s, password:%s", username, password)
	userObj := models.User{}
	err := db.GetDB().QueryRowx(
		`SELECT
            id, uid, username, firstname, lastname , telephone, email
        FROM users
        WHERE
            username = $1 AND password = crypt($2, password)`,
		username, password).StructScan(&userObj)
	if err != nil {
		// fmt.Printf("User:[%v]", err)
		return false, 0
	}
	// fmt.Printf("User:[%v]", userObj)
	return true, userObj.ID
}

func AuthenticateUserToken(token string) (bool, int64) {
	userToken := models.UserToken{}
	err := db.GetDB().QueryRowx(
		`SELECT
            id, user_id, token, is_active
        FROM user_apitoken
        WHERE
            token = $1 AND is_active = TRUE LIMIT 1`,
		token).StructScan(&userToken)
	if err != nil {
		return false, 0
	}
	// fmt.Printf("User:[%v]", userObj)
	return true, userToken.UserID
}

func RespondWithError(code int, message string, c *gin.Context) {
	resp := map[string]string{"error": message}

	c.JSON(code, resp)
	c.Abort()
}
