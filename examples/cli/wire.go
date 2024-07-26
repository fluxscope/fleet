//go:build wireinject
// +build wireinject

package main

import (
	"github.com/fluxscope/fleet"
	"github.com/google/wire"
)

func wireApp() (*fleet.Console, func(), error) {
	wire.Build(ProvideSet)
	return &fleet.Console{}, nil, nil
}
