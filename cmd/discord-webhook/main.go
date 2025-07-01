package main

import (
	"os"

	"github.com/yubiquita/discord-webhook/internal/cli"
)

func main() {
	rootCmd := cli.NewRootCommand()
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}