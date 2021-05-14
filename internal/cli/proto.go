package cli

import (
	"sync"

	"github.com/spf13/cobra"
)

var addProtoCommandOnce sync.Once

func init() {
	addProtoCommandOnce.Do(func() {
		//registerExt(protoCommand)
	})
}
func protoCommand(app *app) *cobra.Command {
	const genFilePath = "internal/protos/gen.go"

	// protoAddCmd represents the add proto command
	var protoAddCmd = &cobra.Command{
		Use:   "add",
		Short: "Create a new proto",
		Run: func(cmd *cobra.Command, args []string) {
		},
	}

	protoAddCmd.Flags().StringP("crud", "x", "", "add crud interface")
	protoAddCmd.Flags().StringP("models", "m", "", "add message")

	var protoRemoveCmd = &cobra.Command{
		Use:   "remove",
		Short: "Remove a proto",
		Run: func(cmd *cobra.Command, args []string) {
		},
	}

	var protoUpdateCmd = &cobra.Command{
		Use:   "update",
		Short: "update proto impl",
		Run: func(cmd *cobra.Command, args []string) {
		},
	}

	// protoCmd represents the proto command
	var protoCmd = &cobra.Command{
		Use:   "proto",
		Short: "Manage proto, including create-mapping",
	}
	protoCmd.AddCommand(protoAddCmd)
	protoCmd.AddCommand(protoRemoveCmd)
	protoCmd.AddCommand(protoUpdateCmd)
	return protoCmd
}
