package main

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func SetInstaBotOptions(c *gin.Context) {

}

func authenticatorMiddleware(c *gin.Context) {
	// Read the authorization header from the request
	authHeader := c.GetHeader("Authorization")

	// Check if the authorization header is present
	if authHeader == "" {
		// If the authorization header is not present, return a 401 Unauthorized response
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	// Split the authorization header into its components
	authComponents := strings.Split(authHeader, " ")
	if len(authComponents) != 3 || authComponents[0] != "Credentials" {
		// If the authorization header is not in the correct format, return a 401 Unauthorized response
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	// Validate the access token
	// This may involve making a request to an OAuth server or using a local token store
	if authComponents[1] != *password && authComponents[2] != *username {
		// If the token is invalid, return a 401 Unauthorized response
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	// If the credential is valid, allow the request to be forwarded to the server
	c.Next()
}

func RedirectToHttps(w http.ResponseWriter, r *http.Request) {
	// Redirect the incoming HTTP request.
	host := strings.Split(r.Host, ":")[0]
	log.Println(r.Host)
	http.Redirect(w, r, "https://"+host+r.RequestURI, http.StatusMovedPermanently)
}

func SetBotConfigs(c *gin.Context) {
	var localbotconfigs BotConfigs
	if err := c.ShouldBindJSON(&localbotconfigs); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Set the bot configs

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func GetBotConfigs(c *gin.Context) {
	// Get the bot configs
	c.JSON(http.StatusOK, botconfigs)
}
