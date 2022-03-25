package middleware

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/pascaldekloe/jwt"
	"log"
	"main/models"
	"main/state"
	"net/http"
	"strings"
	"time"
)

const MsgUnauthorized = "unauthorized"
const MsgNotAcceptedAudience = "not accepted audience"
const MsgInvalidIssuer = "invalid issuer"


func EnableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Allow-Access-Origin", "*")
		next.ServeHTTP(w, r)
	})
}

func ApiAuthCheck(ctx *gin.Context) {

	ctx.Header("Vary", "Authorization")
	authHeader := ctx.GetHeader("Authorization")
	if authHeader == "" {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}
	authHeaderParts := strings.Split(authHeader, " ")
	if len(authHeaderParts) != 2 {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	if authHeaderParts[0] != "Bearer" {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}
	token := authHeaderParts[1]
	claim, err := jwt.HMACCheck([]byte(token), []byte(state.GetApplication().GetConfigValue("SECRET")))
	if err != nil {
		ctx.AbortWithError(http.StatusUnauthorized, err)
		return
	}
	if !claim.Valid(time.Now()) {
		ctx.AbortWithError(http.StatusUnauthorized, errors.New(MsgUnauthorized))
		return
	}

	if !claim.AcceptAudience(state.GetApplication().GetConfigValue("AUDIENCE")) {
		ctx.AbortWithError(http.StatusUnauthorized, errors.New(MsgNotAcceptedAudience))
		return
	}

	if claim.Issuer != state.GetApplication().GetConfigValue("ISSUER") {
		ctx.AbortWithError(http.StatusUnauthorized, errors.New(MsgInvalidIssuer))
		return
	}
	var user models.User
	err = json.Unmarshal([]byte(claim.Subject), &user)
	if err != nil {
		ctx.AbortWithError(http.StatusUnauthorized, err)
		return
	}
	log.Println(user)
	ctx.Next()
}
