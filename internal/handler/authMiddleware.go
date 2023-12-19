package handler

import (
    "encoding/base64"
    "net/http"
    "strings"

    "github.com/gin-gonic/gin"
)

func BasicAuth() gin.HandlerFunc {
    return func(ctx *gin.Context) {
        // Define the authorized user
        user := "user"
        pass := "aymanfaik1"

        // Get the Authorization header
        auth := ctx.Request.Header.Get("Authorization")

        // Check if the Authorization header is empty
        if auth == "" {
            ctx.Header("WWW-Authenticate", `Basic realm="Restricted"`)
            ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is missing"})
            return
        }

        // Decode the Authorization header
        b64auth := strings.SplitN(auth, " ", 2)
        if len(b64auth) != 2 || b64auth[0] != "Basic" {
            ctx.Header("WWW-Authenticate", `Basic realm="Restricted"`)
            ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization header"})
            return
        }

        // Decode the base64 encoded username:password
        payload, _ := base64.StdEncoding.DecodeString(b64auth[1])
        pair := strings.SplitN(string(payload), ":", 2)

        // Check if the username and password match
        if len(pair) != 2 || pair[0] != user || pair[1] != pass {
            ctx.Header("WWW-Authenticate", `Basic realm="Restricted"`)
            ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
            return
        }

        // If everything is ok, continue to the next middleware or handler
        ctx.Next()
    }
}