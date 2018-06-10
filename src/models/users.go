package models

import (
	"time"
	"errors"
	"math/rand"

	"nomad/api/src/resources"
	"golang.org/x/crypto/bcrypt"
	"github.com/dgrijalva/jwt-go"
	"gopkg.in/mgo.v2/bson"
)

type Users struct {
	Id				bson.ObjectId	`bson:"_id,omitempty" json:"-" optional:"true"`
	Status			int				`json:"status" description:"status of the user" default:"0" valid:"required"`
	FirstName		string			`json:"firstname" description:"firstname of the user" default:"john" valid:"required"`
	LastName 		string			`json:"lastname" description:"lastname of the user" default:"brown" valid:"required"`
	Phone 			string			`json:"phone" description:"phone of the user" default:"+1..."`
	Email		 	string			`json:"email" description:"email of the user" default:"example@domain.com" valid:"required"`
	Password		string			`bson:"-" json:"password" description:"password of the user" valid:"required"`
	HashedPassword	[]byte			`bson:"password" json:"-"`
	ConfirmToken	string			`bson:"confirm_token" json:"-"`
	CreatedAt		time.Time		`bson:"created_at" json:"-" optional:"true"`
	UpdatedAt		time.Time		`bson:"updated_at" json:"-" optional:"true"`
}

const USER_STATUS_NEW = 0
const USER_STATUS_ACTIVE = 1
const USER_STATUS_BLOCK = 2

func (u *Users) GenerateHashPassword() ([]byte, error) {

	if u.Password == "" {
		return nil, errors.New("Waiting password")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	u.HashedPassword = hashedPassword

	return hashedPassword, err
}

func (u *Users) GenerateConfirmToken() (string) {

	rand.Seed(time.Now().UnixNano())

	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	b := make([]rune, 32)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	u.ConfirmToken = string(b)

	return string(b)
}

func (u *Users) CompareHashAndPassword(password string) (error) {

	if password == "" {
		return errors.New("Waiting password")
	}

	return bcrypt.CompareHashAndPassword(u.HashedPassword, []byte(password))
}

func (u *Users) GetJWT() (string, error) {

	res, _ := resources.GetInstance()

	mySigningKey := []byte(res.Config.JwtSecret)
	timeout := res.Config.JwtTimeout

	claims := JwtClaims{
		u.Id,
		u.FirstName,
		u.LastName,
		time.Now().Unix(),
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(timeout).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(mySigningKey)
	if err != nil {
		return "", errors.New("Waiting password")
	}

	return tokenString, nil
}
