package token

import (
	"errors"
	"fmt"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	uuid "github.com/satori/go.uuid"
)

var ISSUER = "blueharbortech.com"
var EXPIRES = 60 * 60 * 3
var REFRESHEXPIRES = 60 * 60 * 24 * 30

type JWTCollection struct {
	AccessToken        *string
	AccessTokenClaims  oa2claims
	RefreshToken       *string
	RefreshTokenClaims oa2claims
}

type oa2claims struct {
	jwt.StandardClaims
	IdentityProvider uuid.UUID `json:"ip"`
	App              uuid.UUID `json:"app"`
}

func NewJWTCollection(key string, accessToken *string, refreshToken *string) (JWTCollection, error) {
	collection := JWTCollection{}

	if accessToken != nil {
		collection.AccessToken = accessToken
	}
	if refreshToken != nil {
		collection.RefreshToken = refreshToken
	}
	if ok, err := collection.ValidateTokens(key); !ok || err != nil {
		if err != nil {
			return JWTCollection{}, err
		}
		return JWTCollection{}, errors.New("A Token in the collection is not valid")
	}
	return collection, nil
}

func CreateJWT(key string, subject uuid.UUID, identityProvider uuid.UUID, app uuid.UUID) (*JWTCollection, error) {

	now := time.Now().Unix()

	token, err := createToken(key, oa2claims{
		StandardClaims: jwt.StandardClaims{
			Issuer:    ISSUER,
			ExpiresAt: now + int64(EXPIRES),
			Subject:   subject.String(),
		},
		IdentityProvider: identityProvider,
		App:              app,
	})
	if err != nil {
		return nil, err
	}

	refreshToken, err := createToken(key, oa2claims{
		StandardClaims: jwt.StandardClaims{
			Issuer:    ISSUER,
			ExpiresAt: now + int64(REFRESHEXPIRES),
			Subject:   subject.String(),
		},
		IdentityProvider: identityProvider,
		App:              app,
	})
	if err != nil {
		return nil, err
	}

	return &JWTCollection{
		AccessToken:  &token,
		RefreshToken: &refreshToken,
	}, nil
}

func createToken(key string, claims jwt.Claims) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(key))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func (c *JWTCollection) ValidateTokens(key string) (bool, error) {
	//Parse the Token

	if c.AccessToken != nil {
		claims := oa2claims{}
		token, err := c.validateToken(key, *c.AccessToken, &claims)
		if err != nil {
			return false, err
		}
		//TODO Add Validation of Claims, once we have custom claims for access tokens
		if !token.Valid {
			return false, errors.New("Invalid Token")
		}
		c.AccessTokenClaims = claims
	}
	if c.RefreshToken != nil {
		claims := oa2claims{}
		token, err := c.validateToken(key, *c.RefreshToken, &claims)
		if err != nil {
			return false, err
		}
		//TODO Add Validation of Claims, once we have custom claims for refresh tokens
		if !token.Valid {
			return false, errors.New("Invalid Token")
		}
		c.RefreshTokenClaims = claims
	}
	return true, nil
}

func (c *JWTCollection) validateToken(key string, token string, claims jwt.Claims) (*jwt.Token, error) {
	//Parse the Token
	tkn, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(key), nil
	})

	if ve, ok := err.(*jwt.ValidationError); ok {
		if ve.Errors&jwt.ValidationErrorMalformed != 0 {
			return nil, errors.New("Invalid Token")
		} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
			return nil, errors.New("Token has Expired")
		} else if ve.Errors&jwt.ValidationErrorSignatureInvalid != 0 {
			return nil, errors.New("Invalid access token")
		}
		return nil, err
	}
	return tkn, nil
}
