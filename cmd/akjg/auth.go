package main

import (
	"net/http"
	"time"

	c "github.com/pvik/akjg/internal/config"
	"github.com/pvik/akjg/pkg/httphelper"

	"github.com/golang-jwt/jwt"
	log "github.com/sirupsen/logrus"
)

func apiLogin(w http.ResponseWriter, r *http.Request) {
	apikeyArr, ok := r.URL.Query()["apikey"]

	if !ok || len(apikeyArr) < 1 || len(apikeyArr[0]) < 1 {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	apiKey := apikeyArr[0]

	apiKeyJwtClaims, apiKeyConfigured := c.AppConf.APIKeyJWTDetailsMap[apiKey]

	if !apiKeyConfigured {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	jwtClaims := make(jwt.MapClaims)
	for k, v := range apiKeyJwtClaims {
		jwtClaims[k] = v
	}

	// populate iat and exp claim
	jwtClaims["iat"] = time.Now().UTC().Unix()
	jwtClaims["exp"] = time.Now().Add(time.Duration(c.AppConf.JWTExpiryMins) * time.Minute).UTC().Unix()

	// Create a new token object, specifying signing method and the claims
	// you would like it to contain.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtClaims)

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(c.AppConf.JWTSecret))
	if err != nil {
		log.Errorf("error generating JWT: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Respons with JWT Token
	httphelper.RespondwithJSON(w, 200, map[string]interface{}{
		"token": tokenString,
	})
}
