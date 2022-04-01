package middleware

import (
	"github.com/gin-gonic/gin"
	"log"
	"main/state"
	"main/web"
	"net/http"
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
	user, err := state.GetAuthorizedUserFromHeader(ctx.GetHeader("Authorization"))
	if err != nil {
		web.Write(ctx, http.StatusUnauthorized, err)
		return
	}

	if gin.Mode() == gin.DebugMode {
		log.Printf("Logged user: #%v %v", user.Id, user.Name)
	}
	ctx.Next()
}
