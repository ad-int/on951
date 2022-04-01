package handlers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/pascaldekloe/jwt"
	"golang.org/x/crypto/bcrypt"
	"log"
	dbStructure "main/database/structure"
	"main/state"
	"net/http"
	"time"
)

func GetToken(ctx *gin.Context) {

	user := ctx.DefaultQuery("user", "guest")
	password := ctx.DefaultQuery("password", "not-set")
	audience := ctx.DefaultQuery("audience", state.GetApplication().GetConfigValue("AUDIENCE"))
	issuer := ctx.DefaultQuery("issuer", state.GetApplication().GetConfigValue("ISSUER"))

	hash, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	userRecord, err := json.Marshal(dbStructure.User{
		Id:       1,
		Name:     user,
		Password: string(hash),
	})

	if err != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, err)
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
	jwtBytes, err = claims.HMACSign(jwt.HS512, []byte(state.GetApplication().GetConfigValue("SECRET")))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
	}

	ctx.Status(http.StatusOK)
	_, e := ctx.Writer.Write(jwtBytes)
	if e != nil {
		log.Println(e)
	}

}
