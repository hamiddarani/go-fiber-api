package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/hamiddarani/web-api-fiber/cmd"
	"github.com/spf13/cobra"
)

// @securityDefinitions.apikey AuthBearer
// @in header
// @name Authorization
func main() {
	const description = "Simple web api with golang and fiber"

	root := &cobra.Command{Short: description}

	trap := make(chan os.Signal, 1)
	signal.Notify(trap, syscall.SIGINT, syscall.SIGTERM)

	root.AddCommand(
		cmd.Server{}.Command(trap),
	)

	if err := root.Execute(); err != nil {
		log.Fatalf("failed to execute root command\n%v", err)
	}

}
