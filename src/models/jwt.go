package models

import (
	"errors"
	"time"

	"gopkg.in/mgo.v2/bson"
	"github.com/dgrijalva/jwt-go"
	"nomad/api/src/resources"
)

type JwtClaims struct {
	Id					bson.ObjectId		`json:"id"`
	FirstName			string				`json:"firstname"`
	LastName			string				`json:"lastname"`
	OrigIat				int64				`json:"orig_iat"`
	jwt.StandardClaims
}

func UpdateJWTToken(tokenString string) (string, error) {

	res, _ := resources.GetInstance()

	mySigningKey := []byte(res.Config.JwtSecret)
	timeout := res.Config.JwtTimeout

	claims, err := ParseJWTToken(tokenString)
	if err != nil {
		return "", err
	}

	claims.ExpiresAt = time.Now().Add(timeout).Unix()
	claims.OrigIat = time.Now().Unix()

	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	newTokenString, err := newToken.SignedString(mySigningKey)
	if err != nil {
		return "", errors.New("Waiting password")
	}

	return newTokenString, nil
}
func (jwt *JwtClaims) GetModel() (*Users) {

	res, _ := resources.GetInstance()

	user := Users{}

	collection, session, err := res.Mongo.UserCollectionAndSession();
	if err != nil {
		return nil
	}
	defer session.Close()
	err = collection.Find(bson.M{"_id": jwt.Id}).One(&user)
	if err != nil {
		return nil
	}
	return &user
}

func ParseJWTToken(tokenString string) (*JwtClaims, error) {

	res, _ := resources.GetInstance()

	mySigningKey := []byte(res.Config.JwtSecret)

	token, err := jwt.ParseWithClaims(tokenString, &JwtClaims{}, func(token *jwt.Token) (interface{}, error) {
		return mySigningKey, nil
	})

	if token.Valid {
		if claims, ok := token.Claims.(*JwtClaims); ok{
			return claims, nil
		} else {
			return nil, err
		}
	} else if ve, ok := err.(*jwt.ValidationError); ok {
		if ve.Errors&jwt.ValidationErrorMalformed != 0 {
			return nil, errors.New("That's not even a token")
		} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
			return nil, errors.New("Timing is everything")
		} else {
			return nil, errors.New("Couldn't handle this token: "+err.Error())
		}
	} else {
		return nil, errors.New("Couldn't handle this token: "+err.Error())
	}
}
