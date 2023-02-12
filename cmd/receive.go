package cmd

import (
	"github.com/SamoKopecky/pqcom/main/network"
	"github.com/spf13/cobra"
)

var (
	stdout bool
	dir    string

	receiveCmd = &cobra.Command{
		Use:   "receive",
		Short: "Use receive mode",
		Long: `Use the application in receive mode to receive data.
By default the stdout of the app is the destination of any receiving data.`,
		Run: func(cmd *cobra.Command, args []string) {
			network.Receive(destAddr, srcPort, destPort, dir)
		},
	}
)

func init() {
	appCmd.AddCommand(receiveCmd)
	receiveCmd.Flags().StringVar(&dir, "dir", "", "Receive data and save to files in dir")
}
