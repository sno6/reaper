package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/sno6/reaper/pkg/reaper"
	"github.com/spf13/cobra"
)

func RootCmd() *cobra.Command {
	var searchTerms []string

	rootCmd := &cobra.Command{
		Use:   "reaper",
		Short: "Convert Sauce 3 modules to Sauce 5 modules.",
		Run: func(cmd *cobra.Command, args []string) {
			if len(searchTerms) == 0 {
				fmt.Println("Reaper needs a search term in order to work")
				cmd.Usage()
				os.Exit(1)
			}

			r := reaper.New(&reaper.Config{
				SearchTerms: searchTerms,
				OutDir:      cmd.Flag("out").Value.String(),
			})
			if err := r.Run(); err != nil {
				log.Fatal(err)
			}

		},
	}

	rootCmd.Flags().StringArrayVarP(&searchTerms, "term", "t", []string{}, "Search terms")
	rootCmd.Flags().StringP("out", "o", "./", "Output directory")

	// Min Width.
	// Min Height.
	// Augmentation
	// -v for errors
	// time out flag
	// body size flag

	return rootCmd
}
