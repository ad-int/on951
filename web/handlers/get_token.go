package handlers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/pascaldekloe/jwt"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"on951/application"
	dbStructure "on951/database/structure"
	"on951/models"
	"on951/web"
	"strings"
	"time"
)

func GetToken(ctx *gin.Context) {

	authTokenRequest := &models.AuthTokenRequest{}
	err := ctx.BindJSON(authTokenRequest)
	if err != nil {
		web.WriteBadRequestError(ctx, "missing or invalid body", err)
		return
	}

	if strings.TrimSpace(authTokenRequest.Username) == "" || strings.TrimSpace(authTokenRequest.Password) == "" {
		web.WriteBadRequestError(ctx, "missing credentials")
		return
	}
	log.Println(authTokenRequest)
	audience := ctx.DefaultQuery("audience", application.GetApplication().GetConfigValue("AUDIENCE"))
	issuer := ctx.DefaultQuery("issuer", application.GetApplication().GetConfigValue("ISSUER"))

	user := dbStructure.User{}
	if db := application.GetApplication().GetDatabase(); db != nil {
		application.GetApplication().GetDatabase().
			GetDB().
			Debug().
			Where(dbStructure.User{Name: authTokenRequest.Username}).
			Find(&user)
	}

	e := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(authTokenRequest.Password))
	if e != nil {
		web.WriteMessage(ctx, http.StatusUnauthorized, "invalid login")
		return
	}
	userRecord, _ := json.Marshal(user)

	var claims jwt.Claims
	claims.Subject = string(userRecord)
	claims.Issued = jwt.NewNumericTime(time.Now())
	claims.NotBefore = jwt.NewNumericTime(time.Now())
	claims.Expires = jwt.NewNumericTime(time.Now().Add(24 * time.Hour))
	claims.Issuer = issuer
	claims.Audiences = []string{audience}

	var jwtBytes []byte
	jwtBytes, err = claims.HMACSign(jwt.HS512, []byte(application.GetApplication().GetConfigValue("SECRET")))
	if err != nil {
		web.WriteMessage(ctx, http.StatusInternalServerError, "internal server error", err)
		return
	}

	authTokenResponse := models.AuthTokenResponse{}
	authTokenResponse.AccessToken = string(jwtBytes)
	authTokenResponse.TokenType = "Bearer"
	authTokenResponse.ExpiresIn = 3600
	ctx.IndentedJSON(http.StatusOK, authTokenResponse)
}
