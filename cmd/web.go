/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/rwxd/civitai-search/web"
	"github.com/spf13/cobra"
)

var (
	ip   string
	port string
)

var webCmd = &cobra.Command{
	Use:   "web",
	Short: "Run the webserver",
	Run: func(cmd *cobra.Command, args []string) {
		dbPath, err := cmd.Flags().GetString("db")
		if err != nil {
			panic(err)
		}

		web.StartServer(ip, port, dbPath)
	},
}

func init() {
	rootCmd.AddCommand(webCmd)
	webCmd.PersistentFlags().StringVarP(&ip, "ip", "i", "0.0.0.0", "IP address to bind the server")
	webCmd.PersistentFlags().StringVarP(&port, "port", "p", "8080", "Port for the server")
	webCmd.Flags().String("db", "./db.sqlite", "Path to the sqlite database")
}
