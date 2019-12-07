package cmd

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/sno6/reaper/pkg/reaper"
	"github.com/spf13/cobra"
)

const defTimeout = time.Second * 15

func RootCmd() *cobra.Command {
	var (
		searchTerms   []string
		outDir        string
		verbose       bool
		timeout       time.Duration
		width, height int
	)

	rootCmd := &cobra.Command{
		Use:   "reaper",
		Short: "Generate a dataset for a given topic",
		Run: func(cmd *cobra.Command, args []string) {
			if len(searchTerms) == 0 {
				fmt.Printf("Reaper needs a search term in order to work\n\n")
				cmd.Usage()
				os.Exit(1)
			}

			r := reaper.New(&reaper.Config{
				SearchTerms: searchTerms,
				OutDir:      outDir,
				Verbose:     verbose,
				Timeout:     timeout,
				MinWidth:    width,
				MinHeight:   height,
			})
			if err := r.Run(); err != nil {
				log.Fatal(err)
			}
		},
	}

	rootCmd.Flags().StringArrayVarP(&searchTerms, "search", "s", []string{}, "Search terms")
	rootCmd.Flags().StringVarP(&outDir, "out", "o", "./", "Output directory")
	rootCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Log errors")
	rootCmd.Flags().DurationVarP(&timeout, "timeout", "t", defTimeout, "HTTP request timeout in seconds")
	rootCmd.Flags().IntVar(&width, "mw", -1, "Only allow images width a width >= to this value")
	rootCmd.Flags().IntVar(&height, "mh", -1, "Only allow images width a height >= to this value")

	return rootCmd
}
