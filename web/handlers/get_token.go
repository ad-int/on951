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
	"on951/web"
	"strconv"
	"time"
)

func GetToken(ctx *gin.Context) {

	user := ctx.DefaultQuery("user", "guest")
	password := ctx.DefaultQuery("password", "not-set")
	audience := ctx.DefaultQuery("audience", application.GetApplication().GetConfigValue("AUDIENCE"))
	issuer := ctx.DefaultQuery("issuer", application.GetApplication().GetConfigValue("ISSUER"))

	cost, _ := strconv.Atoi(application.GetApplication().GetConfigValue("BCRYPT_HASH_GENERATION_COST"))
	log.Println("cost", cost)
	hash, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		web.WriteMessage(ctx, http.StatusInternalServerError, "unable to generate token", err)
		return
	}

	userRecord, err := json.Marshal(dbStructure.User{
		Id:       1,
		Name:     user,
		Password: string(hash),
	})

	if err != nil {
		web.WriteMessage(ctx, http.StatusInternalServerError, "internal server error", err)
		return
	}

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
	}

	ctx.Status(http.StatusOK)
	_, e := ctx.Writer.Write(jwtBytes)
	if e != nil {
		log.Println(e)
	}

}
