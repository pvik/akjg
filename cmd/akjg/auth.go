package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	c "github.com/pvik/akjg/internal/config"
	"github.com/pvik/akjg/pkg/httphelper"

	"github.com/golang-jwt/jwt"
	log "github.com/sirupsen/logrus"
)

func apiAuth(w http.ResponseWriter, r *http.Request) {
	authHeaderArr, ok := r.Header["Authorization"]
	if !ok || len(authHeaderArr) < 1 || len(authHeaderArr[0]) < 1 {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if strings.HasPrefix(authHeaderArr[0], "Bearer ") {
		tokenString := strings.Split(authHeaderArr[0], " ")[1]

		// Parse takes the token string and a function for looking up the key
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Don't forget to validate the alg is what you expect:
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}

			return []byte(c.AppConf.JWTSecret), nil
		})
		if err != nil {
			log.Errorf("parse token err: %s", err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			// if JWT has exp claim, check if it has not expired
			//  jwt.Parse verifies token expiry

			httphelper.RespondwithJSON(w, 200, claims)
			return
		}
	}

	w.WriteHeader(http.StatusUnauthorized)
	return
}

func apiJWT(w http.ResponseWriter, r *http.Request) {
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
