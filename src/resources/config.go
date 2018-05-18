package resources

import (
	"log"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Specification struct {
	Debug       	bool
	Port        	int					`default:"7784"`
	User        	string
	Users       	[]string
	Rate        	float32
	Timeout     	time.Duration
	ColorCodes  	map[string]int
	MongoHost  		string			 	`default:"localhost:27017"`
	//MongoHost  		string			 	`default:"db.mvp.nomad.space:27017"`
	MongoDB  		string			 	`default:"db_name"`
	MongoCollUsers  string			 	`default:"users"`
	JwtSecret  		string			 	`default:"jwt_secret"`
	JwtTimeout  	time.Duration	 	`default:"24h""`
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
