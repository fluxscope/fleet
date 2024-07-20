package main

import "github.com/spf13/cobra"

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
			app.Run()
			return nil
		},
	}

	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
