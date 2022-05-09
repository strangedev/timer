package cmd

import (
	"fmt"
	"log"
	"net/rpc"
	"worktimer/lib"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(wakeCommand)
}

var wakeCommand = &cobra.Command{
	Use:   "wake",
	Short: "Wakes the timer",
	Long:  "Starts the timer from suspension, usually you want to use start instead, but this can be used to integrate other programs that track activity.",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := rpc.DialHTTP("tcp", "localhost:1234")
		if err != nil {
			log.Fatal("Dialing:", err)
		}

		wakeArgs := &lib.StartStopArgs{
			Reason: "Manual wake",
		}
		var reply lib.VoidReply
		err = client.Call("Timer.RPC_Wake", wakeArgs, &reply)
		if err != nil {
			log.Fatal("rpc error: ", err)
		}
		fmt.Println("Woke the timer.")
	},
}
