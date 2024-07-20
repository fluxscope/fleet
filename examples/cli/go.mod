module github.com/fluxscope/fleet/examples/cli

go 1.22.0

toolchain go1.22.2

replace github.com/fluxscope/fleet => ../..

require (
	emperror.dev/emperror v0.33.0
	github.com/fluxscope/fleet v0.0.0-20240720132526-fd404c7c6d40
	github.com/google/wire v0.6.0
	github.com/spf13/cobra v1.8.1
)

require (
	emperror.dev/errors v0.8.0 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/oklog/run v1.1.0 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/samber/lo v1.44.0 // indirect
	github.com/samber/slog-common v0.17.0 // indirect
	github.com/samber/slog-zap/v2 v2.6.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	go.uber.org/zap v1.26.0 // indirect
	golang.org/x/text v0.16.0 // indirect
	google.golang.org/protobuf v1.34.1 // indirect
)
