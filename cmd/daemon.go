package cmd

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"
	"timer/lib"

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

		configDir, err := os.UserConfigDir()
		if err != nil {
			log.Fatal(err)
		}

		config := lib.Config{
			CheckpointInterval:  time.Hour,
			StorageDir:          filepath.Join(configDir, ".timer"),
			IdleDuration:        time.Minute * 5,
			IdleWatcherInterval: time.Minute,
			WakeWatcherInterval: time.Second * 30,
		}
		fmt.Printf("Using config %v\n", config)

		if _, err := os.Stat(config.StorageDir); os.IsNotExist(err) {
			if err = os.Mkdir(config.StorageDir, 0755); err != nil {
				log.Fatal(err)
			}
		}

		timer := lib.CreateTimer()
		rpc.Register(timer)
		rpc.HandleHTTP()

		socket, err := net.Listen("tcp", "localhost:1234")
		if err != nil {
			log.Fatal("listen error:", err)
		}

		go func() {
			fmt.Println("Starting RPC server")
			go http.Serve(socket, nil)
		}()

		go func() {
			fmt.Println("Starting checkpointing routine")
			for {
				time.Sleep(config.CheckpointInterval)
				fmt.Println("Checkpointing")
				if checkpointErr := timer.Checkpoint(config); checkpointErr != nil {
					fmt.Printf("%+v\n", timer)
					log.Fatal(checkpointErr)
				}
			}
		}()

		go func() {
			fmt.Println("Starting X11 idle watcher")
			for {
				time.Sleep(config.IdleWatcherInterval)
				if timer.IsStopped() {
					continue
				}

				idleTime, err := lib.GetIdleTime()
				if err != nil {
					log.Fatal(err)
				}

				if idleTime > config.IdleDuration {
					if !timer.IsSuspended() {
						fmt.Println("Suspending timer due to X11 idle")
						timer.Suspend("X11 idle")
					}
				}
			}
		}()
		go func() {
			fmt.Println("Starting X11 wake watcher")
			for {
				time.Sleep(config.WakeWatcherInterval)
				if timer.IsStopped() {
					continue
				}

				currentIdleTime, err := lib.GetIdleTime()
				if err != nil {
					log.Fatal(err)
				}

				if currentIdleTime < config.IdleDuration {
					if timer.IsSuspended() {
						fmt.Println("Waking timer due to X11 activity")
						timer.Wake("X11 activity")
					}
				}
			}
		}()

		<-stop
		if checkpointErr := timer.Checkpoint(config); checkpointErr != nil {
			fmt.Printf("%+v\n", timer)
			log.Fatal(checkpointErr)
		}

		fmt.Println("Exiting")
	},
}
