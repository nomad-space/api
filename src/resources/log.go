package resources


import (
	"github.com/rs/zerolog"
	zlog "github.com/rs/zerolog/log"
	"os"
	"log"
)

func (r *Resources) initLog() error {

	zerolog.TimeFieldFormat = "2006-01-02 15:04:05.000000"

	// Default level for this example is info, unless debug flag is present
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if r.Config.Debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	//instanceResources.Log = zerolog.New(os.Stderr).With().Timestamp().Logger()
	instanceResources.Log = zlog.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	log.SetFlags(0)
	log.SetOutput(instanceResources.Log)

	log.Print("hello world")
	return nil
}
