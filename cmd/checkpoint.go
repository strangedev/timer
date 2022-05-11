package cmd

import (
	"fmt"
	"log"
	"net/rpc"

	"github.com/strangedev/worktimer/lib"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(checkpointCommand)
}

var checkpointCommand = &cobra.Command{
	Use:   "checkpoint",
	Short: "Saves the current state to disk",
	Long:  "Starts a new slice and dumps all past slices to disk.",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := rpc.DialHTTP("tcp", "localhost:1234")
		if err != nil {
			log.Fatal("Dialing:", err)
		}

		checkpointArgs := &lib.VoidArgs{}
		var reply lib.VoidReply
		err = client.Call("Daemon.TriggerCheckpoint", checkpointArgs, &reply)
		if err != nil {
			log.Fatal("rpc error: ", err)
		}
		fmt.Println("Created a checkpoint.")
	},
}
