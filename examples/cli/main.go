package main

import (
	"emperror.dev/emperror"
	"github.com/spf13/cobra"
)

func main() {

	rootCmd := &cobra.Command{
		Use:   "fleet-cli",
		Short: "fleet cli demo",
		RunE: func(cmd *cobra.Command, args []string) error {
			app, cancel, err := wireApp()
			if err != nil {
				panic(err)
			}
			defer cancel()
			return app.RunE(args)
		},
	}

	emperror.Panic(rootCmd.Execute())
}
