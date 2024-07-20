package main

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"strings"
)

func main() {
	v, f := viper.New(), pflag.NewFlagSet("fleet cli demo", pflag.ExitOnError)
	f.String("config", "./configs", "Configuration path")
	f.Bool("version", false, "Show version information")

	configure(v, f)
	app, cancel, err := wireApp()
	if err != nil {
		panic(err)
	}
	defer cancel()
	app.Run()
}

// configure configures some defaults in the Viper instance.
func configure(v *viper.Viper, f *pflag.FlagSet) {
	// Viper settings
	v.AddConfigPath(".")
	v.AddConfigPath("$CONFIG_DIR/")

	// Environment variable settings
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	v.AllowEmptyEnv(true)
	v.AutomaticEnv()
}
