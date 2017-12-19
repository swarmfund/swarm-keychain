package main

import (
	"github.com/spf13/cobra"
	"gitlab.com/swarmfund/keychain"
	"gitlab.com/swarmfund/keychain/config"
	"gitlab.com/swarmfund/keychain/log"
)

var (
	configFile     string
	configInstance config.Config
	rootCmd        = &cobra.Command{
		Use: "keychain",
	}
	runCmd = &cobra.Command{
		Use: "run",
		Run: func(cmd *cobra.Command, args []string) {
			app, err := keychain.NewApp(configInstance)
			if err != nil {
				log.WithField("service", "init").WithError(err).Fatal("failed to init app")
			}
			app.Serve()
		},
	}
)

func main() {
	cobra.OnInitialize(func() {
		c, err := initConfig(configFile)
		if err != nil {
			log.WithField("service", "init").WithError(err).Fatal("failed to init config")
		}
		configInstance = c
	})
	rootCmd.PersistentFlags().StringVar(&configFile, "config", "config.yaml", "config file")
	rootCmd.AddCommand(runCmd)
	rootCmd.AddCommand(migrateCmd)
	rootCmd.Execute()
}

func initConfig(fn string) (config.Config, error) {
	c := config.NewViperConfig(fn)
	if err := c.Init(); err != nil {
		return nil, err
	}
	return c, nil
}
