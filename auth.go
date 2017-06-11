package letsrest

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/nu7hatch/gouuid"
	"time"
)

// LetsRestClaims claims in terminology of jwt just a data that serialized in jwt token
type LetsRestClaims struct {
	UserID string
	jwt.StandardClaims
}

type Auth struct {
	AuthToken string `json:"auth_token"`
}

func createAuthToken() *Auth {
	userID, err := uuid.NewV4()
	Must(err, "uuid.NewV4()")
	expDate := time.Now().Add(time.Hour * 24 * 355 * 10) // 10 years

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, LetsRestClaims{
		UserID: userID.String(),
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expDate.Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "LetsRest",
		},
	})

	// Sign and get the complete encoded token as a string
	tokenString, err := token.SignedString([]byte(secretForJwt))
	Must(err, "token.SignedString([]byte(secretForJwt))")
	return &Auth{AuthToken: tokenString}
}
