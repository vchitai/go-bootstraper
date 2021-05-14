package cli

import (
	"sync"

	"github.com/spf13/cobra"
)

var addServiceCommandOnce sync.Once

func init() {
	addServiceCommandOnce.Do(func() {
		//registerExt(serviceCommand)
	})
}
func serviceCommand(app *app) *cobra.Command {
	// serviceAddCmd represents the add service command
	var serviceAddCmd = &cobra.Command{
		Use:   "add",
		Short: "Create a new service",
		Run: func(cmd *cobra.Command, args []string) {

		},
	}
	serviceAddCmd.Flags().StringP("crud", "x", "", "add crud")
	serviceAddCmd.Flags().StringP("implement", "p", "", "implement of protos")
	serviceAddCmd.Flags().StringP("stores", "s", "", "connected stores")
	serviceAddCmd.Flags().StringP("init", "i", "", "add init cmd")

	var serviceRemoveCmd = &cobra.Command{
		Use:   "remove",
		Short: "Remove a service",
		Run: func(cmd *cobra.Command, args []string) {

		},
	}

	var servicesCmd = &cobra.Command{
		Use:   "services",
		Short: "Do a service",
		Run: func(cmd *cobra.Command, args []string) {

		},
	}
	servicesCmd.AddCommand(serviceAddCmd)
	servicesCmd.AddCommand(serviceRemoveCmd)
	return servicesCmd

}
