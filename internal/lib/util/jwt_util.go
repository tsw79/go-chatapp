package util

import (
	"fmt"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// Create the JWT key used to create the signature
var JWTSecret = []byte(os.Getenv("HMAC_SECRET"))

// Declare the expiration time of the token.
// We have kept it as 5 minutes
var ExpirationTime = time.Now().AddDate(0, 0, 1)

// Create a struct that will be encoded to a JWT.
// We add jwt.StandardClaims as an embedded type, to provide fields like expiry time
type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

// CreateJWTToken generates a JWT signed token for for the given user
func CreateJWTAccessToken(username string) (string, error) {
	// Create the JWT claims, which includes the username and expiry time
	claims := &Claims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: ExpirationTime.Unix(),
		},
	}
	// Declare the token with the algorithm used for signing, and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Create the JWT string
	tokenStr, err := token.SignedString(JWTSecret)
	if err != nil {
		return "", err
	}
	return tokenStr, nil
}

// This function will return an error if the token is invalid, i.e.
// the auth cookie has expired -or- the signature does not match.
func ValidateToken(token string) (Claims, error) {
	// Parse the JWT token string and store the result in `Claims`.
	jwtToken, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return JWTSecret, nil
	})
	// check error
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			// rw.WriteHeader(http.StatusUnauthorized)
			fmt.Println("ValidateToken: Unauthorized!")
			return Claims{}, err
		}
		// rw.WriteHeader(http.StatusBadRequest)
		fmt.Println("ValidateToken - Bad Request")
		return Claims{}, err
	}
	// check token's validity
	if !jwtToken.Valid {
		// rw.WriteHeader(http.StatusUnauthorized)
		fmt.Println("ValidateToken - Status Unauthorized")
		return Claims{}, fmt.Errorf("Status Unauthorized - JWT token not valid!")
	}
	claims, ok := jwtToken.Claims.(*Claims)
	if !ok {
		fmt.Println(claims)
		return Claims{}, fmt.Errorf("Claims not ok!")
		// handle not ok claims
	}
	return *claims, nil
}

// Renew an existing token
func RenewToken(username string) (string, error) {
	return CreateJWTAccessToken(username)
}

// Checks if given token is within `n` seconds of expiry.
// 	`n` = seconds
func TokenIsWithinExpiry(claims *Claims) bool {
	var n time.Duration
	n = 30 // 30 seconds within expiry time
	if time.Unix(claims.ExpiresAt, 0).Sub(time.Now()) > n*time.Second {
		// fmt.Println(http.StatusBadRequest)
		return true
	}
	return false
}
