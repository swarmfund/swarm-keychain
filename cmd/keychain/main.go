package main

import (
	"log"
	"runtime"

	"gitlab.com/distributed_lab/tokend/keychain"
	"gitlab.com/distributed_lab/tokend/keychain/config"
	"github.com/spf13/cobra"
)

var app *keychain.App
var conf config.Config
var version string

var rootCmd *cobra.Command

func main() {
	if version != "" {
		keychain.SetVersion(version)
	}
	runtime.GOMAXPROCS(runtime.NumCPU())
	rootCmd.Execute()
}

func init() {

	rootCmd = &cobra.Command{
		Use: "keychain",
		Run: func(cmd *cobra.Command, args []string) {
			initApp(cmd, args)
			app.Serve()
		},
	}

	conf.DefineConfigStructure(rootCmd)

	rootCmd.AddCommand(dbCmd)
}

func initApp(cmd *cobra.Command, args []string) {
	err := conf.Init()
	if err != nil {
		log.Println("Failed to init config")
		log.Fatal(err.Error())
	}
	app, err = keychain.NewApp(conf)

	if err != nil {
		log.Fatal(err.Error())
	}
}
