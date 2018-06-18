package models

import (
	"time"
)

type Location struct {
	CountryCode			string			`json:"country_code" description:"country_code of the location" valid:"required"`
	Country				string			`json:"country" description:"country of the location" valid:"required"`
	LatinFullName		string			`json:"latin_full_name" description:"latin_full_name of the location" valid:"required"`
	Fullname			string			`json:"full_name" description:"full_name of the location" valid:"required"`
	Clar				string			`json:"clar" description:"clar of the location" valid:"required"`
	LatinClar			string			`json:"latin_clar" description:"latin_clar of the location" valid:"required"`
	City				string			`json:"city" description:"city of the location" valid:"required"`
	LatinCity			string			`json:"latin_city" description:"latin_city of the location" valid:"required"`
	Timezone			string			`json:"timezone" description:"timezone of the location" valid:"required"`
	Timezonesec			int				`json:"timezonesec" description:"timezonesec of the location" valid:"required"`
	LatinCountry		string			`json:"latin_country" description:"latin_country of the location" valid:"required"`
	Id					int				`json:"id" description:"id of the location" valid:"required"`
	CountryId			int				`json:"country_id" description:"country_id of the location" valid:"required"`
	CreatedAt			time.Time		`bson:"created_at" json:"created_at" optional:"true"`
	UpdatedAt			time.Time		`bson:"updated_at" json:"updated_at" optional:"true"`
}