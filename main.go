package main

//go:generate go tool goversioninfo -icon=res/papp.ico -manifest=res/papp.manifest

import (
	"os"
	"path/filepath"

	"github.com/portapps/portapps/v3"
	"github.com/portapps/portapps/v3/pkg/files"
	"github.com/portapps/portapps/v3/pkg/log"
)

type config struct {
	Cleanup bool `yaml:"cleanup" mapstructure:"cleanup"`
}

var (
	app *portapps.App
	cfg *config
)

func init() {
	var err error

	// Default config
	cfg = &config{
		Cleanup: false,
	}

	// Init app
	if app, err = portapps.NewWithCfg("postman-portable", "Postman", cfg); err != nil {
		log.Fatal().Err(err).Msg("Cannot initialize application. See log file for more info.")
	}
}

func main() {
	if err := os.MkdirAll(app.DataPath, 0o755); err != nil {
		log.Fatal().Err(err).Msg("Cannot create data path")
	}
	electronAppPath := app.ElectronAppPath()

	app.Process = filepath.Join(electronAppPath, "Postman.exe")
	app.WorkingDir = electronAppPath
	app.Args = []string{
		"--user-data-dir=" + app.DataPath,
	}

	// Cleanup on exit
	if cfg.Cleanup {
		defer func() {
			files.Cleanup(filepath.Join(os.Getenv("APPDATA"), "Postman"))
		}()
	}

	defer app.Close()
	app.Launch(os.Args[1:])
}
