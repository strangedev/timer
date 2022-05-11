package cmd

import (
	"fmt"
	"log"
	"net/rpc"
	"os"
	"reflect"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/strangedev/worktimer/lib"
)

func init() {
	rootCmd.AddCommand(statusCommand)
}

var statusCommand = &cobra.Command{
	Use:   "status",
	Short: "Show the timer's status",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := rpc.DialHTTP("tcp", "localhost:1234")
		if err != nil {
			log.Fatal("Dialing:", err)
		}

		statusArgs := &lib.VoidArgs{}
		var reply lib.StatusReply
		err = client.Call("Daemon.GetStatus", statusArgs, &reply)

		out := tabwriter.NewWriter(os.Stdout, 0, 2, 1, ' ', 0)

		value := reflect.ValueOf(reply)
		typeOf := value.Type()

		for i := 0; i < value.NumField(); i++ {
			fmt.Fprintf(out, "%s\t%v\n", typeOf.Field(i).Name, value.Field(i).Interface())
		}

		out.Flush()

		if err != nil {
			log.Fatal("rpc error: ", err)
		}
	},
}
