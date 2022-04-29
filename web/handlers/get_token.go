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
	"strconv"
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

	cost, _ := strconv.Atoi(application.GetApplication().GetConfigValue("BCRYPT_HASH_GENERATION_COST"))
	log.Println("cost", cost)
	hash, err := bcrypt.GenerateFromPassword([]byte(authTokenRequest.Password), cost)
	if err != nil {
		web.WriteMessage(ctx, http.StatusInternalServerError, "unable to generate token", err)
		return
	}
	if gin.Mode() == gin.DebugMode {
		u2 := dbStructure.User{
			Name:     authTokenRequest.Username,
			Password: string(hash),
		}
		application.GetApplication().GetDatabase().
			GetDB().Create(&u2)
	}
	user := dbStructure.User{}
	tx := application.GetApplication().GetDatabase().
		GetDB().
		Debug().
		Where("name LIKE ?", authTokenRequest.Username).
		First(&user)

	if tx.RowsAffected == 0 {
		user.Name = authTokenRequest.Username
		user.Password = string(hash)
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
