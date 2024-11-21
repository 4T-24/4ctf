package cmd

import (
	v1 "4ctf/api/v1"
	"4ctf/config"
	"fmt"
	"log"

	"github.com/savsgio/atreugo/v11"
	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the API server",
	Run: func(cmd *cobra.Command, args []string) {
		config := config.LoadConfig()
		server := atreugo.New(atreugo.Config{
			Addr: fmt.Sprintf(":%d", config.Server.Port),
		})

		v1.SetupRoutes(server)

		log.Printf("Starting server on port %d", config.Server.Port)
		if err := server.ListenAndServe(); err != nil {
			log.Fatalf("Error starting server: %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
