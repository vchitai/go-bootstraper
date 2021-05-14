package cli

import (
	"fmt"
	"log"

	"go-bootstraper/configs"
	"go-bootstraper/internal/kits"
	"go-bootstraper/internal/utils"

	"github.com/spf13/cobra"
)

const appName = utils.Name("go-bootstraper")

// An ext is an extension for an app
type ext func(app *app) *cobra.Command

var extensions []ext

func registerExt(ext ext) {
	extensions = append(extensions, ext)
}

type app struct {
	name     utils.Name
	cli      *cobra.Command
	worker   *kits.Worker
	resource *configs.App
}

func (app *app) Name() string {
	return app.name.UpperCamelCase().String()
}

func (app *app) init() error {
	if app.worker == nil {
		return fmt.Errorf("cannot init")
	}
	for _, ext := range extensions {
		app.cli.AddCommand(ext(app))
	}
	return nil
}

func NewApp(cfg *configs.App) (*app, error) {
	var app = &app{
		name:     appName,
		worker:   kits.NewWorker(cfg),
		resource: cfg,
	}
	// RootCmd represents the base command when called without subcommands.
	rootCmd := &cobra.Command{
		Use:   fmt.Sprintf("%s", app.name.LowerDashNotation()),
		Short: fmt.Sprintf("%s is born to help engineer deploy a new services to production in 5 minutes", app.name.UpperCamelCase()),
	}
	app.cli = rootCmd

	// versionCmd represents the version command
	var versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Current version of the application",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Printf("version %s\n", cfg.Version)
		},
	}
	rootCmd.AddCommand(versionCmd)

	if err := app.init(); err != nil {
		return nil, err
	}
	return app, nil
}

func (app *app) Run() {
	defer app.teardown()
	if err := app.cli.Execute(); err != nil {
		log.Fatal("error", err)
	}
}

func (app *app) teardown() {
	if err := app.worker.Close(); err != nil {
		log.Fatal("error", err)
	}
}
