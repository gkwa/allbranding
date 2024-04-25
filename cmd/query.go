package cmd

import (
	"log/slog"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/taylormonacelli/allbranding/query"
)

var (
	releasesURL string
	assetRegex  string
	noCache     bool
	parseHarder bool
	ignoreRegex []string
)

// queryCmd represents the query command
var queryCmd = &cobra.Command{
	Use:   "query",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		query.Run(releasesURL, assetRegex, noCache, parseHarder, ignoreRegex)
	},
}

func init() {
	rootCmd.AddCommand(queryCmd)

	queryCmd.PersistentFlags().StringVar(&releasesURL, "releases-url", "https://api.github.com/repos/gnprice/toml-cli/releases", "URL of the GitHub releases API endpoint")
	err := viper.BindPFlag("releases-url", queryCmd.PersistentFlags().Lookup("releases-url"))
	if err != nil {
		slog.Error("error binding releases-url flag", "error", err)
		os.Exit(1)
	}

	queryCmd.PersistentFlags().StringVar(&assetRegex, "asset-regex", `toml-v\d+\.\d+\.\d+-x86_64-linux\.tar\.gz$`, "Regular expression to match the desired asset")
	err = viper.BindPFlag("asset-regex", queryCmd.PersistentFlags().Lookup("asset-regex"))
	if err != nil {
		slog.Error("error binding asset-regex flag", "error", err)
		os.Exit(1)
	}

	queryCmd.PersistentFlags().BoolVar(&noCache, "no-cache", false, "Disable caching of the releases data")
	err = viper.BindPFlag("no-cache", queryCmd.PersistentFlags().Lookup("no-cache"))
	if err != nil {
		slog.Error("error binding no-cache flag", "error", err)
		os.Exit(1)
	}

	queryCmd.PersistentFlags().BoolVar(&parseHarder, "parse-harder", false, "Use regex to remove non-numeric characters from version strings")
	err = viper.BindPFlag("parse-harder", queryCmd.PersistentFlags().Lookup("parse-harder"))
	if err != nil {
		slog.Error("error binding parse-harder flag", "error", err)
		os.Exit(1)
	}

	queryCmd.PersistentFlags().StringSliceVar(&ignoreRegex, "ignore", []string{}, "Regex patterns to ignore versions (can be specified multiple times)")
	err = viper.BindPFlag("ignore", queryCmd.PersistentFlags().Lookup("ignore"))
	if err != nil {
		slog.Error("error binding ignore flag", "error", err)
		os.Exit(1)
	}
}
