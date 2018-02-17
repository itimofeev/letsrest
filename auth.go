package letsrest

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/rs/xid"
	"time"
)

// LetsRestClaims claims in terminology of jwt just a data that serialized in jwt token
type LetsRestClaims struct {
	UserID string `json:"user_id"`
	jwt.StandardClaims
}

type Auth struct {
	UserID    string `json:"user_id"`
	AuthToken string `json:"auth_token"`
}

func createUser() *User {
	return &User{
		ID:           xid.New().String(),
		RequestLimit: 1000,
	}
}

func createAuthToken(user *User) *Auth {
	expDate := time.Now().Add(time.Hour * 24 * 355 * 10) // 10 years

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, LetsRestClaims{
		UserID: user.ID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expDate.Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "LetsRest",
		},
	})

	// Sign and get the complete encoded token as a string
	tokenString, err := token.SignedString([]byte(secretForJwt))
	Must(err, "token.SignedString([]byte(secretForJwt))")
	return &Auth{UserID: user.ID, AuthToken: tokenString}
}

func userFromAuthToken(authToken string) (*User, error) {
	token, err := jwt.ParseWithClaims(authToken, &LetsRestClaims{}, func(token *jwt.Token) (interface{}, error) {
		return secretForJwt, nil
	})

	if err != nil {
		return nil, err
	}

	claims, _ := token.Claims.(*LetsRestClaims)
	//if !ok || !token.Valid {
	//}
	return &User{ID: claims.UserID}, nil
}
