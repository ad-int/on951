package middleware

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"on951/application"
	"on951/web"
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
	user, err := application.GetAuthorizedUserFromHeader(ctx.GetHeader("Authorization"))
	if err != nil {
		web.Write(ctx, http.StatusUnauthorized, err)
		ctx.Abort()
		return
	}

	if gin.Mode() == gin.DebugMode {
		log.Printf("Logged user: #%v %v", user.Id, user.Name)
	}
	ctx.Next()
}
