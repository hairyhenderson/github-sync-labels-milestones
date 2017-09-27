package cmd

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path"

	"github.com/hairyhenderson/github-sync-labels-milestones/config"
	"github.com/hairyhenderson/github-sync-labels-milestones/sync"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cachePath string

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   os.Args[0],
	Short: "Manage labels and milestones across many repos",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			log.Printf("no config file specified - defaulting to 'config.json'\n")
			args = []string{"config.json"}
		} else if len(args) > 1 {
			log.Printf("error: too many arguments specified - provide only the config filename (%+v)\n", args)
			return errors.New("too many args")
		}
		cachePath, err := homedir.Expand(viper.GetString("cache"))
		if err != nil {
			return err
		}
		opts := sync.Options{
			DryRun:    viper.GetBool("dry-run"),
			NoCache:   viper.GetBool("no-cache"),
			CachePath: cachePath,
		}

		c, err := config.ParseFile(args[0])
		if err != nil {
			return err
		}

		err = sync.Sync(c, opts)
		return err
	},
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	defaultCache := path.Join("~", "."+path.Base(os.Args[0]), "cache")
	RootCmd.Flags().StringP("cache", "", defaultCache, "path to HTTP cache")
	viper.BindPFlag("cache", RootCmd.Flags().Lookup("cache"))

	RootCmd.Flags().BoolP("no-cache", "", false, "bypass HTTP caching")
	viper.BindPFlag("no-cache", RootCmd.Flags().Lookup("no-cache"))

	RootCmd.Flags().BoolP("dry-run", "", false, "do not alter the remote GitHub resources, but show what would be done")
	viper.BindPFlag("dry-run", RootCmd.Flags().Lookup("dry-run"))
}
