package cmd

import (
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
)

var rootFolder string

// Execute exposes the execution trigger of the root command to the outside
// world.
func Execute() error {
	return rootCmd.Execute()
}

var log *logrus.Logger
var verbose bool

func init() {
	log = logrus.New()
	log.SetLevel(logrus.InfoLevel)
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	root, err := findRootFolder(cwd)
	if err != nil {
		panic(err)
	}
	rootCmd.PersistentFlags().StringVar(&rootFolder, "root-folder", root, "Path to the gopher.community root folder")
	rootCmd.PersistentFlags().BoolVar(&verbose, "verbose", false, "Verbose logging")
}

var rootCmd = &cobra.Command{
	Use:   "gocom",
	Short: "A utility for working on and with the gopher.community.",
	Run: func(cmd *cobra.Command, args []string) {
		if verbose {
			log.SetLevel(logrus.DebugLevel)
		}
	},
}
