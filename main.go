package main

import (
	_ "embed"
	"fmt"
	"os"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
	"github.com/urfave/cli/v2"

	"github.com/mipmip/dirtygit/scanner"
	"github.com/mipmip/dirtygit/ui"
)

func getDefaultConfigPath() string {
	home, err := homedir.Dir()
	_ = err // ignore
	return filepath.Join(home, ".dirtygit.yml")
}

//go:embed .dirtygit.yml
var defaultConfig string

func main() {
	app := cli.NewApp()
	app.Name = "dirtygit"
	app.Usage = "Finds git repos in need of commitment"
	app.EnableBashCompletion = true
	app.CommandNotFound = func(c *cli.Context, cmd string) {
		fmt.Printf("ERROR: Unknown command '%s'\n", cmd)
	}
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:    "config",
			Aliases: []string{"c"},
			Usage:   "Location of config file",
			Value:   getDefaultConfigPath(),
		},

		&cli.BoolFlag{
			Name:  "ignore_dir_errors",
			Aliases: []string{"i"},
			Value: true,
			Usage: "Don't halt on errors while finding dirs",
		},
		&cli.BoolFlag{
			Name:  "debug",
			Usage: "show debug output instead of UI",
		},
	}
	app.Action = func(c *cli.Context) error {

		config, err := scanner.ParseConfigFile(c.String("config"), defaultConfig)
		if c.Args().Len() > 0 {

			fmt.Println("Arguments given, skipping config")
			config.ScanDirs.Include = c.Args().Slice()

		} else {
			if err != nil {
				return err
			}
			for i := range config.ScanDirs.Include {
				config.ScanDirs.Include[i] = os.ExpandEnv(config.ScanDirs.Include[i])
			}
			for i := range config.ScanDirs.Exclude {
				config.ScanDirs.Exclude[i] = os.ExpandEnv(config.ScanDirs.Exclude[i])
			}
		}

		if c.Bool("debug") {
			var mgs scanner.MultiGitStatus
			mgs, err = scanner.Scan(config, c.Bool("ignore_dir_errors"))
			if err != nil {
				panic(err)
			}

			for r, st := range mgs {
				fmt.Printf("%-40s %v\n", r, st.ScanTime)
			}
			return nil
		}

		err = ui.Run(config, c.Bool("ignore_dir_errors"))
		if err != nil {
			return err
		}
		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Printf("%+v\n", err)
	}
}
