package resources

import (
	"errors"
	"gopkg.in/mgo.v2"
)

type MongoDB struct {
	DB			*mgo.Database
	Config		Specification
	Session 	*mgo.Session
	Resources	*Resources
}

func (m MongoDB) UserCollectionAndSession() (*mgo.Collection, *mgo.Session, error) {
	if m.Session != nil {
		session := m.Session.Copy()
		return session.DB(m.Config.MongoDB).C(m.Config.MongoCollUsers), session, nil
	} else {
		m.Resources.Log.Debug().Msgf("No original session found")
		return nil, nil, errors.New("No original session found")
	}
}

func (m MongoDB) BookingCollectionAndSession() (*mgo.Collection, *mgo.Session, error) {
	if m.Session != nil {
		session := m.Session.Copy()
		return session.DB(m.Config.MongoDB).C(m.Config.MongoCollBookings), session, nil
	} else {
		m.Resources.Log.Debug().Msgf("No original session found")
		return nil, nil, errors.New("No original session found")
	}
}

func (r *Resources) initMongo() error {
	r.Log.Debug().Msg("Connecting to local mongo server....")

	session, err := mgo.Dial(r.Config.MongoHost)
	if err != nil {
		r.Log.Panic().Msgf("Error mongodb connection to %s: %s", r.Config.MongoHost, err.Error())
		return err
	}

	session.SetMode(mgo.Monotonic, true)

	if err != nil {
		r.Log.Panic().Msgf("Error occured while creating mongodb connection: %s", err.Error())
		return err
	}
	r.Log.Debug().Msgf("Connection established to mongo server %s", r.Config.MongoHost)

	r.Mongo = MongoDB{session.DB(r.Config.MongoDB), r.Config, session, r}
	r.Log.Debug().Msgf("Connection established to mongo db \"%s\"", r.Config.MongoDB)
	r.Log.Debug().Msgf("Connection established to mongo collections \"%s\"", r.Config.MongoCollUsers)

	return nil
}
