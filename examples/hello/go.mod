module github.com/fluxscope/fleet/examples/hello

go 1.22.0

toolchain go1.22.2

replace github.com/fluxscope/fleet => ../..

require (
	github.com/fluxscope/fleet v0.0.0-20240720132526-fd404c7c6d40
	github.com/google/wire v0.6.0
	github.com/gorilla/mux v1.8.1
)

require (
	emperror.dev/emperror v0.33.0 // indirect
	emperror.dev/errors v0.8.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/oklog/run v1.1.0 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/samber/lo v1.44.0 // indirect
	github.com/samber/slog-common v0.17.0 // indirect
	github.com/samber/slog-zap/v2 v2.6.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	go.uber.org/zap v1.26.0 // indirect
	golang.org/x/text v0.16.0 // indirect
	google.golang.org/protobuf v1.34.1 // indirect
)
