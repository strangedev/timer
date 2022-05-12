package cmd

import (
	"fmt"
	"log"
	"net/rpc"

	"github.com/strangedev/worktimer/lib"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(suspendCommand)
}

var suspendCommand = &cobra.Command{
	Use:   "suspend",
	Short: "Suspend the timer",
	Long:  "Stops the timer, but instructs it to automatically restart on activity.",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := rpc.DialHTTP("tcp", "localhost:1234")
		if err != nil {
			log.Fatal("Dialing:", err)
		}

		suspendArgs := &lib.StartStopArgs{
			Reason: "Manual suspend",
		}
		var reply lib.VoidReply
		err = client.Call("Daemon.SuspendTimer", suspendArgs, &reply)
		if err != nil {
			log.Fatal("rpc error: ", err)
		}
		fmt.Println("Suspended the timer.")
	},
}
