package cmd

import (
	"fmt"
	"log"
	"net/rpc"

	"github.com/spf13/cobra"
	"github.com/strangedev/worktimer/lib"
)

func init() {
	rootCmd.AddCommand(stopCommand)
}

var stopCommand = &cobra.Command{
	Use:   "stop",
	Short: "Stop the timer",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := rpc.DialHTTP("tcp", "localhost:1234")
		if err != nil {
			log.Fatal("Dialing:", err)
		}

		stopArgs := &lib.StartStopArgs{
			Reason: "Manual stop",
		}
		var reply lib.VoidReply
		err = client.Call("Daemon.StopTimer", stopArgs, &reply)
		if err != nil {
			log.Fatal("rpc error: ", err)
		}
		fmt.Println("Stopped the timer.")
	},
}
