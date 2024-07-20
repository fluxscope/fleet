//go:build wireinject
// +build wireinject

package main

import (
	"github.com/fluxscope/fleet"
	"github.com/google/wire"
)

//go:generate wire

func wireApp() (*fleet.App, func(), error) {
	wire.Build(ProvideSet)
	return &fleet.App{}, nil, nil
}
