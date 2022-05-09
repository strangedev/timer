package cmd

import (
	"fmt"
	"log"
	"net/rpc"
	"timer/lib"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(startCommand)
}

var startCommand = &cobra.Command{
	Use:   "start",
	Short: "Start the timer",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := rpc.DialHTTP("tcp", "localhost:1234")
		if err != nil {
			log.Fatal("Dialing:", err)
		}

		startArgs := &lib.StartStopArgs{
			Reason: "Manual start",
		}
		var reply lib.VoidReply
		err = client.Call("Timer.RPC_Start", startArgs, &reply)
		if err != nil {
			log.Fatal("rpc error: ", err)
		}
		fmt.Println("Started the timer.")
	},
}
