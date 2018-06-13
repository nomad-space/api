package models

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

type Booking struct {
	Id				bson.ObjectId	`bson:"_id,omitempty" json:"id" optional:"true"`
	FirstName		string			`json:"firstname" description:"firstname of the user" default:"john" valid:"required"`
	LastName 		string			`json:"lastname" description:"lastname of the user" default:"brown" valid:"required"`
	Phone 			string			`json:"phone" description:"phone of the user" default:"+1..."`
	Email		 	string			`json:"email" description:"email of the user" default:"example@domain.com" valid:"required"`
	LocationId		string			`bson:"location_id" json:"location_id" description:"location ID of the booking" valid:"required"`
	HotelId			string			`bson:"hotel_id" json:"hotel_id" description:"hotel ID of the booking" valid:"required"`
	GateId			string			`bson:"gate_id" json:"gate_id" description:"gate ID of the booking" valid:"required"`
	RoomId			string			`bson:"room_id" json:"room_id" description:"room ID of the booking" valid:"required"`
	CheckIn			string			`bson:"checkin_date" json:"checkin_date" description:"checkin of the booking" valid:"required"`
	CheckOut		string			`bson:"checkout_date" json:"checkout_date" description:"checkout of the booking" valid:"required"`
	Adults			string			`bson:"adults" json:"adults" description:"adults of the booking" valid:"required"`
	CreatedAt		time.Time		`bson:"created_at" json:"created_at" optional:"true"`
	UpdatedAt		time.Time		`bson:"updated_at" json:"updated_at" optional:"true"`
}