package cmd

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Short: "Word Battle Ecosystem",
}

func init() {
	cobra.OnInitialize()
}

// Execute the application.
func Execute() {
	godotenv.Load()

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
