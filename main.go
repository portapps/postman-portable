//go:generate go install -v github.com/josephspurrier/goversioninfo/cmd/goversioninfo
//go:generate goversioninfo -icon=res/papp.ico -manifest=res/papp.manifest
package main

import (
	"os"
	"path"

	"github.com/portapps/portapps/v3"
	"github.com/portapps/portapps/v3/pkg/log"
	"github.com/portapps/portapps/v3/pkg/utl"
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
	utl.CreateFolder(app.DataPath)

	electronAppFolder, err := utl.FindElectronAppFolder("app-", app.AppPath)
	if err != nil {
		log.Fatal().Msgf("Electron main folder not found")
	}
	electronBinPath := utl.PathJoin(app.AppPath, electronAppFolder)

	app.Process = utl.PathJoin(electronBinPath, "Postman.exe")
	app.WorkingDir = electronBinPath
	app.Args = []string{
		"--user-data-dir=" + app.DataPath,
	}

	// Cleanup on exit
	if cfg.Cleanup {
		defer func() {
			utl.Cleanup([]string{
				path.Join(os.Getenv("APPDATA"), "Postman"),
			})
		}()
	}

	defer app.Close()
	app.Launch(os.Args[1:])
}
