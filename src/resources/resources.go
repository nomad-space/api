package resources

import (
	"github.com/rs/zerolog"
)

type Resources struct {
	Config 	Specification
	Mongo 	MongoDB
	Log 	zerolog.Logger
	Mail 	Mail
}

var instanceResources *Resources

func GetInstance() (r *Resources, err error) {
	if instanceResources == nil {
		instanceResources = &Resources{}
		if err := instanceResources.initConfig(); err != nil {
			return instanceResources, err
		}
		if err := instanceResources.initLog(); err != nil {
			return instanceResources, err
		}
		if err := instanceResources.initMongo(); err != nil {
			return instanceResources, err
		}
		if err := instanceResources.initMail(); err != nil {
			return instanceResources, err
		}
	}
	return instanceResources, nil
}