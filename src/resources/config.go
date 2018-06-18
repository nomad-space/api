package resources

import (
	"log"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Specification struct {
	Debug       		bool				`default:"true"`
	Port        		int					`default:"7784"`
	FrontURL  			string			 	`default:"http://mvp.nomad.space"`
	MongoHost  			string			 	`default:"127.0.0.1:27017"`
	MongoDB  			string			 	`default:"db_name"`
	MongoCollUsers  	string			 	`default:"users"`
	MongoCollBookings	string			 	`default:"bookings"`
	MongoCollHotels		string			 	`default:"hotels"`
	MongoCollLocations	string			 	`default:"locations"`
	JwtSecret  			string			 	`default:"jwt_secret"`
	JwtTimeout  		time.Duration	 	`default:"24h"`
	SmtpLogin			string				`default:"no-reply@mailman.nomad.space"`
	SmtpPassword		string				`default:"password"`
	SmtpHost			string				`default:"smtp.yandex.ru"`
	SmtpPort			string				`default:"25"`
	SendmailFrom		string				`default:"no-reply@mailman.nomad.space"`
}

func (r *Resources) initConfig() error {
	var s Specification

	err := envconfig.Process("myapp", &s)
	if err != nil {
		log.Fatal(err.Error())
		return err
	}

	instanceResources.Config = s

	return nil
}
