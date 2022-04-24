package middleware

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"on951/application"
	"on951/web"
)

func ApiAuthCheck(ctx *gin.Context) {

	ctx.Header("Vary", "Authorization")
	user, err := application.GetApplication().
		GetAuthorizedUserFromHeader(
			ctx.GetHeader("Authorization"),
		)
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
