package cli

import (
	"sync"
	"time"

	"go-bootstraper/internal/kits"
	"go-bootstraper/internal/store"

	"github.com/spf13/cobra"
)

var addGenerateCommandOnce sync.Once

func init() {
	addGenerateCommandOnce.Do(func() {
		registerExt(generateCommand)
	})
}
func generateCommand(app *app) *cobra.Command {
	return &cobra.Command{
		Use:   "generate",
		Short: "Generate files according to config file",
		Run: func(cmd *cobra.Command, args []string) {
			pre := time.Now()

			drawing, err := store.NewSketch("").Load()
			if err != nil {
				cmd.PrintErr(err)
				return
			}
			b, err := kits.NewBuildTeam(app.resource, drawing)
			if err != nil {
				cmd.PrintErr(err)
				return
			}
			err = b.Build()
			if err != nil {
				cmd.PrintErr(err)
				return
			}

			cmd.Printf("done in %s\n", time.Since(pre).String())
		},
	}
}
