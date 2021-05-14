package cli

import (
	"log"
	"strings"
	"sync"
	"time"

	"github.com/spf13/cobra"
)

var addIncubateCommandOnce sync.Once

func init() {
	addIncubateCommandOnce.Do(func() {
		registerExt(incubateCommand)
	})
}
func incubateCommand(app *app) *cobra.Command {
	worker := app.worker
	// incubateCmd represents the incubate command
	return &cobra.Command{
		Use:   "incubate",
		Short: "Instantly initiate your project, customizable by configs",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) != 1 {
				log.Fatal("only accept single parameter in format [service_domain]/[service_name], etc platform/proj")
			}
			tokens := strings.Split(args[0], "/")
			if len(tokens) != 2 {
				log.Fatal("only accept single parameter in format [service_domain]/[service_name], etc platform/proj")
			}
			domain, name := tokens[0], tokens[1]
			pre := time.Now()
			if err := worker.Boostrap(name, domain); err != nil {
				log.Fatal("cannot incubate ", err)
			}
			log.Printf("done in %s\n", time.Since(pre).String())
		},
	}
}
