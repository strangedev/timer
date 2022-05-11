package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime/debug"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/strangedev/worktimer/lib"
)

var (
	debugging bool
	config    lib.Config
)

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().Bool("debug", true, "print debuggin output")
	viper.BindPFlag("debug", rootCmd.PersistentFlags().Lookup("debug"))
}

func initConfig() {
	userConfigDir, err := os.UserConfigDir()
	if err != nil {
		log.Fatal(err)
	}
	configDir := filepath.Join(userConfigDir, "worktimer")
	configName := "config"
	configType := "toml"
	configPath := filepath.Join(configDir, configName+"."+configType)

	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath(configDir)

	viper.SetDefault("CheckpointInterval", time.Hour)
	viper.SetDefault("StorageDir", configDir)
	viper.SetDefault("IdleDuration", time.Minute*5)
	viper.SetDefault("IdleWatcherInterval", time.Minute)
	viper.SetDefault("WakeWatcherInterval", time.Second*30)
	viper.SetDefault("debug", false)

	viper.SetEnvPrefix("worktimer")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			if err := viper.WriteConfig(); err != nil {
				if _, err := os.Stat(configDir); os.IsNotExist(err) {
					if err = os.Mkdir(configDir, 0755); err != nil {
						log.Fatal(err)
					}
				}

				if err := viper.WriteConfigAs(configPath); err != nil {
					log.Fatalf("Could not write config file: %v", err)
				}
			}
		} else {
			log.Fatalf("Could not load config file: %v", err)
		}
	}

	config = lib.Config{
		CheckpointInterval:  viper.GetDuration("CheckpointInterval"),
		StorageDir:          viper.GetString("StorageDir"),
		IdleDuration:        viper.GetDuration("IdleDuration"),
		IdleWatcherInterval: viper.GetDuration("IdleWatcherInterval"),
		WakeWatcherInterval: viper.GetDuration("WakeWatcherInterval"),
	}
	debugging = viper.GetBool("debug")
}

var rootCmd = &cobra.Command{
	Use:   "worktimer",
	Short: "worktimer counts working hours",
	Long:  `worktimer counts working hours by looking at idle times in X11.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if debugging {
			buildInfo, ok := debug.ReadBuildInfo()
			if !ok {
				fmt.Printf("Failed to read build info")
				return
			}

			fmt.Println(buildInfo)
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
