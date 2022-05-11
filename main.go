package main

import (
	"fmt"
	"runtime/debug"
	"worktimer/cmd"
)

func main() {
	buildInfo, ok := debug.ReadBuildInfo()
	if !ok {
		fmt.Printf("Failed to read build info")
		return
	}

	fmt.Printf("Version: %v\n", buildInfo.Main.Version)

	cmd.Execute()
}
