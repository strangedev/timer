package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/strangedev/worktimer/lib"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(daemonCommmand)
}

var daemonCommmand = &cobra.Command{
	Use:   "daemon",
	Short: "Starts the timer daemon.",
	Run: func(cmd *cobra.Command, args []string) {
		signals := make(chan os.Signal, 1)
		signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
		stop := make(chan bool, 1)
		go func() {
			stopSignal := <-signals
			fmt.Println()
			fmt.Println(stopSignal)
			stop <- true
		}()

		fmt.Printf("Using config %+v\n", config)

		daemon := lib.CreateDaemon(config)
		daemon.Run(stop)

		fmt.Println("Exiting")
	},
}
