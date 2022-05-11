package cmd

import (
	"fmt"
	"log"
	"net/rpc"

	"github.com/strangedev/worktimer/lib"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(noteCommand)
}

var noteCommand = &cobra.Command{
	Use:   "note",
	Short: "Adds a note to the current time slice",
	Long:  "Adds a note to the current time slice, wakes the timer if it is suspended.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client, err := rpc.DialHTTP("tcp", "localhost:1234")
		if err != nil {
			log.Fatal("Dialing:", err)
		}

		noteArgs := &lib.NoteArgs{
			Note: args[0],
		}
		var reply lib.VoidReply
		err = client.Call("Daemon.AddNote", noteArgs, &reply)
		if err != nil {
			log.Fatal("rpc error: ", err)
		}
		fmt.Println("Added the note.")
	},
}
