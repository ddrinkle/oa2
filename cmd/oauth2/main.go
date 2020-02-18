package main

import (
	"github.com/ddrinkle/oa2/internal/handlers"
	"github.com/ddrinkle/oa2/internal/routes"
	"github.com/ddrinkle/platform/app"
	"github.com/ddrinkle/platform/config"
	"github.com/ddrinkle/platform/crypt"
	"encoding/json"
	"os"
	"strconv"

	"github.com/fvbock/endless"
)

type OA2Config struct {
	config.Config
	DBEncryptionKey string `json:"db_encryption_key"`
}

// LoadFile loads a config file, and loads it into the passed in struct
func (c *OA2Config) LoadFile(filename string) error {
	file, err := os.Open(filename)
	defer file.Close()

	if err != nil {
		return err
	}
	decoder := json.NewDecoder(file)

	err = decoder.Decode(c)
	if err != nil {
		return err
	}
	return nil
}

func main() {

	conf := OA2Config{}

	application := handlers.App{}
	application.App = app.New(&conf)
	application.Storage.DB = application.Env.DB

	crypt.SetCryptKeyFunction(func() []byte {
		return []byte(conf.DBEncryptionKey)
	})

	application.AddRoutes(routes.GetRoutes(application))
	mux := application.NewRouter()

	endless.ListenAndServe(":"+strconv.Itoa(conf.App.Port), mux)
}
